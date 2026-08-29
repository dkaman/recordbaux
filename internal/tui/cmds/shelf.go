package cmds

import (
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/db/track"
)

type ShelfLoadIntentMsg struct {
	ID uint
}

type ShelvesLoadIntentMsg struct{}

type ShelvesLoadedMsg struct {
	Shelves []*shelf.Entity
	Err     error
}

type ShelfSaveIntentMsg struct {
	Shelf *shelf.Entity
}

type ShelfSavedMsg struct {
	Err error
}

type ShelfDeleteIntentMsg struct {
	ID uint
}

type ShelfDeletedMsg struct {
	Err error
}

type ShelfAllTracksIntentMsg struct {
	ID uint
}

type ShelvesAllTracksIntentMsg struct {
	IDs []uint
}

type ShelfAllTracksMsg struct {
	ID     uint
	Tracks []*track.Entity
	Err    error
}

type ShelvesAllTracksMsg struct {
	IDs    []uint
	Tracks []*track.Entity
	Err    error
}

func GetShelfCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		return ShelfLoadIntentMsg{ID: id}
	}
}

func GetAllShelvesCmd() tea.Cmd {
	return func() tea.Msg {
		return ShelvesLoadIntentMsg{}
	}
}

func NewShelfCmd(name, shape string, dimX, dimY, binSize, numBins int) tea.Cmd {
	var e *shelf.Entity
	var err error

	if shape == "rect" {
		e, err = shelf.New(name, binSize, shelf.WithShapeRect(dimX, dimY, binSize, bin.SortAlphaByArtist))
		if err != nil {
			return func() tea.Msg {
				return ShelvesLoadedMsg {Err: err}
			}

		}

	} else if shape == "irregular" {
		e, err = shelf.New(name, binSize)
		if err != nil {
			return func() tea.Msg {
				return ShelvesLoadedMsg {Err: err}
			}
		}

		e.AddBin(numBins, bin.SortAlphaByArtist)
	}

	return func() tea.Msg {
		return ShelfSaveIntentMsg{Shelf: e}
	}
}

func SaveShelfCmd(e *shelf.Entity) tea.Cmd {
	return func() tea.Msg {
		return ShelfSaveIntentMsg{Shelf: e}
	}
}

func DeleteShelfCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		return ShelfDeleteIntentMsg{ID: id}
	}
}

func GetAllTracksFromShelfCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		return ShelfAllTracksIntentMsg{ID: id}
	}
}

func GetAllTracksFromShelvesCmd(ids []uint) tea.Cmd {
	return func() tea.Msg {
		return ShelvesAllTracksIntentMsg{IDs: ids}
	}
}
