package shelf

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
	var shelves []*Entity
	err := r.db.
		Preload("Bins").
		Preload("Bins.Records", func(db *gorm.DB) *gorm.DB {
			return db.Order("coordinate ASC")
		}).
		Preload("Bins.Records.Tracklist").
		Find(&shelves).Error

	if err != nil {
		return nil, err
	}

	return shelves, nil
}

// Get looks up a shelf.Entity by its ID.
func (r *Repo) Get(id uint) (*Entity, error) {
	var s Entity
	err := r.db.
		Preload("Bins").
		Preload("Bins.Records", func(db *gorm.DB) *gorm.DB {
			return db.Order("coordinate ASC")
		}).
		Preload("Bins.Records.Tracklist").
		First(&s, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Save creates or updates a shelf.Entity and its nested bins.
func (r *Repo) Save(e *Entity) error {
	// FullSaveAssociations handles nested bin slices
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
}

// Delete removes a shelf.Entity by ID.
func (r *Repo) Delete(id uint) error {
	return r.db.Delete(&Entity{}, "id = ?", id).Error
}
