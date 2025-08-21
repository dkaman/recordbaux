package services

import (
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/track"
)

type trackDB db.Repository[*track.Entity]

type TrackService struct {
	Tracks    trackDB
	AllTracks []*track.Entity
}

func NewTrackService(repo trackDB) *TrackService {
	return &TrackService{
		Tracks: repo,
	}
}
