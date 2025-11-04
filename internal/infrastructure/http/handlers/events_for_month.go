package handlers

import (
	"calendar/internal/event"
	"calendar/internal/infrastructure/http/middleware"
	inmem "calendar/internal/infrastructure/storage/in_memory"
	"calendar/pkg/sl_logger/sl"


	"errors"
	"log/slog"
	"net/http"
	"time"
)



func NewEventsForMonthHandler(log *slog.Logger, svc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		const op = "handlers.event.getformonth"
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

		events, err := svc.ListByMonth(date)
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

