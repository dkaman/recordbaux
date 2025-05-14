package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"

	sm "github.com/dkaman/recordbaux/internal/tui/statemachine"
	css "github.com/dkaman/recordbaux/internal/tui/statemachine/states/createshelf"
	mms "github.com/dkaman/recordbaux/internal/tui/statemachine/states/mainmenu"
)

// Model holds the application state
type Model struct {
	cfg          *config.Config
	stateMachine sm.Machine
}

// New initializes the TUI model
func New(c *config.Config) Model {
	sm, _ := sm.NewMachine(sm.MainMenu, map[sm.StateType]sm.State{
		sm.MainMenu:    mms.New(),
		sm.CreateShelf: css.New(),
	})

	m := Model{
		cfg:          c,
		stateMachine: sm,
	}

	return m
}

// Init is the Bubble Tea initialization command
func (m Model) Init() tea.Cmd {
	return m.stateMachine.Init()
}

// Update routes messages based on the current state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	stateMachineModel, stateMachineCmds := m.stateMachine.Update(msg)
	if sm, ok := stateMachineModel.(sm.Machine); ok {
		m.stateMachine = sm
	}

	cmds = append(cmds, stateMachineCmds)

	return m, tea.Batch(cmds...)
}

// View renders UI based on current state
func (m Model) View() string {
	return m.stateMachine.View()
}
