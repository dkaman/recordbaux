package style

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpEllipsisStyle       = lipgloss.NewStyle().Bold(true)
	helpShortKeyStyle       = lipgloss.NewStyle().Bold(true)
	helpShortDescStyle      = lipgloss.NewStyle().Bold(true)
	helpShortSeparatorStyle = lipgloss.NewStyle().Bold(true)
	helpFullKeyStyle        = lipgloss.NewStyle().Bold(true)
	helpFullDescStyle       = lipgloss.NewStyle().Bold(true)
	helpFullSeparatorStyle  = lipgloss.NewStyle().Bold(true)
)

func DefaultHelpStyles() help.Styles {
	return help.Styles{
		Ellipsis:       helpEllipsisStyle,
		ShortKey:       helpShortKeyStyle,
		ShortDesc:      helpShortDescStyle,
		ShortSeparator: helpShortSeparatorStyle,
		FullKey:        helpFullKeyStyle,
		FullDesc:       helpFullDescStyle,
		FullSeparator:  helpFullSeparatorStyle,
	}
}
