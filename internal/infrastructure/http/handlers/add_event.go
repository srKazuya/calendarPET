// Package handlers содержит обработчики HTTP запросов для календаря
package handlers

import (
	"calendar/internal/event"
	dto "calendar/internal/infrastructure/http/handlers/dto"
	"calendar/internal/infrastructure/http/middleware"
	"calendar/internal/infrastructure/http/request"
	"calendar/internal/infrastructure/http/response"
	"calendar/pkg/sl_logger/sl"
	valResp "calendar/pkg/validator"

	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
)

// NewAddEventHandler создает новый обработчик для добавления события в календарь
// Принимает:
//   - log *slog.Logger: логгер для записи информации о работе обработчика
//   - svc event.Service: сервис для работы с событиями
//
// Возвращает:
//   - http.HandlerFunc: функцию-обработчик HTTP запросов
func NewAddEventHandler(log *slog.Logger, svc event.Service) http.HandlerFunc {
	// Возвращаем функцию-обработчик
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost{
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
			log.Error("err", request.ErrEmptyReqBody, slog.Any("op", op))
			response.WriteJSON(w, http.StatusBadRequest, valResp.Error("request body is empty"))
			return
		}

		if err != nil {
			log.Error("failed to decode request body", slog.Any("op", op), sl.Err(err))
			response.WriteJSON(w, http.StatusBadRequest, valResp.Error("failed to decode request body"))
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
		e := event.Event{
			Date:  req.Date,
			Title: req.Title,
			Desc:  req.Desc,
		}
		// Добавляем событие через сервисный слой
		if err := svc.Add(e); err != nil {
			log.Error("failed to add event", slog.Any("op", op), sl.Err(err))
		}

		// Логируем успешное добавление события
		log.Info("event added", slog.Any("title", e.Title))
		// Отправляем успешный ответ клиенту
		responseOK(w, e.Title)
	}
}

// responseOK отправляет успешный ответ клиенту
// Параметры:
//   - w http.ResponseWriter: интерфейс для записи ответа
//   - title string: заголовок добавленного события
func responseOK(w http.ResponseWriter, title string) {
	r := dto.AddEventResponse{
		Title: title,
	}
	response.WriteJSON(w, http.StatusOK, r)
}
