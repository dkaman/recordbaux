package services

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/playlist"
)

type playlistDB db.Repository[*playlist.Entity]

type PlaylistService struct {
	playlists playlistDB
	logger    *slog.Logger
}

type PlaylistsLoadedMsg struct {
	Playlists []*playlist.Entity
	Err       error
}

type PlaylistSavedMsg struct {
	Err error
}

type PlaylistDeletedMsg struct {
	Err error
}
type PlaylistCheckedOutMsg struct {
	Err error
}

func NewPlaylistService(repo playlistDB, log *slog.Logger) *PlaylistService {
	log.WithGroup("playlistservice")
	return &PlaylistService{
		playlists: repo,
		logger:    log,
	}
}

func (s *PlaylistService) GetAllPlaylistsCmd() tea.Cmd {
	return func() tea.Msg {
		ps, err := s.playlists.All()
		if err != nil {
			s.logger.Error("repo error",
				slog.String("error", err.Error()),
			)
			return PlaylistsLoadedMsg{Playlists: nil, Err: err}
		}

		return PlaylistsLoadedMsg{Playlists: ps, Err: err}
	}
}

func (s *PlaylistService) SavePlaylistCmd(p *playlist.Entity) tea.Cmd {
	return func() tea.Msg {
		err := s.playlists.Save(p)
		if err != nil {
			s.logger.Error("error saving playlist to repo",
				slog.String("error", err.Error()),
			)
		}
		return PlaylistSavedMsg{Err: err}
	}
}

func (s *PlaylistService) DeletePlaylistCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		err := s.playlists.Delete(id)
		if err != nil {
			s.logger.Error("error deleting playlist",
				slog.String("error", err.Error()),
				slog.Any("id", id),
			)
		}
		return PlaylistDeletedMsg{Err: err}
	}
}

// func CheckoutPlaylistCmd(repo recordDB, p *playlist.Entity, logger *slog.Logger) tea.Cmd {
// 	// l := logger.WithGroup("checkoutplaylistcmd")
// 	return func() tea.Msg {
// 		if p == nil || len(p.Tracks) == 0 {
// 			return PlaylistCheckedOutMsg{Err: fmt.Errorf("playlist has no tracks to check out")}
// 		}

// 		// Collect all unique record IDs from the playlist's tracks
// 		recordIDs := make(map[uint]struct{})
// 		for _, track := range p.Tracks {
// 			if track.RecordID != 0 {
// 				recordIDs[track.RecordID] = struct{}{}
// 			}
// 		}

// 		for id := range recordIDs {
// 			rec, err := repo.Get(id)
// 			if err != nil {
// 				return PlaylistCheckedOutMsg{Err: err}
// 			}

// 			rec.CheckedOut = true
// 			err = repo.Save(rec)
// 			if err != nil {
// 				return PlaylistCheckedOutMsg{Err: err}
// 			}

// 		}

// 		return PlaylistCheckedOutMsg{Err: nil}
// 	}
// }
