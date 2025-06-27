package mainmenu

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	keyFmt "github.com/dkaman/recordbaux/internal/tui/key"
)

type MainMenuState struct {
	app       *app.App
	nextState states.StateType

	keys   keyMap
	layout *layout.Div

	logger  *slog.Logger
	shelves list.Model
}

type refreshMsg struct {}

func New(a *app.App, l *layout.Div, log *slog.Logger) MainMenuState {
	log = log.WithGroup("mainmenu")

	delegate := list.NewDefaultDelegate()
	delegate.Styles = style.DefaultItemStyles()

	lst := list.New([]list.Item{}, delegate, 100, 20)
	lst.Title = "select a shelf"
	lst.Styles = style.DefaultListStyles()
	items := make([]list.Item, 0)
	lst.SetItems(items)

	lay, _ := newMainMenuLayout(l, lst)

	return MainMenuState{
		app:       a,
		keys:      defaultKeybinds(),
		layout:    lay,
		shelves:   lst,
		logger:    log,
		nextState: states.Undefined,
	}
}

func (s MainMenuState) Init() tea.Cmd {
	s.logger.Info("mainmenu state init called")
	return tea.Sequence(
		tcmds.GetAllShelvesCmd(s.app.Shelves, s.logger),
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
		shlvs := s.app.AllShelves
		items := make([]list.Item, len(shlvs))
		for i, sh := range shlvs {
			items[i] = shelf.New(sh, s.logger)
		}
		s.shelves.SetItems(items)
		s.layout, _ = newMainMenuLayout(s.layout, s.shelves)

		return s, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.SelectShelf):
			if sel, ok := s.shelves.SelectedItem().(shelf.Model); ok {
				s.nextState = states.LoadedShelf
				s.app.CurrentShelf = sel.PhysicalShelf()
			}
			return s, tea.Batch(cmds...)

		case key.Matches(msg, s.keys.NewShelf):
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
