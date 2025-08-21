package playlist

import (
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) (*Repo, error) {
	return &Repo{db: db}, nil
}

func (r *Repo) All() ([]*Entity, error) {
	var es []*Entity
	err := r.db.
		Preload("Tracks").
		Find(&es).Error

	if err != nil {
		return nil, err
	}

	return es, nil
}

func (r *Repo) Get(id uint) (*Entity, error) {
	var e Entity
	err := r.db.
		First(&e, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *Repo) Save(e *Entity) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
}

func (r *Repo) Delete(id uint) error {
	return r.db.Delete(&Entity{}, "id = ?", id).Error
}
