package style

import (
	"image/color"

	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

const (
	Bullet   = "•"
	Ellipsis = "…"

	DarkenFactor = 0.5
)

var (
	LightGrey    = lipgloss.Color("#AAAAAA")
	DarkBlack    = lipgloss.Color("0")
	DarkRed      = lipgloss.Color("1")
	DarkGreen    = lipgloss.Color("2")
	DarkYellow   = lipgloss.Color("3")
	DarkBlue     = lipgloss.Color("4")
	DarkMagenta  = lipgloss.Color("5")
	DarkCyan     = lipgloss.Color("6")
	DarkWhite    = lipgloss.Color("7")
	LightBlack   = lipgloss.Color("8")
	LightRed     = lipgloss.Color("9")
	LightGreen   = lipgloss.Color("10")
	LightYellow  = lipgloss.Color("11")
	LightBlue    = lipgloss.Color("12")
	LightMagenta = lipgloss.Color("13")
	LightCyan    = lipgloss.Color("14")
	LightWhite   = lipgloss.Color("15")

	LightGreyDimmed    = Darken(LightGrey, DarkenFactor)
	DarkBlackDimmed    = Darken(DarkBlack, DarkenFactor)
	DarkRedDimmed      = Darken(DarkRed, DarkenFactor)
	DarkGreenDimmed    = Darken(DarkGreen, DarkenFactor)
	DarkYellowDimmed   = Darken(DarkYellow, DarkenFactor)
	DarkBlueDimmed     = Darken(DarkBlue, DarkenFactor)
	DarkMagentaDimmed  = Darken(DarkMagenta, DarkenFactor)
	DarkCyanDimmed     = Darken(DarkCyan, DarkenFactor)
	DarkWhiteDimmed    = Darken(DarkWhite, DarkenFactor)
	LightBlackDimmed   = Darken(LightBlack, DarkenFactor)
	LightRedDimmed     = Darken(LightRed, DarkenFactor)
	LightGreenDimmed   = Darken(LightGreen, DarkenFactor)
	LightYellowDimmed  = Darken(LightYellow, DarkenFactor)
	LightBlueDimmed    = Darken(LightBlue, DarkenFactor)
	LightMagentaDimmed = Darken(LightMagenta, DarkenFactor)
	LightCyanDimmed    = Darken(LightCyan, DarkenFactor)
	LightWhiteDimmed   = Darken(LightWhite, DarkenFactor)
)

var (
	BackgroundColor = lipgloss.NewStyle().
			Foreground(DarkBlack)

	TextStyle = lipgloss.NewStyle().
			Foreground(DarkWhite)

	TextStyleDimmed = lipgloss.NewStyle().
			Foreground(DarkWhiteDimmed)

	ActiveTextStyle = lipgloss.NewStyle().
			Foreground(DarkWhite).
			Bold(true)

	ActiveTextStyleDimmed = lipgloss.NewStyle().
				Foreground(DarkWhiteDimmed).
				Bold(true)

	LabelStyle       = TextStyle
	LabelStyleDimmed = TextStyleDimmed

	ActiveLabelStyle       = ActiveTextStyle
	ActiveLabelStyleDimmed = ActiveTextStyleDimmed

	Centered = lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center)
)


// ---------- Darkening helpers ----------

// clamp01 clamps [0,1]
func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func Darken(c color.Color, factor float64) color.Color {
	f := clamp01(factor)
	// Convert to straight (non-premultiplied) NRGBA first.
	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	if n.A == 0 {
		// palette/ANSI may report A=0; treat as opaque UI color
		n.A = 0xFF
	}
	n.R = uint8(float64(n.R)*f + 0.5)
	n.G = uint8(float64(n.G)*f + 0.5)
	n.B = uint8(float64(n.B)*f + 0.5)
	return n
}
