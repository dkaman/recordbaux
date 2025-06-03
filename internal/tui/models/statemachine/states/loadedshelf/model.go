package loadedshelf

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

type refreshLoadedShelfMsg struct{}

type LoadedShelfState struct {
	app       *app.App
	keys      keyMap
	help      help.Model
	nextState states.StateType
	layout    *div.Div

	shelf       shelf.Model
	selectedBin int
}

// New constructs a LoadedShelfState ready to receive a LoadShelfMsg
func New(a *app.App, l *div.Div) LoadedShelfState {
	return LoadedShelfState{
		app:       a,
		keys:      defaultKeybinds(),
		help:      help.New(),
		nextState: states.Undefined,
		layout:    l,
	}
}

func (s LoadedShelfState) Init() tea.Cmd {
	return func() tea.Msg {
		return refreshLoadedShelfMsg{}
	}
}

func (s LoadedShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case refreshLoadedShelfMsg:
		s.shelf = s.app.CurrentShelf.SelectBin(0)

		s.layout, _ = newSelectShelfLayout(s.layout, s.shelf)

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
			s.nextState = states.LoadCollection

		case msg.String() == "enter":
			s.app.CurrentBin = s.shelf.GetSelectedBin()
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
	sh := s.shelf.PhysicalShelf()

	if sh == nil {
		return "no shelf loaded"
	}

	return s.layout.Render()
}

func (s LoadedShelfState) Help() string {
	return s.help.View(s.keys)
}

func (s LoadedShelfState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}

	return states.Undefined, false
}

func (s LoadedShelfState) Transition() {
	s.nextState = states.Undefined
}
