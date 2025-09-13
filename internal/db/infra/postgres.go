package infra

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

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

	return gormDB, nil
}
