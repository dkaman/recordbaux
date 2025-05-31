package layout

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	LayoutSectionNotFoundErr = errors.New("requested layout section not found")
)

type Section int

type layoutOption func(*Node) error

type joinFunc func(map[Section]Renderer) string

func WithSection(sec Section, r Renderer) layoutOption {
	return func(l *Node) error {
		l.sections[sec] = r
		return nil
	}
}

func WithJoinFunc(f joinFunc) layoutOption {
	return func(l *Node) error {
		l.join = f
		return nil
	}
}

type Renderer interface {
	Render() string
	Resize(int, int)
}


type Node struct {
	join     joinFunc
	sections map[Section]Renderer

	Width, Height int
}

func New(width, height int, opts ...layoutOption) (*Node, error) {
	l := &Node{
		Width:    width,
		Height:   height,
		sections: make(map[Section]Renderer),
	}

	for _, o := range opts {
		err := o(l)
		if err != nil {
			return nil, err
		}
	}

	return l, nil
}

func defaultJoinFunc(m map[Section]Renderer) string {
	var parts []string
	for _, r := range m {
		rendered := r.Render()
		parts = append(parts, rendered)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func (n *Node) Render() string {
	if n.join == nil {
		n.join = defaultJoinFunc
	}

	for _, r := range n.sections {
		if _, isNode := r.(*Node); !isNode {
			r.Resize(n.Width, n.Height)
		}
	}

	joined := n.join(n.sections)

	return joined
}

func (n *Node) Resize(w, h int) {
	n.Width = w
	n.Height = h
}

func (n *Node) AddSection(s Section, r Renderer) {
	n.sections[s] = r
}

func (n *Node) GetSection(s Section) (Renderer, error) {
	if r, ok := n.sections[s]; !ok {
		return nil, LayoutSectionNotFoundErr
	} else {
		return r, nil
	}
}

func (n *Node) SetJoinFunc(f joinFunc) {
	n.join = f
}

type TextRenderer struct {
	Body  string
	Style lipgloss.Style
}

func (r *TextRenderer) Render() string {
	return r.Style.Render(r.Body)
}

func (r *TextRenderer) Resize(w, h int) {
	r.Style = r.Style.
		Width(w).
		Height(h)
}

type TeaModelRenderer struct {
	Model tea.Model
	Style lipgloss.Style
}

func (r *TeaModelRenderer) Render() string {
	return r.Style.Render(r.Model.View())
}

func (r *TeaModelRenderer) Resize(w, h int) {
	r.Style = r.Style.
		Width(w).
		Height(h)
}

// PrettyPrint builds and returns the tree structure as a string.
func (n *Node) PrettyPrint() string {
	var sb strings.Builder
	n.buildPrettyString(&sb, 0)
	return sb.String()
}

// buildPrettyString recursively writes Node and its children into the builder.
func (n *Node) buildPrettyString(sb *strings.Builder, level int) {
	indent := strings.Repeat("\t", level)
	sb.WriteString(fmt.Sprintf("%sNode (Width=%d, Height=%d)\n", indent, n.Width, n.Height))
	for sec, r := range n.sections {
		sb.WriteString(fmt.Sprintf("%s\tSection %v: %T\n", indent, sec, r))
		if child, ok := r.(*Node); ok {
			child.buildPrettyString(sb, level+1)
		}
	}
}

// // textModel is a trivial tea.Model for static overlay text.
// type textModel struct{ text string }

// func newTextModel(t string) textModel                   { return textModel{text: t} }
// func (t textModel) Init() tea.Cmd                       { return nil }
// func (t textModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
// func (t textModel) View() string                        { return t.text }

// type TwoBarViewportLayout struct {
// 	width, height  int
// 	styleTopBar    lipgloss.Style
// 	styleStatusBar lipgloss.Style
// 	styleViewport  lipgloss.Style
// 	styleOverlay   lipgloss.Style
// 	styleHelpBar   lipgloss.Style
// }

// // newLayout constructs a default layout with styled bars and bordered viewport.
// func NewTwoBarViewportLayout() *TwoBarViewportLayout {
// 	barStyle := lipgloss.NewStyle().
// 		Background(darkBlue).
// 		Foreground(darkBlack).
// 		Bold(true)

// 	viewportStyle := lipgloss.NewStyle().
// 		Border(lipgloss.NormalBorder()).
// 		AlignHorizontal(lipgloss.Center).
// 		AlignVertical(lipgloss.Center).
// 		Margin(0)

// 	overlayStyle := lipgloss.NewStyle().
// 		Border(lipgloss.NormalBorder()).
// 		Margin(0)

// 	helpStyle := lipgloss.NewStyle().
// 		Background(lightBlack).
// 		Foreground(lightWhite).
// 		Height(1).
// 		Align(lipgloss.Left)

// 	return &TwoBarViewportLayout{
// 		styleTopBar:    barStyle,
// 		styleStatusBar: barStyle,
// 		styleViewport:  viewportStyle,
// 		styleOverlay:   overlayStyle,
// 		styleHelpBar:   helpStyle,
// 	}
// }

// // SetSize updates the available terminal space.
// func (l *TwoBarViewportLayout) SetSize(w, h int) {
// 	l.width = w
// 	l.height = h
// }

// // Render composes the three regions in a vertical stack.
// // top:   content for the TopBar
// // body:  main application view content
// // status: content for the StatusBar
// func (l *TwoBarViewportLayout) Render(top, body, status, ov, help string) string {
// 	// Top bar: height 1
// 	topBar := l.styleTopBar.Width(l.width).Height(1).Render(top)

// 	// Status bar: height 1
// 	statusBar := l.styleStatusBar.Width(l.width).Height(1).Render(status)

// 	helpBar := ""
// 	helpBarHeight := 0
// 	if help != "" {
// 		helpBar = l.styleHelpBar.Width(l.width).Height(1).Render(help)
// 		helpBarHeight = 1
// 	}

// 	vpRegion := l.height - 2 - helpBarHeight
// 	if vpRegion < 0 {
// 		vpRegion = 0
// 	}

// 	// Compute inner content height (subtract borders)
// 	contentHeight := vpRegion - 2
// 	if contentHeight < 0 {
// 		contentHeight = 0
// 	}

// 	// Compute inner content width (subtract borders)
// 	contentWidth := l.width - 2
// 	if contentWidth < 0 {
// 		contentWidth = 0
// 	}

// 	viewport := l.styleViewport.Width(contentWidth).Height(contentHeight).Render(body)

// 	// Assemble sections in order, omitting help if empty
// 	parts := []string{topBar, viewport}
// 	if helpBarHeight > 0 {
// 		parts = append(parts, helpBar)
// 	}

// 	parts = append(parts, statusBar)

// 	// Join vertically to fill available height
// 	bg := lipgloss.JoinVertical(lipgloss.Left, parts...)

// 	// Overlay on top if provided
// 	if ov != "" {
// 		minOverlayWidth := (l.width) / 3
// 		maxOverlayWidth := (l.width * 3) / 4
// 		minOverlayHeight := (l.height) / 3
// 		maxOverlayHeight := (l.height * 3) / 4

// 		constrained := l.styleOverlay.
// 			Width(minOverlayWidth).
// 			MaxWidth(maxOverlayWidth).
// 			Height(minOverlayHeight).
// 			MaxHeight(maxOverlayHeight).
// 			Render(ov)

// 		f := newTextModel(constrained)
// 		b := newTextModel(bg)
// 		overlayModel := overlay.New(f, b, overlay.Center, overlay.Center, 0, 0)
// 		bg = overlayModel.View()
// 	}

// 	// Stack vertically
// 	return bg
// }
