package services

import (
	"fmt"

	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/record"
)

type recordDB db.Repository[*record.Entity]

type RecordService struct {
	Records      recordDB
}

func NewRecordService(repo recordDB) *RecordService {
	return &RecordService{
		Records: repo,
	}
}

func (s *RecordService) UpdateCheckedOutStatus(records []*record.Entity, status bool) error {
	for _, r := range records {
		r.CheckedOut = status

		err := s.Records.Save(r)
		if err != nil {
			return fmt.Errorf("error checking out record: %w", err)
		}
	}

	return nil
}
