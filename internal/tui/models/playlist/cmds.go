package playlist

import (
	tea "charm.land/bubbletea/v2"
	"github.com/dkaman/recordbaux/internal/db/playlist"
)

type LoadPlaylistMsg struct {
	Phy *playlist.Entity
}

func WithPhysicalPlaylist(p *playlist.Entity) tea.Cmd {
	return func() tea.Msg { return LoadPlaylistMsg{Phy: p} }
}
