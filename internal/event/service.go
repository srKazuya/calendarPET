//Package event provides ...
package event

import "time"

type Service interface {
	Add(e Event) error
	Update(e Event) error
	Delete(uuid uint64) error
	ListByDay(t time.Time) ([]Event, error)
	ListByWeek(t time.Time) ([]Event, error)
	ListByMonth(t time.Time) ([]Event, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) Add(e Event) error {
	return s.storage.Add(e)
}

func (s *service) Update(e Event) error {
	return s.storage.Update(e)
}

func (s *service) Delete(id uint64) error {
	return s.storage.Delete(id)
}

func (s *service) ListByDay(t time.Time) ([]Event, error) {
	return s.storage.ListByDay(t)
}

func (s *service) ListByWeek(t time.Time) ([]Event, error) {
	return s.storage.ListByWeek(t)
}

func (s *service) ListByMonth(t time.Time) ([]Event, error) {
	return s.storage.ListByMonth(t)
}
