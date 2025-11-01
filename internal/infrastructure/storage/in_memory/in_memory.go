// Package inmem provides
package inmem

import (
	"calendar/internal/event"
	"errors"
	"fmt"
	"sync"
)

var (
	errNoValue = errors.New("нет значения: ")
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
func (s *Storage) Get(id uint64) (event.Event, error) {
	const op = "infra.storage.in_memory.get"
	s.mu.Lock()
	defer s.mu.Unlock()
	if event, ok := s.db[id]; !ok {
		return event, fmt.Errorf("%s: error: %w, %v", op, errNoValue, id)
	} else {
		return event, nil
	}

}
func (s *Storage) Delete(id uint64) error {
	const op = "infra.storage.in_memory.delete"
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.db[id]; !ok {
		return fmt.Errorf("%s: error: %w, %v", op, errNoValue, id)
	}

	delete(s.db, id)
	return nil
}

func (s *Storage) Update(e event.Event) error {
	const op = "infra.storage.in_memory.update"
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.db[e.UUID]; !ok {
		return fmt.Errorf("%s: error: %w, %v", op, errNoValue, e.UUID)
	} else {
		s.db[e.UUID] = e
	}
	return nil
}

func (s *Storage) List() ([]event.Event, error) {
	const op = "infra.storage.in_memory.list"
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]event.Event, 0, len(s.db))
	for _, event := range s.db {
		result = append(result, event)
	}

	return result, nil
}
