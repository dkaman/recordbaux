package style

import (
	lipgloss "charm.land/lipgloss/v2"
)

var (
	ProgressStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)
