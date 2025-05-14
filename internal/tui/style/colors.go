package style


import (
	"github.com/charmbracelet/lipgloss"
)

const accentColor = lipgloss.Color("99")
const yellowColor = lipgloss.Color("#ECFD66")
const whiteColor = lipgloss.Color("255")
const grayColor = lipgloss.Color("241")
const darkGrayColor = lipgloss.Color("236")
const lightGrayColor = lipgloss.Color("247")

var (
	ActiveTextStyle = lipgloss.NewStyle().Foreground(whiteColor)
	TextStyle       = lipgloss.NewStyle().Foreground(lightGrayColor)

	ActiveLabelStyle = lipgloss.NewStyle().Foreground(accentColor)
	LabelStyle       = lipgloss.NewStyle().Foreground(grayColor)

	PlaceholderStyle = lipgloss.NewStyle().Foreground(darkGrayColor)
	CursorStyle      = lipgloss.NewStyle().Foreground(whiteColor)

	PaddedStyle = lipgloss.NewStyle().Padding(1)

	ErrorHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#FF5F87")).Bold(true).Padding(0, 1).SetString("ERROR")
	ErrorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87"))
	CommentStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#757575")).PaddingLeft(1)

	SendButtonActiveStyle   = lipgloss.NewStyle().Background(accentColor).Foreground(yellowColor).Padding(0, 2)
	SendButtonInactiveStyle = lipgloss.NewStyle().Background(darkGrayColor).Foreground(lightGrayColor).Padding(0, 2)
	SendButtonStyle         = lipgloss.NewStyle().Background(darkGrayColor).Foreground(grayColor).Padding(0, 2)

	InlineCodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87")).Background(lipgloss.Color("#3A3A3A")).Padding(0, 1)
	LinkStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#00AF87")).Underline(true)
)
