package infra

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/playlist"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/db/track"
)

type PostgresConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	DBName   string `koanf:"dbname"`
}

// NewPostgresRepo opens a connection and migrates the schema based on your Entity definitions.
func NewPostgresConnection(host string, port int, user string, password string, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		host, port, user, password, dbName,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}

	err = gormDB.AutoMigrate(
		&shelf.Entity{},
		&bin.Entity{},
		&record.Entity{},
		&track.Entity{},
		&playlist.Entity{},
	)
	if err != nil {
		return nil, fmt.Errorf("auto-migrate failed: %w", err)
	}

	return gormDB, nil
}
