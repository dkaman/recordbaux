package app

import (
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

type App struct {
	Shelves      []shelf.Model
	CurrentShelf shelf.Model
	CurrentBin   bin.Model
}

func NewApp() *App {
	return &App{}
}

func (a *App) AddShelf(sh shelf.Model) {
	a.Shelves = append(a.Shelves, sh)
}

func (a *App) SelectShelf(sh shelf.Model) {
	a.CurrentShelf = sh

	phy := a.CurrentShelf.PhysicalShelf()

	if phy != nil {
		b := phy.Bins[0]
		a.SelectBin(bin.New(b, style.ActiveTextStyle))
	}
}

func (a *App) SelectBin(b bin.Model) {
	a.CurrentBin = b
}
