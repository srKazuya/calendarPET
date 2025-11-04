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


func NewUpdateEventHandler(log *slog.Logger, svc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		const op = "handlers.event.update"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(r)),
		)

		var req dto.UpdateEventRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if errors.Is(err, io.EOF) {
			log.Error("bad request",
				slog.String("type", request.ErrEmptyReqBody.Error()),
				sl.Err(err),
			)
			updateEventResponse(w, request.ErrEmptyReqBody.Error())
			return
		}
		if err != nil {
			log.Error("bad request",
				slog.String("type", request.ErrFailedToDecodeReqBody.Error()),
				sl.Err(err),
			)
			updateEventResponse(w, request.ErrFailedToDecodeReqBody.Error())
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			response.WriteJSON(w, http.StatusBadRequest, valResp.ValidationError(validateErr))
			return
		}

		reqEvent := event.Event{
			UUID:     req.UUID,
			UserUUID: req.UserUUID,
			Date:     req.Date,
			Title:    req.Title,
			Desc:     req.Desc,
		}

		if err := svc.Update(reqEvent); err != nil {
			switch {
			case errors.Is(err, inmem.ErrNoValue):
				log.Error("failed to delete event", sl.Err(err))
				return
			default:
				log.Error("unexpected error adding event", sl.Err(err))
				return
			}
		}

		log.Info("event update", slog.Any("title", req.UUID))

		updateEventResponseOK(w, req.UUID)
	}
}

func updateEventResponseOK(w http.ResponseWriter, id uint64) {
	r := dto.UpdateEventResponse{
		ValidationResponse: valResp.OK(),
		UUID:               id,
	}
	response.WriteJSON(w, http.StatusOK, r)
}

func updateEventResponse(w http.ResponseWriter, e string) {
	r := dto.UpdateEventResponse{
		ValidationResponse: valResp.Error(e),
	}
	response.WriteJSON(w, http.StatusBadRequest, r)
}
