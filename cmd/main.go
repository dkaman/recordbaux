package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbletea"
	"golang.org/x/term"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

func main() {
	cfg, err := config.New(
		config.WithEnv("SHELF_"),
		config.WithFile("config/shelf.yaml"),
	)
	if err != nil {
		log.Fatalf("error in model configuration: %v", err)
	}

	totalW, totalH, _ := term.GetSize(int(os.Stdout.Fd()))

	l, err := layout.New(totalW, totalH)
	if err != nil {
		log.Fatalf("error measuring screen size %s", err)
	}

	t := tui.New(cfg, l)

	p := tea.NewProgram(t)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
