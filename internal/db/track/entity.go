package track

import (
	discogs "github.com/dkaman/discogs-golang"
)

type Entity struct {
	ID       uint `gorm:"primaryKey"`
	RecordID uint `gorm:"index"`

	Duration string
	Position string
	Title    string
	Key      string
	BPM      int
}

func New(t discogs.Track) (*Entity, error) {
	return &Entity{
		Title:    t.Title,
		Duration: t.Duration,
		Position: t.Position,
	}, nil
}

// implementing the tabler interface to change default name so it's not
// entitites
func (e *Entity) TableName() string {
	return "tracks"
}
