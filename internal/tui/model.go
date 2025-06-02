package tui

import (
	"fmt"

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
	vp := l.Find("viewport")

	a := app.NewApp()
	sm, _ := statemachine.New(a, c, vp)

	m := Model{
		cfg:          c,
		app:          a,
		keys:         defaultKeybinds(),
		help:         h,
		helpVisible:  false,
		layout:       l,
		stateMachine: sm,
	}

	_ = addTopBarText(l, "recordbaux - organize your record collection")
	_ = addStatusBarText(m.layout, fmt.Sprintf("current state: %s", m.stateMachine.CurrentStateType()))

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

			helpBar := m.layout.Find("helpbar")
			if m.helpVisible {
				helpBar.Show()
			} else {
				helpBar.Hide()
			}

			w, h := m.layout.Width(), m.layout.Height()
			m.layout.Resize(w, h)
		}

	case tea.WindowSizeMsg:
		w, h := msg.Width, msg.Height
		m.layout.Resize(w, h)
		return m, nil
	}

	// update state machine
	stateMachineModel, stateMachineCmds := m.stateMachine.Update(msg)
	if sm, ok := stateMachineModel.(statemachine.Model); ok {
		m.stateMachine = sm
	}
	cmds = append(cmds, stateMachineCmds)

	// update bars
	statusBarText := fmt.Sprintf("current state: %s", m.stateMachine.CurrentStateType())
	_ = addStatusBarText(m.layout, statusBarText)

	helpText := fmt.Sprintf("global: %s statemachine: %s", m.help.View(m.keys), m.stateMachine.Help())
	_ = addHelpBarText(m.layout, helpText)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.layout.Render()
}
