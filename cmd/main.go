package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui"
)

func main() {
	cfg, err := config.New(
		config.WithEnv("SHELF_"),
		config.WithFile("config/shelf.yaml"),
	)
	if err != nil {
		log.Fatalf("error in model configuration: %v", err)
	}

	t := tui.New(cfg, nil)
	p := tea.NewProgram(t)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
