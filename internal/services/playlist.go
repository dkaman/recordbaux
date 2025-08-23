package services

import (
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/playlist"
)

type playlistDB db.Repository[*playlist.Entity]

type PlaylistService struct {
	Playlists playlistDB
	AllPlaylists []*playlist.Entity
	CurrentPlaylist *playlist.Entity
}

func NewPlaylistService(repo playlistDB) *PlaylistService {
	return &PlaylistService{
		Playlists: repo,
	}
}
