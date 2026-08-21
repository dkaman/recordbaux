package services

import (
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/record"
)

type recordDB db.Repository[*record.Entity]

type RecordService struct {
	records recordDB
}


func NewRecordService(repo recordDB) *RecordService {
	return &RecordService{
		records: repo,
	}
}

func (s *RecordService) GetRecord(id uint) (*record.Entity, error) {
	return s.records.Get(id)
}

func (s *RecordService) GetAllRecords() ([]*record.Entity, error) {
	return s.records.All()
}

func (s *RecordService) SaveRecord(r *record.Entity) error {
	return s.records.Save(r)
}

func (s *RecordService) DeleteRecord(id uint) error {
	return s.records.Delete(id)
}
