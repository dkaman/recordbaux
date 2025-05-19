package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"

	css "github.com/dkaman/recordbaux/internal/tui/statemachine/states/createshelf"
	lss "github.com/dkaman/recordbaux/internal/tui/statemachine/states/loadedshelf"
	mms "github.com/dkaman/recordbaux/internal/tui/statemachine/states/mainmenu"
)

type Model struct {
	// global application config
	cfg          *config.Config

	// tea models
	stateMachine statemachine.Model

	// styling/layout
	layout       *layouts.TallLayout
}

func New(c *config.Config) Model {
	tallLayout := layouts.NewTallLayout()

	tallLayout.WithSection(layouts.StatusLine, "state: main menu")

	// Initialize state machine, passing bg into CreateShelfState
	sm, _ := statemachine.New(
		// our initial state is main menu
		statemachine.MainMenu,

		map[statemachine.StateType]statemachine.State{
			statemachine.MainMenu:    mms.New(tallLayout),
			statemachine.CreateShelf: css.New(tallLayout),
			statemachine.LoadedShelf: lss.New(tallLayout),
		},

		tallLayout,
	)

	return Model{
		cfg:          c,
		stateMachine: sm,
		layout:       tallLayout,
	}
}

func (m Model) Init() tea.Cmd {
	return m.stateMachine.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	stateMachineModel, stateMachineCmds := m.stateMachine.Update(msg)
	if sm, ok := stateMachineModel.(statemachine.Model); ok {
		m.stateMachine = sm
	}

	cmds = append(cmds, stateMachineCmds)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	_ = m.stateMachine.View()
	return m.layout.Render()
}
