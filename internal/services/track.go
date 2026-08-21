package services

import (
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/track"
)

type trackDB db.Repository[*track.Entity]

type TrackService struct {
	tracks trackDB
}

func NewTrackService(repo trackDB) *TrackService {
	return &TrackService{
		tracks: repo,
	}
}

func (s *TrackService) GetTrack(id uint) (*track.Entity, error) {
	return s.tracks.Get(id)
}

func (s *TrackService) GetAllTracks() ([]*track.Entity, error) {
	return s.tracks.All()
}

func (s *TrackService) SaveTrack(e *track.Entity) error {
	return s.tracks.Save(e)
}

func (s *TrackService) DeleteTrack(id uint) error {
	return s.tracks.Delete(id)
}
