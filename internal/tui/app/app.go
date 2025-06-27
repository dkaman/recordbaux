package app

import (
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/shelf"
)

type shelfDB db.Repository[*shelf.Entity]

type App struct {
	Shelves      shelfDB
	AllShelves  []*shelf.Entity
	CurrentShelf *shelf.Entity
	CurrentBin   *bin.Entity
}

func NewApp(repo shelfDB) *App {
	return &App{
		Shelves: repo,
	}
}
