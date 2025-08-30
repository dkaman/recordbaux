package style

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

const (
	DarkBlack    = lipgloss.ANSIColor(0)
	DarkRed      = lipgloss.ANSIColor(1)
	DarkGreen    = lipgloss.ANSIColor(2)
	DarkYellow   = lipgloss.ANSIColor(3)
	DarkBlue     = lipgloss.ANSIColor(4)
	DarkMagenta  = lipgloss.ANSIColor(5)
	DarkCyan     = lipgloss.ANSIColor(6)
	DarkWhite    = lipgloss.ANSIColor(7)
	LightBlack   = lipgloss.ANSIColor(8)
	LightRed     = lipgloss.ANSIColor(9)
	LightGreen   = lipgloss.ANSIColor(10)
	LightYellow  = lipgloss.ANSIColor(11)
	LightBlue    = lipgloss.ANSIColor(12)
	LightMagenta = lipgloss.ANSIColor(13)
	LightCyan    = lipgloss.ANSIColor(14)
	LightWhite   = lipgloss.ANSIColor(15)

	Bullet   = "•"
	Ellipsis = "…"
)

var (
	LightGrey    = lipgloss.Color("#AAAAAA")
)

var (
	BackgroundColor = lipgloss.NewStyle().
			Foreground(DarkBlack)

	TextStyle = lipgloss.NewStyle().
			Foreground(DarkWhite)

	ActiveTextStyle = lipgloss.NewStyle().
			Foreground(DarkWhite).
			Bold(true)

	LabelStyle       = TextStyle

	ActiveLabelStyle = ActiveTextStyle

	Centered = lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center)
)
