package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui"
	"github.com/dkaman/recordbaux/internal/db/shelf"
)

func main() {
	cfg, err := config.New(
		config.WithEnv("SHELF_"),
		config.WithFile("config/shelf.yaml"),
	)
	if err != nil {
		log.Fatalf("error in model configuration: %v", err)
	}

	db := shelf.NewMemoryRepo()

	t := tui.New(cfg, nil, db)
	p := tea.NewProgram(t)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
