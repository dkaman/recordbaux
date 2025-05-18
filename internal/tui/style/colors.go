package style

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	darkBlack    = lipgloss.ANSIColor(0)
	darkRed      = lipgloss.ANSIColor(1)
	darkGreen    = lipgloss.ANSIColor(2)
	darkYellow   = lipgloss.ANSIColor(3)
	darkBlue     = lipgloss.ANSIColor(4)
	darkMagenta  = lipgloss.ANSIColor(5)
	darkCyan     = lipgloss.ANSIColor(6)
	darkWhite    = lipgloss.ANSIColor(7)
	lightBlack   = lipgloss.ANSIColor(8)
	lightRed     = lipgloss.ANSIColor(9)
	lightGreen   = lipgloss.ANSIColor(10)
	lightYellow  = lipgloss.ANSIColor(11)
	lightBlue    = lipgloss.ANSIColor(12)
	lightMagenta = lipgloss.ANSIColor(13)
	lightCyan    = lipgloss.ANSIColor(14)
	lightWhite   = lipgloss.ANSIColor(15)

	bullet   = "•"
	ellipsis = "…"
)

var (
	BackgroundColor = lipgloss.NewStyle().
			Foreground(darkBlack)

	TextStyle = lipgloss.NewStyle().
			Foreground(darkWhite)

	ActiveTextStyle = lipgloss.NewStyle().
			Foreground(darkWhite).
			Bold(true)

	LabelStyle       = TextStyle
	ActiveLabelStyle = ActiveTextStyle

	TableTitleBarStyle = lipgloss.NewStyle().
				Padding(0, 0, 1, 2)

	TableTitleStyle = lipgloss.NewStyle().
			Foreground(darkBlue).
			Background(lightWhite).
			Padding(0, 1)

	TableSpinnerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#8E8E8E", Dark: "#747373"})

	TableFilterPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#ECFD65"})

	TableFilterCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})

	TableDefaultFilterCharacterMatchStyle = lipgloss.NewStyle().Underline(true)

	TableStatusBarStyle = lipgloss.NewStyle().
				Foreground(lightBlue).
				Padding(0, 0, 1, 2)

	TableStatusEmptyStyle = lipgloss.NewStyle().Foreground(lightWhite)

	TableStatusBarActiveFilterStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"})

	TableStatusBarFilterCountStyle = lipgloss.NewStyle().Foreground(lightWhite)

	TableNoItemsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#909090", Dark: "#626262"})

	TableArabicPaginationStyle = lipgloss.NewStyle().Foreground(darkWhite)

	TablePaginationStyleStyle = lipgloss.NewStyle().PaddingLeft(2) //nolint:mnd

	TableHelpStyleStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2) //nolint:mnd

	TableActivePaginationDotStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"}).
					SetString(bullet)

	TableInactivePaginationDotStyle = lipgloss.NewStyle().
					Foreground(darkWhite).
					SetString(bullet)

	TableDividerDotStyle = lipgloss.NewStyle().
				Foreground(darkWhite).
				SetString(" " + bullet + " ")
)
