package services

import (
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/shelf"
)

type shelfDB db.Repository[*shelf.Entity]

type ShelfService struct {
	Shelves      shelfDB
	AllShelves  []*shelf.Entity
	CurrentShelf *shelf.Entity
	CurrentBin   *bin.Entity
}

func NewShelfService(repo shelfDB) *ShelfService {
	return &ShelfService{
		Shelves: repo,
	}
}
