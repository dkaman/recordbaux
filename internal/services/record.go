package services

import (
	"fmt"
	"log/slog"

	tea "charm.land/bubbletea/v2"
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/playlist"
	"github.com/dkaman/recordbaux/internal/db/record"
)

type recordDB db.Repository[*record.Entity]

type RecordService struct {
	logger  *slog.Logger
	records recordDB
}

type PlaylistCheckedOutMsg struct {
	Err error
}

func NewRecordService(repo recordDB, log *slog.Logger) *RecordService {
	logger := log.WithGroup("recordservice")
	return &RecordService{
		logger:  logger,
		records: repo,
	}
}

func (s *RecordService) UpdateCheckedOutStatus(records []*record.Entity, status bool) error {
	for _, r := range records {
		r.CheckedOut = status

		err := s.records.Save(r)
		if err != nil {
			return fmt.Errorf("error checking out record: %w", err)
		}
	}

	return nil
}

func (s *RecordService) CheckoutPlaylistCmd(p *playlist.Entity) tea.Cmd {
	return func() tea.Msg {
		if p == nil || len(p.Tracks) == 0 {
			return PlaylistCheckedOutMsg{Err: fmt.Errorf("playlist has no tracks to check out")}
		}

		recordIDs := make(map[uint]struct{})
		for _, track := range p.Tracks {
			if track.RecordID != 0 {
				recordIDs[track.RecordID] = struct{}{}
			}
		}

		var recordsToCheckout []*record.Entity

		for id := range recordIDs {
			rec, err := s.records.Get(id)
			if err != nil {
				return PlaylistCheckedOutMsg{Err: fmt.Errorf("failed to fetch record %d: %w", id, err)}
			}
			if rec.CheckedOut {
				return PlaylistCheckedOutMsg{Err: fmt.Errorf("cannot checkout playlist: record '%s' is already checked out", rec.Title)}
			}
			recordsToCheckout = append(recordsToCheckout, rec)
		}

		// Execution Pass: Mark them all as checked out and save
		for _, rec := range recordsToCheckout {
			rec.CheckedOut = true
			if err := s.records.Save(rec); err != nil {
				return PlaylistCheckedOutMsg{Err: fmt.Errorf("failed to save record %d: %w", rec.ID, err)}
			}
		}

		return PlaylistCheckedOutMsg{Err: nil}
	}
}
