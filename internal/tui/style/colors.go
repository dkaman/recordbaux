package style

import (
	"github.com/charmbracelet/huh"
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

func FormTheme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(lightGreen)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(lightMagenta)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(lightMagenta)
	t.Focused.Description = t.Focused.Description.Foreground(lightGreen)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(lightRed)
	t.Focused.Directory = t.Focused.Directory.Foreground(lightMagenta)
	t.Focused.File = t.Focused.File.Foreground(darkWhite)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(lightRed)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(lightYellow)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(lightYellow)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(lightYellow)
	t.Focused.Option = t.Focused.Option.Foreground(darkWhite)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(lightYellow)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(lightGreen)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(lightGreen)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(darkWhite)
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(darkWhite)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(lightYellow).Background(lightMagenta).Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(darkWhite).Background(lightBlack)

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(lightYellow)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(darkWhite)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(lightYellow)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}
