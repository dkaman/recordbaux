package infra

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	// "github.com/dkaman/recordbaux/internal/db/bin"
	// "github.com/dkaman/recordbaux/internal/db/playlist"
	// "github.com/dkaman/recordbaux/internal/db/record"
	// "github.com/dkaman/recordbaux/internal/db/shelf"
	// "github.com/dkaman/recordbaux/internal/db/track"
)

type SQLiteConfig struct {
	Path string `koanf:"path"`
}

func NewSQLiteConnection(path string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(path))
}
