package event

type Storage interface {
    Add(e Event) error
    Get(uuid uint64) (Event, error)
    Update(e Event) error
    Delete(uuid uint64) error
    List() ([]Event, error)
}
