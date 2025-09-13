package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/db/infra"
	"github.com/dkaman/recordbaux/internal/db/playlist"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/db/track"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui"
)

// top level namespace for use in finding config file and env vars
const (
	Namespace = "recordbaux"
)

// subsections of the config that can be looked up by name
var (
	ConfigDB   = "db"
	ConfigLogs = "logs"
)

func main() {
	// read config from both a file and from env vars, based on the
	// configured namespace.
	cfg, err := config.New(
		config.WithFile(fmt.Sprintf("config/%s.yaml", Namespace)),
		config.WithEnv(Namespace),
	)
	if err != nil {
		log.Fatalf("error in config construction: %v", err)
	}

	// logging configuration
	var lConf logConfig
	err = cfg.Unmarshal(ConfigLogs, &lConf)
	if err != nil {
		log.Fatalf("error unmarshalling logging config: %v", err)
	}

	logger, err := initLogger(lConf)
	if err != nil {
		log.Fatalf("error initializing logger: %v", err)
	}

	// database configuration
	var dbConf dbConfig
	err = cfg.Unmarshal(ConfigDB, &dbConf)
	if err != nil {
		log.Fatalf("error unmarshalling database config: %v", err)
	}

	svcs, err := initServices(dbConf, logger)
	if err != nil {
		log.Fatalf("error constructing shelf repo: %v", err)
	}

	// create the root tea.Model to begin execution loop
	t, err := tui.New(cfg, logger, svcs)
	if err != nil {
		log.Fatalf("error constructing tui model: %v", err)
	}

	p := tea.NewProgram(t)
	if _, err := p.Run(); err != nil {
		log.Fatalf("error running tea program: %v", err)
	}
}

type logConfig struct {
	Level    string `koanf:"level"`
	Format   string `koanf:"format"`
	Location string `koanf:"location"`
}

func initLogger(c logConfig) (*slog.Logger, error) {
	lvl, err := parseLogLevel(c.Level)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(c.Location)
	if err != nil {
		return nil, err
	}

	var w io.Writer
	switch u.Scheme {
	case "file":
		f, err := os.OpenFile(u.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}
		w = f
	case "stream":
		w = os.Stdout
	}

	var h slog.Handler
	switch c.Format {
	case "text":
		h = slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: lvl,
		})
	case "json":
		h = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: lvl,
		})
	}

	return slog.New(h), nil
}

func parseLogLevel(levelStr string) (slog.Level, error) {
	upper := strings.ToUpper(levelStr)

	var lvl slog.Level
	err := lvl.UnmarshalText([]byte(upper))
	if err != nil {
		return 0, fmt.Errorf("failed to parse log level '%s': %w", levelStr, err)
	}

	return lvl, nil
}

type dbConfig struct {
	Driver   string `koanf:"driver"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	DBName   string `koanf:"dbname"`
}

func initServices(c dbConfig, l *slog.Logger) (*services.AllServices, error) {
	switch c.Driver {
	case "postgres":
		gormDB, err := infra.NewPostgresConnection(c.Host, c.Port, c.User, c.Password, c.DBName)
		if err != nil {
			l.Error("database error",
				slog.Any("error", err),
			)
		}

		shelfRepo, err := shelf.NewRepo(gormDB)
		if err != nil {
			l.Error("database error",
				slog.Any("error", err),
			)
		}

		trackRepo, err := track.NewRepo(gormDB)
		if err != nil {
			l.Error("database error",
				slog.Any("error", err),
			)
		}

		playlistRepo, err := playlist.NewRepo(gormDB)
		if err != nil {
			l.Error("database error",
				slog.Any("error", err),
			)
		}

		recordRepo, err := record.NewRepo(gormDB)
		if err != nil {
			l.Error("database error",
				slog.Any("error", err),
			)
		}

		svc := services.New(l, shelfRepo, recordRepo, playlistRepo, trackRepo)

		return svc, nil
	}

	return nil, fmt.Errorf("invalid db driver: %s", c.Driver)
}
