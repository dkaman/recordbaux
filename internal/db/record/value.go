package record

import (
	"github.com/dkaman/discogs-golang"
)

// Entity represents a physical record we are going to place on the physical
// shelf.
type Entity struct {
	// lol
	ID         uint `gorm:"primaryKey"`
	BinID      uint `gorm:"index"`
	InstanceID int
	ReleaseID  int
	MasterID   int
	Position   int `gorm:"index"`

	Title         string
	Artists       StringArray `gorm:"type:jsonb"`
	CatalogNumber string
}

func New(ri discogs.ReleaseInstance) (*Entity, error) {
	var artists []string
	for _, a := range ri.BasicInfo.Artists {
		artists = append(artists, a.Name)
	}

	return &Entity{
		InstanceID: ri.InstanceID,
		ReleaseID:  ri.ID,
		MasterID:   ri.BasicInfo.MasterID,

		Title:         ri.BasicInfo.Title,
		Artists:       artists,
		CatalogNumber: ri.BasicInfo.Labels[0].CatNo,
	}, nil
}

// implementing the tabler interface to change default name so it's not
// entitites
func (e *Entity) TableName() string {
	return "records"
}
