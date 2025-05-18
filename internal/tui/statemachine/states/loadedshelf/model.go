package loadedshelf

import (
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/key"

    "github.com/dkaman/recordbaux/internal/physical"
    "github.com/dkaman/recordbaux/internal/tui/statemachine"
)

type LoadShelfMsg struct {
    Shelf *physical.Shelf
}

func WithShelf(shelf *physical.Shelf) tea.Cmd {
	return func() tea.Msg {
		return LoadShelfMsg{
			Shelf: shelf,
		}
	}
}

type binKey = key.Binding

type keyMap struct {
	Next binKey
	Prev binKey
	Back binKey
}

func defaultKeys() keyMap {
	return keyMap{
		Next: key.NewBinding(key.WithKeys("right", "l")),
		Prev: key.NewBinding(key.WithKeys("left", "h")),
		Back: key.NewBinding(key.WithKeys("esc")),
	}
}

type LoadedShelfState struct {
	shelf       *physical.Shelf
	selectedBin int
	keys        keyMap
	nextState   statemachine.StateType
}

// New constructs a LoadedShelfState ready to receive a LoadShelfMsg
func New() LoadedShelfState {
	return LoadedShelfState{
		keys:      defaultKeys(),
		nextState: statemachine.LoadedShelf,
	}
}

func (s LoadedShelfState) Init() tea.Cmd {
	return nil
}

func (s LoadedShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoadShelfMsg:
		s.shelf = msg.Shelf
		s.selectedBin = 0
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Next):
			if s.shelf != nil {
				s.selectedBin = (s.selectedBin + 1) % len(s.shelf.Bins)
			}
		case key.Matches(msg, s.keys.Prev):
			if s.shelf != nil {
				s.selectedBin = (s.selectedBin - 1 + len(s.shelf.Bins)) % len(s.shelf.Bins)
			}
		case key.Matches(msg, s.keys.Back):
			s.nextState = statemachine.MainMenu
		case msg.String() == "enter":
			// TODO: handle bin selection
		}
	}

	return s, nil
}

func (s LoadedShelfState) View() string {
	if s.shelf == nil {
		return "\n  No shelf loaded\n"
	}
	var parts []string
	for i := range s.shelf.Bins {
		if i == s.selectedBin {
			parts = append(parts, "[*]")
		} else {
			parts = append(parts, "[ ]")
		}
	}
	return "\n\n" + strings.Join(parts, " ") + "\n"
}

func (s LoadedShelfState) Next(_ tea.Msg) (*statemachine.StateType, error) {
	return &s.nextState, nil
}
