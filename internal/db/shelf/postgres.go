package shelf

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/record"
)

type PostgresConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	DBName   string `koanf:"dbname"`
}

func confToDSN(c PostgresConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}

// PostgresRepo persists shelf.Entities directly via GORM.
type PostgresRepo struct {
	db *gorm.DB
}

// NewPostgresRepo opens a connection and migrates the schema based on your Entity definitions.
func NewPostgresRepo(cnf PostgresConfig) (*PostgresRepo, error) {
	dsn := confToDSN(cnf)
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}
	// AutoMigrate your domain entity and bin entity types.
	if err := gormDB.AutoMigrate(&Entity{}, &bin.Entity{}, &record.Entity{}); err != nil {
		return nil, fmt.Errorf("auto-migrate failed: %w", err)
	}
	return &PostgresRepo{db: gormDB}, nil
}

// All returns every shelf.Entity stored in the database.
func (r *PostgresRepo) All() ([]*Entity, error) {
	var shelves []*Entity
	err := r.db.
		Preload("Bins").
		Preload("Bins.Records", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Find(&shelves).Error

	if err != nil {
		return nil, err
	}

	return shelves, nil
}

// Get looks up a shelf.Entity by its ID.
func (r *PostgresRepo) Get(id uint) (*Entity, error) {
	var s Entity
	err := r.db.
		Preload("Bins").
		Preload("Bins.Records", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		First(&s, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Save creates or updates a shelf.Entity and its nested bins.
func (r *PostgresRepo) Save(e *Entity) error {
	// FullSaveAssociations handles nested bin slices
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(e).Error
}

// Delete removes a shelf.Entity by ID.
func (r *PostgresRepo) Delete(id uint) error {
	return r.db.Delete(&Entity{}, "id = ?", id).Error
}
