package createplaylist

import (
	lipgloss "charm.land/lipgloss/v2"
)

func (s CreatePlaylistState) renderModel() string {
	s.list.SetWidth(s.width)
	s.list.SetHeight(s.height)

	tracksLayer := lipgloss.NewLayer(s.list.View())

	layers := []*lipgloss.Layer{
		tracksLayer.X(0).Y(0),
	}

	comp := lipgloss.NewCompositor(layers...)

	if s.namingPlaylist {
		formView := s.nameForm.View()

		formW := s.width / 8
		formH := 5

		formStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(formW).
			Height(formH)

		formLayer := lipgloss.NewLayer(formStyle.Render(formView))

		formX := (s.width - formW) / 2
		formY := (s.height - formH) / 2

		comp.AddLayers(formLayer.
			X(formX).Y(formY).Z(1),
		)
	}

	return comp.Render()
}
