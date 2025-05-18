package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/huh"
)

func FormTheme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(lightGreen)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(lightBlue)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(lightBlue)
	t.Focused.Description = t.Focused.Description.Foreground(lightGreen)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(lightRed)
	t.Focused.Directory = t.Focused.Directory.Foreground(lightBlue)
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
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(lightYellow).Background(lightBlue).Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(darkWhite).Background(lightBlack)

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(lightYellow)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lightBlack)
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
