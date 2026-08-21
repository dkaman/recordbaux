package services

import (
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/bin"
)

type binDB db.Repository[*bin.Entity]

type BinService struct {
	bins binDB
}

func NewBinService(repo binDB) *BinService {
	return &BinService{
		bins: repo,
	}
}

func (s *BinService) GetBin(id uint) (*bin.Entity, error) {
	return s.bins.Get(id)
}

func (s *BinService) GetAllBins() ([]*bin.Entity, error) {
	return s.bins.All()
}

func (s *BinService) SaveBin(e *bin.Entity) error {
	return s.bins.Save(e)
}

func (s *BinService) DeleteBin(id uint) error {
	return s.bins.Delete(id)
}
