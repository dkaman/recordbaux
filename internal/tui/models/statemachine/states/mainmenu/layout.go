package mainmenu

import (
	"github.com/dkaman/recordbaux/internal/tui/style"

	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

func (s MainMenuState) renderModel() string {
	canvas := lipgloss.NewCanvas()

	if len(s.shelves.Items()) == 0 && len(s.playlists.Items()) == 0 {
		empty := lipgloss.NewLayer(style.ActiveTextStyle.Render("no shelves or playlists defined, 'o' to create shelf..."))
		canvas.AddLayers(empty)
		return canvas.Render()
	}

	// Create styles for focused and blurred list containers
	focusedStyle := lipgloss.NewStyle().BorderForeground(style.LightGreen)
	blurredStyle := lipgloss.NewStyle().BorderForeground(style.DarkWhite)

	var shelfBoxStyle, playlistBoxStyle lipgloss.Style

	if s.focus == shelvesView {
		shelfBoxStyle = focusedStyle
		playlistBoxStyle = blurredStyle
	} else {
		shelfBoxStyle = blurredStyle
		playlistBoxStyle = focusedStyle
	}

	boxW := s.width / 2
	boxH := s.height

	s.shelves.SetSize(boxW, boxH)
	s.playlists.SetSize(boxW, boxH)

	shelfBox := lipgloss.NewLayer(shelfBoxStyle.Render(s.shelves.View()))
	playlistBox := lipgloss.NewLayer(playlistBoxStyle.Render(s.playlists.View()))

	canvas.AddLayers(shelfBox.
		Width(boxW).
		Height(boxH).
		X(0).Y(0),
	)

	canvas.AddLayers(playlistBox.
		Width(boxW).
		Height(boxH).
		X(boxW).Y(0),
	)

	return canvas.Render()
}
