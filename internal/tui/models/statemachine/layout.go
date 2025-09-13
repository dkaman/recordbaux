package statemachine

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

func (m Model) renderModel() string {
	canvas := lipgloss.NewCanvas()

	viewportStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height)

	viewPort := lipgloss.NewLayer(viewportStyle.Render(m.currentState.View()))

	canvas.AddLayers(viewPort.
		X(0).Y(0),
	)

	return canvas.Render()
}
