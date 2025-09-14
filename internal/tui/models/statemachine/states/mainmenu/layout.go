package mainmenu

import (
	"log/slog"

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

	if s.creating {
		s.shelves.Styles = style.DefaultListStylesDimmed()
		s.playlists.Styles = style.DefaultListStylesDimmed()

		focusedStyle = focusedStyle.BorderForeground(style.LightGreenDimmed)
		blurredStyle = blurredStyle.BorderForeground(style.DarkWhiteDimmed)
	} else {
		s.shelves.Styles = style.DefaultListStyles()
		s.playlists.Styles = style.DefaultListStyles()
	}

	s.shelves.SetSize(boxW-2, boxH-2)
	s.playlists.SetSize(boxW-2, boxH-2)

	s.logger.Debug("box dims",
		slog.Any("width", s.width),
		slog.Any("height", s.height),
	)

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

	if s.creating {
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
			X(formX).Y(formY).Z(1),
		)
	}

	return canvas.Render()
}
