package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui/state"
	css "github.com/dkaman/recordbaux/internal/tui/state/create_shelf_state"
	mms "github.com/dkaman/recordbaux/internal/tui/state/main_menu_state"
)

// Model holds the application state
type Model struct {
	cfg          *config.Config
	stateMachine state.Machine
}

// New initializes the TUI model
func New(c *config.Config) Model {
	sm, _ := state.NewMachine(state.MainMenu, map[state.StateType]state.State{
		state.MainMenu:    mms.New(),
		state.CreateShelf: css.New(),
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
	if sm, ok := stateMachineModel.(state.Machine); ok {
		m.stateMachine = sm
	}

	cmds = append(cmds, stateMachineCmds)

	return m, tea.Batch(cmds...)
}

// View renders UI based on current state
func (m Model) View() string {
	return m.stateMachine.View()
}
