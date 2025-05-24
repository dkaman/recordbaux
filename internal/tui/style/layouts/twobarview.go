package layouts

import (
	"github.com/charmbracelet/lipgloss"

	overlay "github.com/rmhubbert/bubbletea-overlay"
)

const (
	TopBar Section = iota
	Viewport
	StatusBar
)

type TwoBarViewportLayout struct {
	width, height  int
	styleTopBar    lipgloss.Style
	styleStatusBar lipgloss.Style
	styleViewport  lipgloss.Style
	styleOverlay   lipgloss.Style
	styleHelpBar   lipgloss.Style
}

// newLayout constructs a default layout with styled bars and bordered viewport.
func NewTwoBarViewportLayout() *TwoBarViewportLayout {
	barStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Bold(true)

	viewportStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Margin(0)

	overlayStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Margin(0)

	helpStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("15")).
		Height(1).
		Align(lipgloss.Left)

	return &TwoBarViewportLayout{
		styleTopBar:    barStyle,
		styleStatusBar: barStyle,
		styleViewport:  viewportStyle,
		styleOverlay:   overlayStyle,
		styleHelpBar:   helpStyle,
	}
}

// SetSize updates the available terminal space.
func (l *TwoBarViewportLayout) SetSize(w, h int) {
	l.width = w
	l.height = h
}

// Render composes the three regions in a vertical stack.
// top:   content for the TopBar
// body:  main application view content
// status: content for the StatusBar
func (l *TwoBarViewportLayout) Render(top, body, status, ov, help string) string {
	// Top bar: height 1
	topBar := l.styleTopBar.Width(l.width).Height(1).Render(top)

	// Status bar: height 1
	statusBar := l.styleStatusBar.Width(l.width).Height(1).Render(status)

	helpBar := ""
	helpBarHeight := 0
	if help != "" {
		helpBar = l.styleHelpBar.Width(l.width).Height(1).Render(help)
		helpBarHeight = 1
	}

	vpRegion := l.height - 2 - helpBarHeight
	if vpRegion < 0 {
		vpRegion = 0
	}

	// Compute inner content height (subtract borders)
	contentHeight := vpRegion - 2
	if contentHeight < 0 {
		contentHeight = 0
	}

	// Compute inner content width (subtract borders)
	contentWidth := l.width - 2
	if contentWidth < 0 {
		contentWidth = 0
	}

	viewport := l.styleViewport.Width(contentWidth).Height(contentHeight).Render(body)

	// Assemble sections in order, omitting help if empty
	parts := []string{topBar, viewport}
	if helpBarHeight > 0 {
		parts = append(parts, helpBar)
	}

	parts = append(parts, statusBar)

	// Join vertically to fill available height
	bg := lipgloss.JoinVertical(lipgloss.Left, parts...)

	// Overlay on top if provided
	if ov != "" {
		f := newTextModel(l.styleOverlay.Render(ov))
		b := newTextModel(bg)
		overlayModel := overlay.New(f, b, overlay.Center, overlay.Center, 0, 0)
		bg = overlayModel.View()
	}

	// Stack vertically
	return bg
}
