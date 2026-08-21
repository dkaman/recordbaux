package cmds

import (
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/db/track"
)

type ShelfLoadIntentMsg struct {
	ID uint
}

type ShelvesLoadIntentMsg struct {}

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

type ShelfAllTracksMsg struct {
	ID     uint
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
