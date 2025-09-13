package createshelf

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

func (s CreateShelfState) renderModel() string {
	canvas := lipgloss.NewCanvas()

	formView := s.createShelfForm.View()
	formW := lipgloss.Width(formView)
	formH := lipgloss.Height(formView)

	borderedForm := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(formW).
		Height(formH).
		Render(formView)

	formLayer := lipgloss.NewLayer(borderedForm)

	formX := (s.width - formW) / 2
	formY := (s.height - formH) / 2

	canvas.AddLayers(formLayer.
		X(formX).Y(formY),
	)

	return canvas.Render()
}
