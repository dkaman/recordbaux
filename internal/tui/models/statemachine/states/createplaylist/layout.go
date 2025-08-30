package createplaylist

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)


func (s CreatePlaylistState) renderModel() string {
	canvas := lipgloss.NewCanvas()

	var tracksLayer *lipgloss.Layer

	if len(s.table.Rows()) == 0 {
		tracksLayer = lipgloss.NewLayer("no tracks defined, load some into a shelf...")
	} else {
		s.table.SetWidth(s.width)
		s.table.SetHeight(s.height)
		tracksLayer = lipgloss.NewLayer(s.table.View())
	}

	canvas.AddLayers(tracksLayer.
		Width(s.width).
		Height(s.height).
		X(0).Y(0),
	)

	return canvas.Render()
}
