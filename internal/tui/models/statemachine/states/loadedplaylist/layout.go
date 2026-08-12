package loadedplaylist

import (
	tea "charm.land/bubbletea/v2"
	// lipgloss "charm.land/lipgloss/v2"
)

func (s LoadedPlaylistState) renderModel() tea.View {
	return s.playlist.View()
}
