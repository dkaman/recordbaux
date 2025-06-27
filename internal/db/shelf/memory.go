package shelf

import (
	"fmt"
	"slices"
	"sync"

	"github.com/dkaman/recordbaux/internal/db"
)

type MemoryRepo struct {
	mu      sync.Mutex
	shelves []*Entity
}

func NewMemoryRepo() *MemoryRepo {
	s := make([]*Entity, 0)
	return &MemoryRepo{
		shelves: s,
	}
}

func (r *MemoryRepo) All() ([]*Entity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]*Entity, len(r.shelves))
	copy(out, r.shelves)
	return  out, nil
}

func (r *MemoryRepo) Get(id db.ID) (*Entity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range r.shelves {
		if s.ID.String() == id.String() {
			return s, nil
		}
	}

	return &Entity{}, fmt.Errorf("shelf '%s' not found", id)
}

func (r *MemoryRepo) Save(shelf *Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, s := range r.shelves {
		if s.ID.String() == shelf.ID.String() {
			r.shelves[i] = shelf
			return nil
		}
	}

	r.shelves = append(r.shelves, shelf)

	return nil
}

func (r *MemoryRepo) Delete(id db.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, s := range r.shelves {
		if s.ID.String() == id.String() {
			r.shelves = slices.Delete(r.shelves, i, i)
			return nil
		}
	}

	return fmt.Errorf("shelf '%s' not found to delete", id)

}
