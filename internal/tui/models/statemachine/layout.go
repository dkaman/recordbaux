package statemachine

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

func (m Model) renderModel() tea.View {
	currentState := m.currentState.View()

	viewportStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height)

	content := viewportStyle.Render(currentState.Content)

	return tea.NewView(content)
}
