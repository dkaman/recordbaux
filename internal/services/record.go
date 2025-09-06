package services

import (
	"fmt"
	"log/slog"

	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/record"
)

type recordDB db.Repository[*record.Entity]

type RecordService struct {
	logger  *slog.Logger
	records recordDB
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
