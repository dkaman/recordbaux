package services

import "log/slog"

type AllServices struct {
	*ShelfService
	*RecordService
	*PlaylistService
	*TrackService
}

func New(log *slog.Logger, s shelfDB, r recordDB, p playlistDB, t trackDB) *AllServices {
	return &AllServices{
		NewShelfService(s, log),
		NewRecordService(r, log),
		NewPlaylistService(p, log),
		NewTrackService(t, log),
	}
}
