package services

import (
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/shelf"
)

type shelfDB db.Repository[*shelf.Entity]

type ShelfService struct {
	shelves shelfDB
}

func NewShelfService(repo shelfDB) *ShelfService {
	return &ShelfService{
		shelves: repo,
	}
}

func (s *ShelfService) GetAllShelves() ([]*shelf.Entity, error) {
	return s.shelves.All()
}

func (s *ShelfService) GetShelf(id uint) (*shelf.Entity, error) {
	return s.shelves.Get(id)
}

func (s *ShelfService) SaveShelf(e *shelf.Entity) error {
	return s.shelves.Save(e)
}

func (s *ShelfService) DeleteShelf(id uint) error {
	return s.shelves.Delete(id)
}
