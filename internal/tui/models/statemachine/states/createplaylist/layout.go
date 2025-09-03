package createplaylist
import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)
func (s CreatePlaylistState) renderModel() string {
	canvas := lipgloss.NewCanvas()
	tracksLayer := lipgloss.NewLayer(s.list.View())
	canvas.AddLayers(tracksLayer.
		Width(s.width).
		Height(s.height).
		X(0).Y(0),
	)

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

		canvas.AddLayers(formLayer.
			X(formX).Y(formY).Z(1),
		)
	}

	return canvas.Render()
}
