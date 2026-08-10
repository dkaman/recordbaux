package loadedplaylist

import (
	lipgloss "charm.land/lipgloss/v2"
)

func (s LoadedPlaylistState) renderModel() string {
	s.trackTable.SetWidth(s.width)
	s.trackTable.SetHeight(s.height)

	layers := []*lipgloss.Layer{
		lipgloss.NewLayer(s.trackTable.View()).X(0).Y(0),
	}

	comp := lipgloss.NewCompositor(layers...)

	return comp.Render()
}
