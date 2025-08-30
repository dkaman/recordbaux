package createshelf

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

func (s CreateShelfState) renderModel() string {
	canvas := lipgloss.NewCanvas()

	formLayer := lipgloss.NewLayer(s.createShelfForm.View())

	canvas.AddLayers(formLayer.
		Width(s.width/3).
		Height(s.height/3).
		X(s.width/3).Y(s.height/3),
	)

	return canvas.Render()
}
