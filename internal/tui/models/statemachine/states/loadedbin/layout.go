package loadedbin

import (
	lipgloss "charm.land/lipgloss/v2"
)

func (s LoadedBinState) renderModel() string {
	boxW := s.width/2
	boxH := s.height

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(boxW).
		Height(boxH)

	s.records.SetWidth(boxW-2)
	s.records.SetHeight(boxH-2)
	s.selectedRecord = s.selectedRecord.SetSize(boxW-2, boxH-2)

	left := lipgloss.NewLayer(boxStyle.Render(s.records.View()))
	right := lipgloss.NewLayer(boxStyle.Render(s.selectedRecord.View().Content))

	layers := []*lipgloss.Layer{
		left.X(0).Y(0),
		right.X(boxW).Y(0),
	}

	comp := lipgloss.NewCompositor(layers...)

	return comp.Render()
}
