package statemachine

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"
)

type StateTransitionMsg struct {
	NextState StateType
}

func WithNextState(t StateType) tea.Cmd {
	return func() tea.Msg {
		return StateTransitionMsg{
			NextState: t,
		}
	}
}

type LoadShelfMsg struct {
	Shelf *physical.Shelf
}

type NewShelfMsg struct {
	Shelf *physical.Shelf
}

func WithLoadShelf(shelf *physical.Shelf) tea.Cmd {
	return func() tea.Msg {
		return LoadShelfMsg{
			Shelf: shelf,
		}
	}
}

func WithNewShelf(shelf *physical.Shelf) tea.Cmd {
	return func() tea.Msg {
		return NewShelfMsg{
			Shelf: shelf,
		}
	}
}
