// Package handlers содержит обработчики HTTP запросов для календаря
package handlers

import (
	"calendar/internal/event"
	dto "calendar/internal/infrastructure/http/handlers/dto"
	"calendar/internal/infrastructure/http/middleware"
	"calendar/internal/infrastructure/http/request"
	"calendar/internal/infrastructure/http/response"
	inmem "calendar/internal/infrastructure/storage/in_memory"
	"calendar/pkg/sl_logger/sl"
	valResp "calendar/pkg/validator"

	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
)

// NewAddEventHandler создает новый обработчик для добавления события в календарь POST
// Принимает:
//   - log *slog.Logger: логгер для записи информации о работе обработчика
//   - svc event.Service: сервис для работы с событиями
//
// Возвращает:
//   - http.HandlerFunc: функцию-обработчик HTTP запросов
func NewAddEventHandler(log *slog.Logger, svc event.Service) http.HandlerFunc {
	// Возвращаем функцию-обработчик
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Константа для идентификации операции в логах
		const op = "handlers.event.add"

		// Добавляем в логгер информацию об операции и ID запроса
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(r)),
		)

		// Структура для декодирования запроса
		var req dto.AddEventRequest

		// Декодируем тело запроса в структуру
		err := json.NewDecoder(r.Body).Decode(&req)
		// Проверяем, не пустое ли тело запроса
		if errors.Is(err, io.EOF) {
			log.Error("bad request",
				slog.String("type", request.ErrEmptyReqBody.Error()),
				sl.Err(err),
			)
			addEventResponseErr(w, request.ErrEmptyReqBody.Error())
			
			return
		}
		if err != nil {
			log.Error("bad request",
				slog.String("type", request.ErrFailedToDecodeReqBody.Error()),
				sl.Err(err),
			)
			addEventResponseErr(w,  request.ErrFailedToDecodeReqBody.Error())
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		// Валидируем данные запроса
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			response.WriteJSON(w, http.StatusBadRequest, valResp.ValidationError(validateErr))
			return
		}

		// Создаем объект события из данных запроса
		respEvent := event.Event{
			Date:  req.Date,
			Title: req.Title,
			Desc:  req.Desc,
		}
		// Добавляем событие через сервисный слой
		if err := svc.Add(respEvent); err != nil {
			switch {
			case errors.Is(err, inmem.ErrNoValue):
				log.Error("failed to add event", sl.Err(err))
				return
			default:
				log.Error("unexpected error adding event", sl.Err(err))
				return
			}
		}

		// Логируем успешное добавление события
		log.Info("event added", slog.Any("title", respEvent.Title))

		// Отправляем успешный ответ клиенту
		addEventResponseOK(w, respEvent.Title)
	}
}

// responseOK отправляет успешный ответ клиенту
// Параметры:
//   - w http.ResponseWriter: интерфейс для записи ответа
//   - title string: заголовок добавленного события
func addEventResponseOK(w http.ResponseWriter, title string) {
	r := dto.AddEventResponse{
		ValidationResponse: valResp.OK(),
		Title:              title,
	}
	response.WriteJSON(w, http.StatusOK, r)
}

func addEventResponseErr(w http.ResponseWriter, e string) {
	r := dto.AddEventResponse{
		ValidationResponse: valResp.Error(e),
	}
	response.WriteJSON(w, http.StatusBadRequest, r)
}