package tui

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/tui/style"
)

func (m Model) renderModel() tea.View {
	if !m.ready {
		return tea.NewView("\n initializing...")
	}

	numBars := 2
	if m.helpVisible {
		numBars = 3
	}

	barStyle := style.BarStyle.
		Width(m.width).
		Height(1)

	helpStyle := style.HelpBarStyle.
		Width(m.width).
		Height(1)

	viewportStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width).
		Height(max(0, m.height-numBars))

	topBar := barStyle.Render(m.topBarText)
	statusBar := barStyle.Render(m.statusBarText)
	viewPort := viewportStyle.Render(m.stateMachine.View().Content)

	var content string

	if m.helpVisible {
		helpBar := helpStyle.Render(m.Help())
		content = lipgloss.JoinVertical(lipgloss.Left,
			topBar,
			viewPort,
			helpBar,
			statusBar,
		)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Left,
			topBar,
			viewPort,
			statusBar,
		)
	}

	return tea.NewView(content)
}
