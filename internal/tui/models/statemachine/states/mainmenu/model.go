package mainmenu

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	keyFmt "github.com/dkaman/recordbaux/internal/tui/key"
)

type MainMenuState struct {
	shelfService *services.ShelfService
	nextState    states.StateType

	keys   keyMap
	layout *layout.Div

	logger  *slog.Logger
	shelves list.Model
}

type refreshMsg struct{}

func New(s *services.ShelfService, l *layout.Div, log *slog.Logger) MainMenuState {
	log = log.WithGroup("mainmenu")

	delegate := list.NewDefaultDelegate()
	delegate.Styles = style.DefaultItemStyles()

	lst := list.New([]list.Item{}, delegate, 0, 0)
	lst.Title = "select a shelf"
	lst.Styles = style.DefaultListStyles()

	lay, _ := newMainMenuLayout(l, lst)

	return MainMenuState{
		shelfService: s,
		keys:         defaultKeybinds(),
		layout:       lay,
		shelves:      lst,
		logger:       log,
		nextState:    states.Undefined,
	}
}

func (s MainMenuState) Init() tea.Cmd {
	s.logger.Debug("mainmenu state init",
		slog.Int("numShelves", len(s.shelves.Items())),
	)

	return tea.Sequence(
		tcmds.GetAllShelvesCmd(s.shelfService.Shelves, s.logger),
		s.refresh(),
	)
}

func (s MainMenuState) refresh() tea.Cmd {
	return func() tea.Msg {
		return refreshMsg{}
	}
}

func (s MainMenuState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case refreshMsg:
		s.logger.Debug("refreshing shelves from service")

		shlvs := s.shelfService.AllShelves
		items := make([]list.Item, len(shlvs))

		for i, sh := range shlvs {
			items[i] = shelf.New(sh, s.logger)
		}

		s.shelves.SetItems(items)
		s.layout, _ = newMainMenuLayout(s.layout, s.shelves)
		return s, nil

	case tea.KeyMsg:
		s.logger.Debug("key pressed",
			slog.String("key", msg.Type.String()),
		)

		switch {
		case key.Matches(msg, s.keys.SelectShelf):
			if sel, ok := s.shelves.SelectedItem().(shelf.Model); ok {
				s.logger.Debug("selected shelf", slog.Any("id", sel.ID()))
				s.nextState = states.LoadedShelf
				s.shelfService.CurrentShelf = sel.PhysicalShelf()
				return s, tea.Batch(cmds...)
			}

			s.logger.Warn("somehow a shelf was selected that wasn't a shelf.Model")
			return s, tea.Batch(cmds...)

		case key.Matches(msg, s.keys.NewShelf):
			s.logger.Debug("create shelf selected")
			s.nextState = states.CreateShelf
			return s, tea.Batch(cmds...)
		}
	}

	listModel, listUpdateCmds := s.shelves.Update(msg)
	cmds = append(cmds, listUpdateCmds)

	s.shelves = listModel

	s.layout, _ = newMainMenuLayout(s.layout, s.shelves)
	return s, tea.Batch(cmds...)
}

func (s MainMenuState) View() string {
	return s.layout.Render()
}

func (s MainMenuState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}

	return states.Undefined, false
}

func (s MainMenuState) Transition() states.State {
	s.nextState = states.Undefined
	return s
}

func (s MainMenuState) Help() string {
	return keyFmt.FmtKeymap(s.keys.ShortHelp())
}
