package services

import (
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/playlist"
)

type playlistDB db.Repository[*playlist.Entity]

type PlaylistService struct {
	playlists playlistDB
}

func NewPlaylistService(repo playlistDB) *PlaylistService {
	return &PlaylistService{
		playlists: repo,
	}
}

func (s *PlaylistService) GetPlaylist(id uint) (*playlist.Entity, error) {
	return s.playlists.Get(id)
}

func (s *PlaylistService) GetAllPlaylists() ([]*playlist.Entity, error) {
	return s.playlists.All()
}

func (s *PlaylistService) SavePlaylist(p *playlist.Entity) error {
		return s.playlists.Save(p)
}

func (s *PlaylistService) DeletePlaylist(id uint) error {
		return s.playlists.Delete(id)
}
