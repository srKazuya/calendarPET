// Package handlerdto provides
package handlerdto

import (
	"calendar/internal/event"
	resp "calendar/pkg/validator"
	"time"
)

type UserEvent struct {
	Date  time.Time `json:"date"`
	Title string    `json:"title"`
	Desc  string    `json:"description"`
}

type AddEventRequest struct {
	Date  time.Time `json:"date" validate:"required"`
	Title string    `json:"title" validate:"required"`
	Desc  string    `json:"desc" validate:"required"`
}
type AddEventResponse struct {
	resp.ValidationResponse
	Title string `json:"title"`
}

type DeleteEventRequest struct {
	UUID uint64 `json:"UUID" validate:"required"`
}

type DeleteEventResponse struct {
	resp.ValidationResponse
	UUID uint64 `json:"UUID" validate:"required"`
}

type UpdateEventRequest struct {
	UUID     uint64    `json:"UUID" validate:"required"`
	UserUUID uint64    `json:"userUUID" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
	Title    string    `json:"title" validate:"required"`
	Desc     string    `json:"description" validate:"required"`
}

type UpdateEventResponse struct {
	resp.ValidationResponse
	UUID uint64 `json:"UUID" validate:"required"`
}

type GetEventRequest struct {
	Date time.Time `json:"date" validate:"required"`
}

type GetEventResponse struct {
	resp.ValidationResponse
	Events []UserEvent
}

func FromEvents(events []event.Event) []UserEvent {
	res := make([]UserEvent, 0, len(events))
	for _, ev := range events {
		res = append(res, UserEvent{
			Date:  ev.Date,
			Title: ev.Title,
			Desc:  ev.Desc,
		})
	}
	return res
}
