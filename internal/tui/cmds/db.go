package cmds

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/playlist"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/db/track"
	"github.com/dkaman/recordbaux/internal/db/record"
)

type ShelvesLoadedMsg struct {
	Shelves []*shelf.Entity
	Err     error
}

type ShelfLoadedMsg struct {
	Shelf *shelf.Entity
	Err   error
}

type ShelfSavedMsg struct {
	Err error
}

type ShelfDeletedMsg struct {
	Err error
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

type AllTracksLoadedMsg struct {
	Tracks []*track.Entity
	Err    error
}

type PlaylistCheckedOutMsg struct {
	Err error
}

type shelfDB db.Repository[*shelf.Entity]
type playlistDB db.Repository[*playlist.Entity]
type trackDB db.Repository[*track.Entity]
type recordDB db.Repository[*record.Entity]

func GetAllShelvesCmd(repo shelfDB, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("getallshelvescmd")
	return func() tea.Msg {
		ss, err := repo.All()
		if err != nil {
			l.Error("repo error",
				slog.String("error", err.Error()),
			)
			return ShelvesLoadedMsg{Shelves: nil, Err: err}
		}

		return ShelvesLoadedMsg{Shelves: ss, Err: err}
	}
}

func GetShelfCmd(repo shelfDB, id uint, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("getshelfcmd")
	return func() tea.Msg {
		s, err := repo.Get(id)
		l.Info("s",
			"s", s,
		)
		return ShelfLoadedMsg{Shelf: s, Err: err}
	}
}

func SaveShelfCmd(repo shelfDB, e *shelf.Entity, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("saveshelfcmd")
	return func() tea.Msg {
		err := repo.Save(e)
		if err != nil {
			l.Error("error saving shelf to repo",
				slog.String("error", err.Error()),
			)
		}

		return ShelfSavedMsg{Err: err}
	}
}

func DeleteShelfCmd(repo shelfDB, id uint, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("deleteshelfcmd")
	return func() tea.Msg {
		err := repo.Delete(id)
		if err != nil {
			l.Error("error deleting shelf",
				slog.String("error", err.Error()),
				slog.Any("id", id),
			)

		}
		return ShelfDeletedMsg{Err: err}
	}
}

func GetAllPlaylistsCmd(repo playlistDB, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("getallplaylistscmd")
	return func() tea.Msg {
		ps, err := repo.All()
		if err != nil {
			l.Error("repo error",
				slog.String("error", err.Error()),
			)
			return PlaylistsLoadedMsg{Playlists: nil, Err: err}
		}

		return PlaylistsLoadedMsg{Playlists: ps, Err: err}
	}
}

func SavePlaylistCmd(repo playlistDB, p *playlist.Entity, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("saveplaylistcmd")
	return func() tea.Msg {
		err := repo.Save(p)
		if err != nil {
			l.Error("error saving playlist to repo",
				slog.String("error", err.Error()),
			)
		}
		return PlaylistSavedMsg{Err: err}
	}
}

func DeletePlaylistCmd(repo playlistDB, id uint, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("deleteplaylistcmd")
	return func() tea.Msg {
		err := repo.Delete(id)
		if err != nil {
			l.Error("error deleting playlist",
				slog.String("error", err.Error()),
				slog.Any("id", id),
			)
		}
		return PlaylistDeletedMsg{Err: err}
	}
}

func GetAllTracksCmd(repo trackDB, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("getalltrackscmd")
	return func() tea.Msg {
		ts, err := repo.All()
		if err != nil {
			l.Error("repo error", slog.String("error", err.Error()))
			return AllTracksLoadedMsg{Tracks: nil, Err: err}
		}
		return AllTracksLoadedMsg{Tracks: ts, Err: nil}
	}
}

func CheckoutPlaylistCmd(repo recordDB, p *playlist.Entity, logger *slog.Logger) tea.Cmd {
	// l := logger.WithGroup("checkoutplaylistcmd")
	return func() tea.Msg {
		if p == nil || len(p.Tracks) == 0 {
			return PlaylistCheckedOutMsg{Err: fmt.Errorf("playlist has no tracks to check out")}
		}

		// Collect all unique record IDs from the playlist's tracks
		recordIDs := make(map[uint]struct{})
		for _, track := range p.Tracks {
			if track.RecordID != 0 {
				recordIDs[track.RecordID] = struct{}{}
			}
		}

		for id := range recordIDs {
			rec, err := repo.Get(id)
			if err != nil {
				return PlaylistCheckedOutMsg{Err: err}
			}

			rec.CheckedOut = true
			err = repo.Save(rec)
			if err != nil {
				return PlaylistCheckedOutMsg{Err: err}
			}

		}

		return PlaylistCheckedOutMsg{Err: nil}
	}
}
