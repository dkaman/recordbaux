package loadedshelf

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

func (s LoadedShelfState) renderModel() string {
	canvas := lipgloss.NewCanvas()

	shelf := lipgloss.NewLayer(s.shelf.View())

	canvas.AddLayers(shelf.
		X(0).Y(0),
	)

	return canvas.Render()
}
