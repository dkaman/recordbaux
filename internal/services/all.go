package services

import "log/slog"

type AllServices struct {
	*ShelfService
	*BinService
	*RecordService
	*PlaylistService
	*TrackService
}

func New(log *slog.Logger, s shelfDB, b binDB, r recordDB, p playlistDB, t trackDB) *AllServices {
	return &AllServices{
		NewShelfService(s),
		NewBinService(b),
		NewRecordService(r),
		NewPlaylistService(p),
		NewTrackService(t),
	}
}
