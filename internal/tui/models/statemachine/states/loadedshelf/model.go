package loadedshelf

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/key"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	keyFmt "github.com/dkaman/recordbaux/internal/tui/key"
)

type LoadedShelfState struct {
	shelfService *services.ShelfService
	keys         keyMap
	nextState    states.StateType

	shelf       shelf.Model
	selectedBin int

	logger *slog.Logger
	width, height int
}

// New constructs a LoadedShelfState ready to receive a LoadShelfMsg
func New(s *services.ShelfService, log *slog.Logger) LoadedShelfState {
	logGroup := log.WithGroup(states.LoadedShelf.String())

	return LoadedShelfState{
		shelfService: s,
		keys:         defaultKeybinds(),
		nextState:    states.Undefined,
		logger:       logGroup,
	}
}

func (s LoadedShelfState) Init() tea.Cmd {
	s.logger.Debug("loadedshelf state init")
	return nil
}

func (s LoadedShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height

	case shelf.LoadShelfMsg:
		sh := msg.Phy

		s.shelf = shelf.New(sh, s.logger).
			SetSize(s.width, s.height).
			SelectBin(0)

		return s, tea.Batch(cmds...)

	case tea.KeyMsg:
		sh := s.shelf.PhysicalShelf()
		if sh == nil {
			return s, nil
		}

		switch {
		case key.Matches(msg, s.keys.Next):
			s.shelf = s.shelf.SelectNextBin()

		case key.Matches(msg, s.keys.Prev):
			s.shelf = s.shelf.SelectPrevBin()

		case key.Matches(msg, s.keys.Back):
			return s, tcmds.WithNextState(states.MainMenu, nil, nil)

		case key.Matches(msg, s.keys.Load):
			return s, tcmds.WithNextState(
				states.LoadCollection,
				nil,
				[]tea.Cmd{shelf.WithPhysicalShelf(s.shelf.PhysicalShelf())},
			)

		case msg.String() == "enter":
			b := s.shelf.GetSelectedBin().PhysicalBin()
			return s, tcmds.WithNextState(
				states.LoadedBin,
				nil,
				[]tea.Cmd{bin.WithPhysicalBin(b)},
			)
		}
	}

	shelfModel, shelfCmds := s.shelf.Update(msg)
	if sh, ok := shelfModel.(shelf.Model); ok {
		s.shelf = sh
	}
	cmds = append(cmds, shelfCmds)

	return s, tea.Batch(cmds...)
}

func (s LoadedShelfState) View() string {
	return s.renderModel()
}

func (s LoadedShelfState) Help() string {
	return keyFmt.FmtKeymap(s.keys.ShortHelp())
}

func (s LoadedShelfState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}

	return states.Undefined, false
}

func (s LoadedShelfState) Transition() states.State {
	s.nextState = states.Undefined
	return s
}
