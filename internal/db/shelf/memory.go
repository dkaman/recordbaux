package shelf

import (
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
)

type MemoryRepo struct {
	mu      sync.Mutex
	shelves []*Entity
	nextID  uint32
}

func NewMemoryRepo() *MemoryRepo {
	s := make([]*Entity, 0)
	return &MemoryRepo{
		shelves: s,
		nextID: 1,
	}
}

func (r *MemoryRepo) All() ([]*Entity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]*Entity, len(r.shelves))
	copy(out, r.shelves)
	return out, nil
}

func (r *MemoryRepo) Get(id uint) (*Entity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range r.shelves {
		if s.ID == id {
			return s, nil
		}
	}

	return &Entity{}, fmt.Errorf("shelf '%s' not found", id)
}

func (r *MemoryRepo) Save(shelf *Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if shelf.ID == 0 {
		shelf.ID = uint(atomic.AddUint32(&r.nextID, 1))
		r.shelves = append(r.shelves, shelf)
		return nil
	}

	for i, s := range r.shelves {
		if s.ID == shelf.ID {
			r.shelves[i] = shelf
			return nil
		}
	}

	return fmt.Errorf("error updating or inserting shelf '%d'", shelf.ID)
}

func (r *MemoryRepo) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, s := range r.shelves {
		if s.ID == id {
			r.shelves = slices.Delete(r.shelves, i, i+1)
			return nil
		}
	}

	return fmt.Errorf("shelf '%s' not found to delete", id)

}
