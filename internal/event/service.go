//Package event provides ...
package event


type Service interface {
	Add(e Event) error
	Get(id uint64) (Event, error)
	Delete(id uint64) error
	List() ([]Event, error)
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

func (s *service) Get(id uint64) (Event, error) {
	return s.storage.Get(id)
}

func (s *service) Delete(id uint64) error {
	return s.storage.Delete(id)
}

func (s *service) List() ([]Event, error) {
	return s.storage.List()
}
