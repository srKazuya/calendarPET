package handlers

import (
	"calendar/internal/event"
	dto "calendar/internal/infrastructure/http/handlers/dto"
	"calendar/internal/infrastructure/http/middleware"
	"calendar/internal/infrastructure/http/response"
	inmem "calendar/internal/infrastructure/storage/in_memory"
	"calendar/pkg/sl_logger/sl"
	valResp "calendar/pkg/validator"

	"errors"
	"log/slog"
	"net/http"
	"time"
)

var (
	errMissingDateParam  = errors.New("mising date parameter")
	errInvalidDateFormat = errors.New("invalid date format")
)

func NewEventsForDayHandler(log *slog.Logger, svc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		const op = "handlers.event.getforday"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(r)),
		)

		dateS := r.URL.Query().Get("date")
		if dateS == "" {
			log.Error("bad request",
				slog.String("type", errMissingDateParam.Error()),
			)
			getEventForDayResponseErr(w, errMissingDateParam.Error())
			return
		}
		date, err := time.Parse("2006-01-02", dateS)
		if err != nil {
			log.Error("bad request",
				slog.String("type", errInvalidDateFormat.Error()),
				sl.Err(err),
			)
			getEventForDayResponseErr(w, errInvalidDateFormat.Error())
			return
		}

		events, err := svc.ListByDay(date)
		if err != nil {
			switch {
			case errors.Is(err, inmem.ErrNoValue):
				log.Error("failed to get event", sl.Err(err))
			default:
				log.Error("unexpected error adding event", sl.Err(err))
			}
		}

		log.Info("events getted")
		getEventForDayResponseOK(w, events)
	}
}

func getEventForDayResponseOK(w http.ResponseWriter, e []event.Event) {
	r := dto.GetEventResponse{
		ValidationResponse: valResp.OK(),
		Events:             dto.FromEvents(e),
	}
	response.WriteJSON(w, http.StatusOK, r)
}

func getEventForDayResponseErr(w http.ResponseWriter, e string) {
	r := dto.GetEventResponse{
		ValidationResponse: valResp.Error(e),
	}
	response.WriteJSON(w, http.StatusBadRequest, r)
}
