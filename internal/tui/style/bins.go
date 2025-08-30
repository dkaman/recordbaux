package style

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

var (
	BaseBinStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)
