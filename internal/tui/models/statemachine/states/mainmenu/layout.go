package mainmenu

import (
	"github.com/dkaman/recordbaux/internal/tui/style"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

func (s MainMenuState) renderModel() tea.View {
	var background string

	if len(s.shelves.Items()) == 0 && len(s.playlists.Items()) == 0 {
		emptyText := style.ActiveTextStyle.Render("no shelves or playlists defined, 'o' to create shelf...")
		background = lipgloss.Place(
			s.width, s.height,
			lipgloss.Center, lipgloss.Center,
			emptyText,
		)
	} else {
		shelfW := s.shelves.Width()
		shelfH := s.shelves.Height()
		left := lipgloss.NewStyle().
			Width(shelfW).
			Height(shelfH).
			Render(s.shelves.View().Content)

		playlistW := s.playlists.Width()
		playlistH := s.playlists.Height()
		right := lipgloss.NewStyle().
			Width(playlistW).
			Height(playlistH).
			Render(s.playlists.View().Content)

		background = lipgloss.JoinHorizontal(lipgloss.Top, left, right)

		background = lipgloss.Place(
			s.width, s.height,
			lipgloss.Left, lipgloss.Top,
			background,
		)
	}

	return tea.NewView(background)
}
