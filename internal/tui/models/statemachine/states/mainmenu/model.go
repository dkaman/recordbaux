package mainmenu

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

type MainMenuState struct {
	app    *app.App
	keys   keyMap
	help   help.Model
	layout *div.Div

	nextState states.StateType
}

func New(a *app.App, l *div.Div) MainMenuState {
	lay, _ := newMainMenuLayout(l)

	return MainMenuState{
		app:       a,
		keys:      defaultKeybinds(),
		help:      help.New(),
		layout:    lay,
		nextState: states.Undefined,
	}
}

func (s MainMenuState) Init() tea.Cmd {
	return nil
}

func (s MainMenuState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.SelectShelf):
			s.nextState = states.SelectShelf

		case key.Matches(msg, s.keys.NewShelf):
			s.nextState = states.CreateShelf
		}
	}

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

func (s MainMenuState) Transition() {
	s.nextState = states.Undefined
}

func (s MainMenuState) Help() string {
	return s.help.View(s.keys)
}
