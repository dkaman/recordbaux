package mainmenu

import (
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
)

type MainMenuState struct {
	app       *app.App
	keys      keyMap
	nextState statemachine.StateType
}

func New(a *app.App) MainMenuState {
	return MainMenuState{
		app:       a,
		nextState: statemachine.Undefined,
		keys:      defaultKeybinds(),
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
			s.nextState = statemachine.SelectShelf
		case key.Matches(msg, s.keys.NewShelf):
			s.nextState = statemachine.CreateShelf
		}
	}

	return s, tea.Batch(cmds...)
}

func (s MainMenuState) View() string {
	return ""
}

func (s MainMenuState) Next() (statemachine.StateType, bool) {
	if s.nextState != statemachine.Undefined {
		return s.nextState, true
	}

	return statemachine.Undefined, false
}

func (s MainMenuState) Transition() {
	s.nextState = statemachine.Undefined
}
