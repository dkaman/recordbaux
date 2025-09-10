package tui

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/util"
)

func getHandlers() *handlers.Registry {
	r := handlers.NewRegistry()

	handlers.Register(r, handleTeaWindowSizeMsg)
	handlers.Register(r, handleTeaKeyPressMsg)

	return r
}

func handleTeaWindowSizeMsg(m Model, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd, tea.Msg) {
	var cmds []tea.Cmd

	m.width, m.height = msg.Width, msg.Height

	if !m.ready {
		m.logger.Debug("ready")
		m.ready = true
		cmds = append(cmds, m.stateMachine.Init())
	}

	numBars := 2
	if m.helpVisible {
		numBars = 3
	}

	smWidth := m.width - 2
	smHeight := m.height - numBars - 2

	sizeMsg := tea.WindowSizeMsg{
		Width:  smWidth,
		Height: smHeight,
	}

	return m, tea.Batch(cmds...), sizeMsg
}

func handleTeaKeyPressMsg(m Model, msg tea.KeyPressMsg) (tea.Model, tea.Cmd, tea.Msg) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit, nil

	case key.Matches(msg, m.keys.ToggleHelp):
		m.helpVisible = !m.helpVisible

		numBars := 2
		if m.helpVisible {
			numBars = 3
		}

		smWidth := m.width - 2
		smHeight := m.height - numBars - 2
		sizeMsg := tea.WindowSizeMsg{Width: smWidth, Height: smHeight}

		var sizeCmd tea.Cmd
		m.stateMachine, sizeCmd = util.UpdateModel(m.stateMachine, sizeMsg)

		return m, sizeCmd, nil
	}

	return m, nil, msg
}
