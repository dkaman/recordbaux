package infra

const (
	DriverPostgres = "postgres"
	DriverSQLite = "sqlite"
)

type Config struct {
	Driver   string         `koanf:"driver"`
	Postgres PostgresConfig `koanf:"postgres"`
	SQLite   SQLiteConfig   `koanf:"sqlite"`
}
