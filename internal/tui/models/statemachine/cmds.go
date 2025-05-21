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

type BroadcastLoadShelfMsg struct {
	Shelf *physical.Shelf
}

type NewShelfMsg struct {
	Shelf *physical.Shelf
}

type LoadBinMsg struct {
	Bin *physical.Bin
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

func WithLoadShelfBroadcast(shelf *physical.Shelf) tea.Cmd {
	return func() tea.Msg {
		return BroadcastLoadShelfMsg{
			Shelf: shelf,
		}
	}
}

func WithLoadBin(bin *physical.Bin) tea.Cmd {
	return func() tea.Msg {
		return LoadBinMsg {
			Bin: bin,
		}
	}
}
