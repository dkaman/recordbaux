package loadedbin

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

func (s LoadedBinState) renderModel() string {
	canvas := lipgloss.NewCanvas()

	boxW := s.width/2
	boxH := s.height

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(boxW).
		Height(boxH)

	s.records.SetWidth(boxW-2)
	s.records.SetHeight(boxH-2)
	s.selectedRecord.SetSize(boxW-2, boxH-2)

	left := lipgloss.NewLayer(boxStyle.Render(s.records.View()))
	right := lipgloss.NewLayer(boxStyle.Render(s.selectedRecord.View()))

	canvas.AddLayers(left.
		X(0).Y(0),
	)

	canvas.AddLayers(right.
		X(boxW).Y(0),
	)

	return canvas.Render()
}
