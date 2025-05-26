package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/huh"
)

func DefaultFormStyles() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(LightGreen)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(LightBlue)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(LightBlue)
	t.Focused.Description = t.Focused.Description.Foreground(LightGreen)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(LightRed)
	t.Focused.Directory = t.Focused.Directory.Foreground(LightBlue)
	t.Focused.File = t.Focused.File.Foreground(DarkWhite)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(LightRed)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(LightYellow)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(LightYellow)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(LightYellow)
	t.Focused.Option = t.Focused.Option.Foreground(DarkWhite)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(LightYellow)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(LightGreen)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(LightGreen)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(DarkWhite)
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(DarkWhite)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(LightYellow).Background(LightBlue).Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(DarkWhite).Background(LightBlack)

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(LightYellow)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(LightBlack)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(LightYellow)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}
