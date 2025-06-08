package style

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	BarStyle = lipgloss.NewStyle().
		Background(DarkBlue).
		Foreground(DarkBlack).
		Bold(true)

	HelpBarStyle = lipgloss.NewStyle().
		Background(LightGrey).
		Foreground(LightBlack)
)
