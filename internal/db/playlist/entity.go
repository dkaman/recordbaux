package playlist

import "github.com/dkaman/recordbaux/internal/db/track"

type Entity struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Tracks []*track.Entity `gorm:"many2many:playlist_tracks;"`
}

// TableName overrides the default table name to be `playlists`.
func (e *Entity) TableName() string {
	return "playlists"
}
