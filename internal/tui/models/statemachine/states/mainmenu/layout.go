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

	boxW := s.width / 2
	boxH := s.height

	// Create styles for focused and blurred list containers
	focusedStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(style.LightGreen).
		Width(boxW).
		Height(boxH)

	blurredStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(style.DarkWhite).
		Width(boxW).
		Height(boxH)

	s.shelves.SetSize(boxW-2, boxH-2)
	s.playlists.SetSize(boxW-2, boxH-2)

	var shelfBoxStyle, playlistBoxStyle lipgloss.Style

	if s.focus == shelvesView {
		shelfBoxStyle = focusedStyle
		playlistBoxStyle = blurredStyle
	} else {
		shelfBoxStyle = blurredStyle
		playlistBoxStyle = focusedStyle
	}

	shelfBox := lipgloss.NewLayer(shelfBoxStyle.Render(s.shelves.View()))
	playlistBox := lipgloss.NewLayer(playlistBoxStyle.Render(s.playlists.View()))

	canvas.AddLayers(shelfBox.
		X(0).Y(0),
	)

	canvas.AddLayers(playlistBox.
		X(boxW).Y(0),
	)

	return canvas.Render()
}
