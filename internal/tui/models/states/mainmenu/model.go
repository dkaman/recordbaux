package mainmenu

import (
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
)

type MainMenuState struct {
	keys keyMap
}

func New() MainMenuState {
	return MainMenuState{
		keys: defaultKeybinds(),
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
			cmds = append(cmds,
				statemachine.WithNextState(statemachine.SelectShelf),
			)
		case key.Matches(msg, s.keys.NewShelf):
			cmds = append(cmds,
				statemachine.WithNextState(statemachine.CreateShelf),
			)
		}
	}

	return s, tea.Batch(cmds...)
}

func (s MainMenuState) View() string {
	return ""
}
