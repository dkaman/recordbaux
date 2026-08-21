package cmds

import (
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db/playlist"
)

type PlaylistLoadIntentMsg struct {
	ID uint
}

type PlaylistsLoadIntentMsg struct {}

type PlaylistsLoadedMsg struct {
	Err       error
	Playlists []*playlist.Entity
}

type PlaylistSaveIntentMsg struct {
	Entity *playlist.Entity
}

type PlaylistSavedMsg struct {
	Err error
}

type PlaylistDeleteIntentMsg struct {
	ID uint
}

type PlaylistDeletedMsg struct {
	Err error
}

func GetPlaylistCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		return PlaylistLoadIntentMsg{ID: id}
	}
}

func GetAllPlaylistsCmd() tea.Cmd {
	return func() tea.Msg {
		return PlaylistsLoadIntentMsg{}
	}
}

func SavePlaylistCmd(e *playlist.Entity) tea.Cmd {
	return func() tea.Msg {
		return PlaylistSaveIntentMsg{Entity: e}
	}
}


func DeletePlaylistCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		return PlaylistDeleteIntentMsg{ID: id}
	}
}
