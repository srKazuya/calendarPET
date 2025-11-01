// Package handlerdto provides
package handlerdto

import (
	resp "calendar/pkg/validator"
	"time"
)

type AddEventRequest struct {
	Date  time.Time `json:"date" validate:"required"`
	Title string    `json:"title" validate:"required"`
	Desc  string    `json:"desc" validate:"required"`
}

type AddEventResponse struct {
	resp.ValidationResponse
	Title string `json:"title"`
}
