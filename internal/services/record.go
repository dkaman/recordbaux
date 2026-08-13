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

type PlaylistCheckoutMsg struct {
	Err error
	Status bool
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

func (s *RecordService) SetCheckoutCmd(p *playlist.Entity, c bool) tea.Cmd {
	return func() tea.Msg {
		if p == nil || len(p.Tracks) == 0 {
			return PlaylistCheckoutMsg{Err: fmt.Errorf("playlist has no tracks to check in or out")}
		}

		recordIDs := make(map[uint]struct{})
		for _, track := range p.Tracks {
			if track.RecordID != 0 {
				recordIDs[track.RecordID] = struct{}{}
			}
		}

		var records []*record.Entity

		for id := range recordIDs {
			rec, err := s.records.Get(id)
			if err != nil {
				return PlaylistCheckoutMsg{Err: fmt.Errorf("failed to fetch record %d: %w", id, err)}
			}
			if c && rec.CheckedOut {
				return PlaylistCheckoutMsg{Err: fmt.Errorf("cannot check out playlist: record '%s' is already checked out", rec.Title)}
			}
			records = append(records, rec)
		}

		for _, rec := range records {
			rec.CheckedOut = c
			if err := s.records.Save(rec); err != nil {
				return PlaylistCheckoutMsg{Err: fmt.Errorf("failed to save record %d: %w", rec.ID, err)}
			}
		}

		return PlaylistCheckoutMsg{Err: nil, Status: c}
	}
}
