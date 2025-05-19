package loadedshelf

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"
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
		Back: key.NewBinding(key.WithKeys("q")),
	}
}

type LoadedShelfState struct {
	shelf       *physical.Shelf
	selectedBin int
	keys        keyMap

	layout *layouts.TallLayout
}

// New constructs a LoadedShelfState ready to receive a LoadShelfMsg
func New(l *layouts.TallLayout) LoadedShelfState {
	return LoadedShelfState{
		keys:      defaultKeys(),
		layout:    l,
	}
}

func (s LoadedShelfState) Init() tea.Cmd {
	return nil
}

func (s LoadedShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

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
			cmds = append(cmds, statemachine.WithNextState(statemachine.MainMenu))
			s.shelf = nil
		case msg.String() == "enter":
			// TODO: handle bin selection
		}
	}

	return s, tea.Batch(cmds...)
}

func (s LoadedShelfState) View() string {
	var view string

	if s.shelf == nil {
		view = "no shelf loaded"
	} else {
		var parts []string
		for i := range s.shelf.Bins {
			if i == s.selectedBin {
				parts = append(parts, "[*]")
			} else {
				parts = append(parts, "[ ]")
			}
		}

		view = "\n" + strings.Join(parts, " ") + "\n"
	}


	s.layout.WithSection(layouts.BottomWindow, view)

	return view
}
