package cmds

import (
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db/playlist"
)

type PlaylistCheckoutIntentMsg struct {
	Playlist *playlist.Entity
	Status bool
}

type PlaylistCheckoutMsg struct {
	Err    error
	Status bool
}

func SetCheckoutCmd(p *playlist.Entity, status bool) tea.Cmd {
	return func() tea.Msg {
		return PlaylistCheckoutIntentMsg{
			Playlist: p,
			Status: status,
		}
	}
}
