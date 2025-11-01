package event

import "time"

type Event struct {
	UUID     uint64    `json:"UUID"`
	UserUUID uint64    `json:"userUUID"`
	Date     time.Time `json:"date"`
	Title    string    `json:"title"`
	Desc     string    `json:"description"`
}
