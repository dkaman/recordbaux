package services

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/track"
)

type trackDB db.Repository[*track.Entity]

type TrackService struct {
	tracks trackDB
	logger *slog.Logger
}

type AllTracksLoadedMsg struct {
	Tracks []*track.Entity
	Err    error
}

func NewTrackService(repo trackDB, log *slog.Logger) *TrackService {
	logger := log.WithGroup("trackservice")
	return &TrackService{
		logger: logger,
		tracks: repo,
	}
}

func (s *TrackService) GetAllTracksCmd() tea.Cmd {
	return func() tea.Msg {
		ts, err := s.tracks.All()
		if err != nil {
			s.logger.Error("repo error", slog.String("error", err.Error()))
			return AllTracksLoadedMsg{Tracks: nil, Err: err}
		}
		return AllTracksLoadedMsg{Tracks: ts, Err: nil}
	}
}
