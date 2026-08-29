package tui

import (
	"context"
	"fmt"
	"log/slog"

	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/db/track"
	"github.com/dkaman/recordbaux/internal/db/playlist"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

// window size commands

func (m *Model) refreshWindowSize() tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: m.width, Height: m.height}
	}
}

// discogs commands

func (m *Model) getDiscogsFolders() tea.Cmd {
	return func() tea.Msg {
		folders, _ := m.discogsClient.Collection.ListFolders(context.TODO(), m.discogsUsername)
		return tcmds.GetDiscogsFoldersMsg{
			Folders: folders,
		}
	}
}

func (m *Model) retrieveDiscogsCollection(folder string) tea.Cmd {
	return func() tea.Msg {
		var recs []*record.Entity
		var msg tcmds.NewDiscogsCollectionMsg

		releaseInstances, err := m.discogsClient.Collection.GetReleasesByFolder(context.TODO(), m.discogsUsername, 0)
		if err != nil {
			msg.Err = fmt.Errorf("error getting releases from discogs: %w", err)
			return msg
		}

		m.logger.Debug("got releases from discogs",
			slog.Int("count", len(releaseInstances)),
		)

		for _, ri := range releaseInstances {
			m.logger.Debug("processing release",
				slog.Int("id", ri.ID),
			)

			r, err := record.New(ri)
			if err != nil {
				msg.Err = fmt.Errorf("error constructing record entity: %w", err)
				return msg
			}

			recs = append(recs, r)
		}

		msg.Releases = recs
		msg.Err = nil

		return msg
	}
}

func (m *Model) enrichReleaseInstance(rec *record.Entity) tea.Cmd {
	return func() tea.Msg {
		rel, err := m.discogsClient.Database.GetRelease(context.TODO(), rec.ReleaseID)
		if err != nil {
			return tcmds.NewDiscogsEnrichRecordMsg{
				Record: nil,
				Err: err,
			}
		}

		var tracks []*track.Entity
		for _, trk := range rel.Tracklist {
			t, err := track.New(trk)
			if err != nil {
				return tcmds.NewDiscogsEnrichRecordMsg{
					Record: nil,
					Err: err,
				}
			}

			tracks = append(tracks, t)
		}

		rec.Tracklist = tracks

		return tcmds.NewDiscogsEnrichRecordMsg{
			Record: rec,
			Err:    nil,
		}
	}
}

// shelf commands

func (m *Model) getShelfCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		shlf, err := m.svcs.ShelfService.GetShelf(id)
		return tcmds.ShelvesLoadedMsg{Shelves: []*shelf.Entity{shlf}, Err: err}
	}
}

func (m *Model) getAllShelvesCmd() tea.Cmd {
	return func() tea.Msg {
		allShelves, err := m.svcs.ShelfService.GetAllShelves()
		return tcmds.ShelvesLoadedMsg{Shelves: allShelves, Err: err}
	}
}

func (m *Model) saveShelfCmd(e *shelf.Entity) tea.Cmd {
	return func() tea.Msg {
		err := m.svcs.ShelfService.SaveShelf(e)
		return tcmds.ShelfSavedMsg{Err: err}
	}
}

func (m *Model) deleteShelfCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		err := m.svcs.ShelfService.DeleteShelf(id)
		return tcmds.ShelfDeletedMsg{Err: err}
	}
}

func (m *Model) getAllTracksFromShelfCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		shlf, err := m.svcs.ShelfService.GetShelf(id)
		if err != nil {
			return tcmds.ShelfAllTracksMsg{
				ID:  id,
				Err: err,
			}
		}

		return tcmds.ShelfAllTracksMsg{
			ID:     id,
			Tracks: shlf.AllTracks(),
		}
	}
}

func (m *Model) getAllTracksFromShelvesCmd(ids []uint) tea.Cmd {
	return func() tea.Msg {
		var tracks []*track.Entity

		for _, id := range ids {
			shlf, err := m.svcs.ShelfService.GetShelf(id)
			if err != nil {
				return tcmds.ShelfAllTracksMsg{
					ID:  id,
					Err: err,
				}
			}

			tracks = append(tracks, shlf.AllTracks()...)
		}

		return tcmds.ShelvesAllTracksMsg{
			IDs:     ids,
			Tracks: tracks,
		}
	}
}

// bin commands

func (m *Model) getBinCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		b, err := m.svcs.BinService.GetBin(id)
		return tcmds.BinsLoadedMsg{
			Bins: []*bin.Entity{b},
			Err:  err,
		}
	}
}

// playlist commands

func (m *Model) savePlaylistCmd(p *playlist.Entity) tea.Cmd {
	return func() tea.Msg {
		err := m.svcs.PlaylistService.SavePlaylist(p)
		return tcmds.PlaylistSavedMsg{Err: err}
	}
}

func (m *Model) getPlaylistCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		p, err := m.svcs.PlaylistService.GetPlaylist(id)
		return tcmds.PlaylistsLoadedMsg{
			Err: err,
			Playlists: []*playlist.Entity{p},
		}
	}
}

func (m *Model) getAllPlaylistsCmd() tea.Cmd {
	return func() tea.Msg {
		ps, err := m.svcs.PlaylistService.GetAllPlaylists()

		return tcmds.PlaylistsLoadedMsg{
			Err:       err,
			Playlists: ps,
		}
	}
}

func (m *Model) deletePlaylistsCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		err := m.svcs.PlaylistService.DeletePlaylist(id)
		return tcmds.PlaylistDeletedMsg{
			Err: err,
		}
	}
}

// record commands

func (m *Model) setCheckoutCmd(p *playlist.Entity, status bool) tea.Cmd {
	return func() tea.Msg {
		if p == nil || len(p.Tracks) == 0 {
			return tcmds.PlaylistCheckoutMsg{Err: fmt.Errorf("playlist has no tracks to check in or out")}
		}

		recordIDs := make(map[uint]struct{})
		for _, track := range p.Tracks {
			if track.RecordID != 0 {
				recordIDs[track.RecordID] = struct{}{}
			}
		}

		var records []*record.Entity

		for id := range recordIDs {
			rec, err := m.svcs.RecordService.GetRecord(id)
			if err != nil {
				return tcmds.PlaylistCheckoutMsg{Err: fmt.Errorf("failed to fetch record %d: %w", id, err)}
			}
			if status && rec.CheckedOut {
				return tcmds.PlaylistCheckoutMsg{Err: fmt.Errorf("cannot check out playlist: record '%s' is already checked out", rec.Title)}
			}
			records = append(records, rec)
		}

		for _, rec := range records {
			rec.CheckedOut = status
			if err := m.svcs.RecordService.SaveRecord(rec); err != nil {
				return tcmds.PlaylistCheckoutMsg{Err: fmt.Errorf("failed to save record %d: %w", rec.ID, err)}
			}
		}

		return tcmds.PlaylistCheckoutMsg{Err: nil, Status: status}
	}
}
