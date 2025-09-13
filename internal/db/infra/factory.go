package infra

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/playlist"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/db/track"
)

func New(conf Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Select the appropriate driver based on the configuration
	switch conf.Driver {
	case DriverPostgres:
		pgConf := conf.Postgres
		db, err = NewPostgresConnection(pgConf.Host, pgConf.Port, pgConf.User, pgConf.Password, pgConf.DBName)
	case DriverSQLite:
		db, err = NewSQLiteConnection(conf.SQLite.Path)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", conf.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Centralized AutoMigrate for all supported databases
	err = db.AutoMigrate(
		&shelf.Entity{},
		&bin.Entity{},
		&record.Entity{},
		&track.Entity{},
		&playlist.Entity{},
	)
	if err != nil {
		return nil, fmt.Errorf("auto-migrate failed: %w", err)
	}

	return db, nil
}
