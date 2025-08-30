package loadedshelf

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/key"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

	keyFmt "github.com/dkaman/recordbaux/internal/tui/key"
)

type refreshMsg struct{}

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
	return s.refresh()
}

func (s LoadedShelfState) refresh() tea.Cmd {
	return func() tea.Msg {
		return refreshMsg{}
	}
}

func (s LoadedShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height

	case refreshMsg:
		sh := s.shelfService.CurrentShelf

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
			s.nextState = states.MainMenu

		case key.Matches(msg, s.keys.Load):
			s.logger.Debug("loading colleciton",
				slog.Any("shelf", s.shelf.PhysicalShelf()),
			)

			s.nextState = states.LoadCollection

		case msg.String() == "enter":
			s.logger.Debug("bin selected",
				slog.Any("bin", s.shelf.GetSelectedBin().PhysicalBin()),
			)

			s.shelfService.CurrentBin = s.shelf.GetSelectedBin().PhysicalBin()
			s.nextState = states.LoadedBin
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
