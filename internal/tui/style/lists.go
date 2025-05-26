package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	listTitleBarStyle = lipgloss.NewStyle().
				Padding(0, 0, 1, 3).
				Background(LightBlue).
				Foreground(DarkBlack)

	listTitleStyle = lipgloss.NewStyle().
			Bold(true)

	listSpinnerStyle                     = lipgloss.NewStyle().Foreground(LightGreen)
	listFilterPromptStyle                = lipgloss.NewStyle().Foreground(LightGreen)
	listFilterCursorStyle                = lipgloss.NewStyle().Foreground(LightGreen)
	listDefaultFilterCharacterMatchStyle = lipgloss.NewStyle().Foreground(LightGreen)
	listStatusBarStyle                   = lipgloss.NewStyle().Foreground(LightGreen)
	listStatusEmptyStyle                 = lipgloss.NewStyle().Foreground(LightGreen)
	listStatusBarActiveFilterStyle       = lipgloss.NewStyle().Foreground(LightGreen)
	listStatusBarFilterCountStyle        = lipgloss.NewStyle().Foreground(LightGreen)
	listNoItemsStyle                     = lipgloss.NewStyle().Foreground(LightGreen)
	listPaginationStyleStyle             = lipgloss.NewStyle().Foreground(LightGreen)
	listHelpStyleStyle                   = lipgloss.NewStyle().Foreground(LightGreen)
	listActivePaginationDotStyle         = lipgloss.NewStyle().Foreground(LightGreen)
	listInactivePaginationDotStyle       = lipgloss.NewStyle().Foreground(LightGreen)
	listArabicPaginationStyle            = lipgloss.NewStyle().Foreground(LightGreen)
	listDividerDotStyle                  = lipgloss.NewStyle().Foreground(LightGreen)
)

func DefaultListStyles() list.Styles {
	s := list.DefaultStyles()

	s.TitleBar = listTitleBarStyle
	s.Title = listTitleStyle
	s.Spinner = listSpinnerStyle
	s.FilterPrompt = listFilterPromptStyle
	s.FilterCursor = listFilterCursorStyle

	s.DefaultFilterCharacterMatch = listDefaultFilterCharacterMatchStyle

	s.StatusBar = listStatusBarStyle
	s.StatusEmpty = listStatusEmptyStyle
	s.StatusBarActiveFilter = listStatusBarActiveFilterStyle
	s.StatusBarFilterCount = listStatusBarFilterCountStyle

	s.NoItems = listNoItemsStyle

	s.PaginationStyle = listPaginationStyleStyle
	s.HelpStyle = listHelpStyleStyle

	s.ActivePaginationDot = listActivePaginationDotStyle
	s.InactivePaginationDot = listInactivePaginationDotStyle
	s.ArabicPagination = listArabicPaginationStyle
	s.DividerDot = listDividerDotStyle

	return s
}
