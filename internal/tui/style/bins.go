package style

import (
	lipgloss "charm.land/lipgloss/v2"
)

var (
	BaseBinStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)
