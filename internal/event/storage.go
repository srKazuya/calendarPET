package event

import "time"

type Storage interface {
	Add(e Event) error
	Update(e Event) error
	Delete(uuid uint64) error
	ListByDay(t time.Time) ([]Event, error)
	ListByWeek(t time.Time) ([]Event, error)
	ListByMonth(t time.Time) ([]Event, error)
}
