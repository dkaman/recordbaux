package loadedshelf

import (
	"github.com/charmbracelet/bubbles/v2/key"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

func getHandlers() *handlers.Registry {
	r := handlers.NewRegistry()

	handlers.Register(r, handleTeaWindowSizeMsg)
	handlers.Register(r, handleTeaKeyPressMsg)
	handlers.Register(r, handleLoadShelfMsg)

	return r
}

func handleTeaWindowSizeMsg(s LoadedShelfState, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.width, s.height = msg.Width, msg.Height
	return s, nil, nil
}

func handleTeaKeyPressMsg(s LoadedShelfState, msg tea.KeyPressMsg) (tea.Model, tea.Cmd, tea.Msg) {
	sh := s.shelf.PhysicalShelf()
	if sh == nil {
		return s, nil, nil
	}

	switch {
	case key.Matches(msg, s.keys.Next):
		s.shelf = s.shelf.SelectNextBin()

	case key.Matches(msg, s.keys.Prev):
		s.shelf = s.shelf.SelectPrevBin()

	case key.Matches(msg, s.keys.Back):
		return s, tcmds.Transition(states.MainMenu, nil, nil), nil

	case key.Matches(msg, s.keys.Load):
		return s, tcmds.Transition(
			states.LoadCollection,
			nil,
			[]tea.Cmd{shelf.WithPhysicalShelf(s.shelf.PhysicalShelf())},
		), nil

	case msg.String() == "enter":
		b := s.shelf.GetSelectedBin().PhysicalBin()
		return s, tcmds.Transition(
			states.LoadedBin,
			nil,
			[]tea.Cmd{bin.WithPhysicalBin(b)},
		), nil
	}

	return s, nil, msg
}

func handleLoadShelfMsg(s LoadedShelfState, msg shelf.LoadShelfMsg) (tea.Model, tea.Cmd, tea.Msg) {
	sh := msg.Phy

	s.shelf = shelf.New(sh, s.logger).
		SetSize(s.width, s.height).
		SelectBin(0)

	return s, nil, nil
}
