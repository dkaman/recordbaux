package layout

import (
	"errors"

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
