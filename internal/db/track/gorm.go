package track

import (
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

// NewPostgresRepo opens a connection and migrates the schema based on your Entity definitions.
func NewRepo(db *gorm.DB) (*Repo, error) {
	return &Repo{db: db}, nil
}

// All returns every shelf.Entity stored in the database.
func (r *Repo) All() ([]*Entity, error) {
	var tracks []*Entity
	err := r.db.
		Find(&tracks).Error

	if err != nil {
		return nil, err
	}

	return tracks, nil
}

// Get looks up a shelf.Entity by its ID.
func (r *Repo) Get(id uint) (*Entity, error) {
	var t Entity
	err := r.db.
		First(&t, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Save creates or updates a shelf.Entity and its nested bins.
func (r *Repo) Save(e *Entity) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
}

// Delete removes a shelf.Entity by ID.
func (r *Repo) Delete(id uint) error {
	return r.db.Delete(&Entity{}, "id = ?", id).Error
}
