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
}

// newLayout constructs a default layout with styled bars and bordered viewport.
func NewTwoBarViewportLayout() *TwoBarViewportLayout {
	barStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Bold(true)

	viewportStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Margin(0)

	overlayStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Margin(0)

	return &TwoBarViewportLayout{
		styleTopBar:    barStyle,
		styleStatusBar: barStyle,
		styleViewport:  viewportStyle,
		styleOverlay:   overlayStyle,
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
func (l *TwoBarViewportLayout) Render(top, body, status, ov string) string {
	// Top bar: height 1
	topBar := l.styleTopBar.Width(l.width).Height(1).Render(top)

	// Status bar: height 1
	statusBar := l.styleStatusBar.Width(l.width).Height(1).Render(status)

	// Viewport: remaining height (height - 2 for bars)
	vpHeight := l.height - 2
	if vpHeight < 0 {
		vpHeight = 0
	}

	contentWidth := l.width - 2
	if contentWidth < 0 {
		contentWidth = 0
	}

	contentHeight := vpHeight - 2
	if contentHeight < 0 {
		contentHeight = 0
	}

	viewport := l.styleViewport.Width(contentWidth).Height(contentHeight).Render(body)

	bg := lipgloss.JoinVertical(lipgloss.Left, topBar, viewport, statusBar)

	if ov != "" {
		f := newTextModel(l.styleOverlay.Render(ov))
		b := newTextModel(bg)
		overlayModel := overlay.New(f, b, overlay.Center, overlay.Center, 0, 0)

		bg = overlayModel.View()
	}

	// Stack vertically
	return bg
}
