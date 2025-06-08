package style

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	BaseTableStyle = lipgloss.NewStyle().
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	tableHeaderStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Bold(true).
				Foreground(LightCyan)

	tableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	tableSelectedStyle = lipgloss.NewStyle().
				Bold(true)
)

func DefaultTableStyles() table.Styles {
	return table.Styles{
		Header:   tableHeaderStyle,
		Cell:     tableCellStyle,
		Selected: tableSelectedStyle,
	}
}
