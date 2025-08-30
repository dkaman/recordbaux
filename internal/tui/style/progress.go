package style

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

var (
	ProgressStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)
