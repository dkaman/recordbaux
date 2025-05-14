package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/dkaman/shelf/internal/config"
	"github.com/dkaman/shelf/internal/tui"
)

func main() {
	cfg, err := config.New(
		config.WithEnv("SHELF_"),
		config.WithFile("config/shelf.yaml"),
	)
	if err != nil {
		log.Fatalf("error in model configuration: %v", err)
	}

	p := tea.NewProgram(tui.New(cfg))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
