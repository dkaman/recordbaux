package loadedplaylist

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

func (s PlaylistLoadedState) renderModel() string {
	canvas := lipgloss.NewCanvas()

	s.trackTable.SetWidth(s.width)
	s.trackTable.SetHeight(s.height)

	tracks := lipgloss.NewLayer(s.trackTable.View())

	canvas.AddLayers(tracks.
		X(0).Y(0),
	)

	return canvas.Render()
}
