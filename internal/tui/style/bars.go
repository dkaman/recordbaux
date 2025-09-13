package style

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
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
