// Package inmem provides
package inmem

import (
	"calendar/internal/event"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNoValue = errors.New("no value")
)

type Storage struct {
	mu     sync.Mutex
	db     map[uint64]event.Event
	lastID uint64
}

func New() *Storage {
	db := make(map[uint64]event.Event)
	return &Storage{db: db}
}

func (s *Storage) Add(e event.Event) error {
	const op = "infra.storage.in_memory.save"
	s.mu.Lock()
	defer s.mu.Unlock()
	if e.UUID == 0 {
		s.lastID++
		e.UUID = s.lastID
	}
	s.db[e.UUID] = e

	return nil
}

func (s *Storage) Update(e event.Event) error {
	const op = "infra.storage.in_memory.update"
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.db[e.UUID]; !ok {
		return fmt.Errorf("%s: error: %w, %v", op, ErrNoValue, e.UUID)
	} else {
		s.db[e.UUID] = e
	}
	return nil
}

func (s *Storage) Delete(id uint64) error {
	const op = "infra.storage.in_memory.delete"
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.db[id]; !ok {
		return fmt.Errorf("%s: error: %w, %v", op, ErrNoValue, id)
	}

	delete(s.db, id)
	return nil
}

func (s *Storage) ListByDay(t time.Time) ([]event.Event, error) {
	const op = "infra.storage.in_memory.list_by_day"
	result := []event.Event{}
	if len(s.db) == 0 {
		return nil, fmt.Errorf("%s: error: %w", op, ErrNoValue)
	}
	y, m, d := t.Date()
	for _, event := range s.db {
		y1, m1, d1 := event.Date.Date()
		if y == y1 && m == m1 && d == d1 {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) ListByWeek(t time.Time) ([]event.Event, error) {
	const op = "infra.storage.in_memory.list_by_week"
	result := []event.Event{}
	if len(s.db) == 0 {
		return nil, fmt.Errorf("%s: error: %w", op, ErrNoValue)
	}
	y, w := t.ISOWeek()
	for _, event := range s.db {
		y1, w1 := event.Date.ISOWeek()
		if y == y1 && w == w1 {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) ListByMonth(t time.Time) ([]event.Event, error) {
	const op = "infra.storage.in_memory.list_by_month"
	result := []event.Event{}
	if len(s.db) == 0 {
		return nil, fmt.Errorf("%s: error: %w", op, ErrNoValue)
	}
	y, m, _ := t.Date()
	for _, event := range s.db {
		y1, m1, _ := event.Date.Date()
		if y == y1 && m == m1 {
			result = append(result, event)
		}
	}
	return result, nil
}
