package main

import (
	"io"
	"log"
	"log/slog"
	"net/url"
	"os"

	"github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/tui"
)

type logConfig struct {
	Level    string `koanf:"level"`
	Format   string `koanf:"format"`
	Location string `koanf:"location"`
}

func initLogger(c logConfig) (*slog.Logger, error) {
	log.Printf("log config: %v", c)
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
		h = slog.NewTextHandler(w, nil)
	case "json":
		h = slog.NewJSONHandler(w, nil)
	}

	return slog.New(h), nil

}

func main() {
	cfg, err := config.New(
		config.WithEnv("SHELF_"),
		config.WithFile("config/shelf.yaml"),
	)
	if err != nil {
		log.Fatalf("error in model configuration: %v", err)
	}

	var lConf logConfig
	err = cfg.Unmarshal("recordbaux.logs", &lConf)
	if err != nil {
		log.Fatalf("error unmarshalling logging config: %v", err)
	}

	logger, err := initLogger(lConf)
	if err != nil {
		log.Fatalf("error initializing logger: %v", err)
	}

	var pConf shelf.PostgresConfig
	err = cfg.Unmarshal("recordbaux.db", &pConf)
	if err != nil {
		logger.Error("error unmarshalling database config: %v", err)
		os.Exit(1)
	}

	db, err := shelf.NewPostgresRepo(pConf)
	if err != nil {
		logger.Error("error constructing shelf repo",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	t := tui.New(cfg, nil, db)
	p := tea.NewProgram(t)
	if _, err := p.Run(); err != nil {
		logger.Error("error running tea program",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}
