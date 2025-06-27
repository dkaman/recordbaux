package loadedshelf

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	keyFmt "github.com/dkaman/recordbaux/internal/tui/key"
)

type refreshMsg struct{}

type LoadedShelfState struct {
	shelfService *services.ShelfService
	keys         keyMap
	nextState    states.StateType
	layout       *layout.Div

	shelf       shelf.Model
	selectedBin int

	logger *slog.Logger
}

// New constructs a LoadedShelfState ready to receive a LoadShelfMsg
func New(s *services.ShelfService, l *layout.Div, log *slog.Logger) LoadedShelfState {
	logGroup := log.WithGroup(states.LoadedShelf.String())
	return LoadedShelfState{
		shelfService: s,
		keys:         defaultKeybinds(),
		nextState:    states.Undefined,
		layout:       l,
		logger:       logGroup,
	}
}

func (s LoadedShelfState) Init() tea.Cmd {
	s.logger.Info("loadedshelf state init called")
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
	case refreshMsg:
		contentWidth := s.layout.Width() - 2
		contentHeight := s.layout.Height() - 2

		sh := s.shelfService.CurrentShelf

		s.shelf = shelf.New(sh, s.logger).
			SelectBin(0).
			SetSize(contentWidth, contentHeight)

		s.layout, _ = newSelectShelfLayout(s.layout, s.shelf)

		return s, tea.Batch(cmds...)

	case tea.WindowSizeMsg:
		s.layout.Resize(msg.Width, msg.Height)
		msg.Width = msg.Width - 2
		msg.Height = msg.Height - 2

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
			s.logger.Info("loading colleciton",
				slog.Any("shelf", s.shelf.PhysicalShelf()),
			)
			s.nextState = states.LoadCollection

		case msg.String() == "enter":
			s.logger.Info("bin selected",
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

	s.layout, _ = newSelectShelfLayout(s.layout, s.shelf)

	return s, tea.Batch(cmds...)
}

func (s LoadedShelfState) View() string {
	return ""
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
