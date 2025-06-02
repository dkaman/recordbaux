package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

type Model struct {
	// global application config/state
	cfg  *config.Config
	app  *app.App
	keys keyMap
	help help.Model

	// tea models
	stateMachine statemachine.Model

	// styling/layout
	helpVisible bool
	layout      *div.Div
}

func New(c *config.Config) Model {
	h := help.New()
	h.Styles = style.DefaultHelpStyles()

	l, _ := newTUILayout()

	m := Model{
		cfg:         c,
		app:         app.NewApp(),
		keys:        defaultKeybinds(),
		help:        h,
		helpVisible: false,
		layout:      l,
	}

	sm, _ := statemachine.New(m.app, c, m.layout)

	m.stateMachine = sm

	return m
}

func (m Model) Init() tea.Cmd {
	return m.stateMachine.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.ToggleHelp):
			m.helpVisible = !m.helpVisible
		}

	case tea.WindowSizeMsg:
		w, h := msg.Width, msg.Height
		m.layout.Resize(w, h)
		return m, nil
	}

	stateMachineModel, stateMachineCmds := m.stateMachine.Update(msg)
	if sm, ok := stateMachineModel.(statemachine.Model); ok {
		m.stateMachine = sm
	}
	cmds = append(cmds, stateMachineCmds)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.layout.Render()
}
