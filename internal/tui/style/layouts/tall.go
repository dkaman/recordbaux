package layouts

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	tea "github.com/charmbracelet/bubbletea"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

type Section int

const (
	SideBar Section = iota
	TopWindow
	BottomWindow
	Overlay
	StatusLine
)

// textModel is a trivial tea.Model for static overlay text.
type textModel struct{ text string }

func newTextModel(t string) textModel                   { return textModel{text: t} }
func (t textModel) Init() tea.Cmd                       { return nil }
func (t textModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t textModel) View() string                        { return t.text }

// TallLayout represents a three-section layout (A, B, C) plus a status line.
type TallLayout struct {
	container         lipgloss.Style
	sectionA          lipgloss.Style
	sectionB          lipgloss.Style
	sectionC          lipgloss.Style
	sectionOverlay    lipgloss.Style
	statusLineStyle   lipgloss.Style
	innerContentWidth int

	sectionContent map[Section]string
}

// NewTallLayout constructs a TallLayout based on current terminal size.
// A = 1/3 width; B = 2/3 width × 30% height; C = 2/3 width × 70% height.
// Also initializes a default status line style (one line high, full inner width).
func NewTallLayout() *TallLayout {
	// Outer container style
	container := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1)

	// Ghost style to measure section frame sizes
	ghost := lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	// Frame sizes for container and sections
	cFrameW := container.GetHorizontalFrameSize()
	cFrameH := container.GetVerticalFrameSize()
	sFrameW := ghost.GetHorizontalFrameSize()
	sFrameH := ghost.GetVerticalFrameSize()

	// Detect terminal size; fallback to 80×24
	totalW, totalH, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		totalW, totalH = 80, 24
	}

	// Compute inner "content" dimensions (subtract container + section frames)
	innerW := totalW - (cFrameW + sFrameW*2)
	innerH := totalH - (cFrameH + sFrameH*3)

	// Split inner area: A vs B+C and B vs C
	contentAW := innerW / 3
	contentRightW := innerW - contentAW
	contentBH := innerH * 3 / 10
	contentCH := innerH - contentBH

	// Build section styles
	sectionA := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Width(contentAW).Height(innerH+ghost.GetVerticalFrameSize()).
		Align(lipgloss.Center, lipgloss.Center)

	sectionB := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Width(contentRightW).Height(contentBH).
		Align(lipgloss.Center, lipgloss.Center)

	sectionC := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Width(contentRightW).Height(contentCH).
		Align(lipgloss.Center, lipgloss.Center)

	sectionOverlay := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Width(100)

	// Default status line style: one line, full inner width
	statusLineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(innerW)

	return &TallLayout{
		container:         container,
		sectionA:          sectionA,
		sectionB:          sectionB,
		sectionC:          sectionC,
		sectionOverlay:    sectionOverlay,
		statusLineStyle:   statusLineStyle,
		innerContentWidth: innerW,
		sectionContent:    make(map[Section]string),
	}
}

// Render composes sections A, B, C and a status line inside the outer container.
// Render expects three tea.Models for A/B/C and a status string.
func (l *TallLayout) Render() string {
	aView, ok := l.sectionContent[SideBar]
	if !ok {
		aView = ""
	}

	bView, ok := l.sectionContent[TopWindow]
	if !ok {
		bView = ""
	}

	cView, ok := l.sectionContent[BottomWindow]
	if !ok {
		cView = ""
	}

	oView, ok := l.sectionContent[Overlay]
	if !ok {
		oView = ""
	}

	a := l.sectionA.Render(aView)
	b := l.sectionB.Render(bView)
	c := l.sectionC.Render(cView)

	rightCol := lipgloss.JoinVertical(lipgloss.Top, b, c)
	body := lipgloss.JoinHorizontal(lipgloss.Top, a, rightCol)

	s, ok := l.sectionContent[StatusLine]
	if !ok {
		s = ""
	}

	statusLine := l.statusLineStyle.Render(s)
	combined := lipgloss.JoinVertical(lipgloss.Top, body, statusLine)

	if oView == "" {
		return l.container.Render(combined)
	} else {
		o := l.sectionOverlay.Render(oView)
		bg := newTextModel(l.container.Render(combined))
		fg := newTextModel(o)
		overlayModel := overlay.New(fg, bg, overlay.Center, overlay.Center, 0, 0)
		return overlayModel.View()
	}
}

func (l *TallLayout) WithSection(sec Section, content string) *TallLayout {
	l.sectionContent[sec] = content
	return l
}
