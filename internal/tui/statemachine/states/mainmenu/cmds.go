package main_menu_state

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"
)

type NewShelfMsg struct {
	Shelf *physical.Shelf
}

func  WithShelf(shelf *physical.Shelf) tea.Cmd {
	return func() tea.Msg {
		return NewShelfMsg{
			Shelf: shelf,
		}
	}
}
