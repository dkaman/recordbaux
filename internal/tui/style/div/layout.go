package div

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type JoinDirection int

const (
	Row JoinDirection = iota
	Column
)

type Renderer interface {
	Render() string
}

type Resizer interface {
	Resize(width, height int)
}

type Node interface {
	Renderer
	Resizer
	Children() []Node
	AddChild(Node)
}

type divOption func(*Div) error

type TopRightBottomLeft struct{ Top, Right, Bottom, Left int }

type Div struct {
	name     string
	children []Node

	style lipgloss.Style

	direction JoinDirection

	margin  TopRightBottomLeft
	padding TopRightBottomLeft
	border  bool

	hidden bool

	width, height  int
	fixedW, fixedH bool
}

func New(d JoinDirection, s lipgloss.Style, opts ...divOption) (*Div, error) {
	div := &Div{
		name:      "",
		children:  nil,
		style:     s,
		direction: d,
		margin:    TopRightBottomLeft{0, 0, 0, 0},
		padding:   TopRightBottomLeft{0, 0, 0, 0},
		border:    false,
		hidden:    false,
		width:     0,
		height:    0,
		fixedW:    false,
		fixedH:    false,
	}

	for _, o := range opts {
		err := o(div)
		if err != nil {
			return nil, fmt.Errorf("error in div option: %s", err)
		}
	}

	return div, nil
}

func (d *Div) ApplyOption(opts ...divOption) error {
	for _, o := range opts {
		err := o(d)
		if err != nil {
			return fmt.Errorf("error applying option to div: %s", err)
		}
	}

	return nil
}

func WithName(name string) divOption {
	return func(d *Div) error {
		d.name = name
		return nil
	}
}

func WithFixedHeight(h int) divOption {
	return func(d *Div) error {
		d.height = h
		d.fixedH = true
		return nil
	}
}

func WithFixedWidth(w int) divOption {
	return func(d *Div) error {
		d.width = w
		d.fixedW = true
		return nil
	}
}

func WithChild(child Node) divOption {
	return func(d *Div) error {
		d.children = append(d.children, child)
		return nil
	}
}

func WithMargin(top, right, bottom, left int) divOption {
	return func(d *Div) error {
		d.margin = TopRightBottomLeft{top, right, bottom, left}
		return nil
	}
}

func WithPadding(top, right, bottom, left int) divOption {
	return func(d *Div) error {
		d.padding = TopRightBottomLeft{top, right, bottom, left}
		return nil
	}
}

func WithBorder(border bool) divOption {
	return func(d *Div) error {
		d.border = border
		return nil
	}
}

func WithHidden(h bool) divOption {
	return func(d *Div) error {
		d.hidden = h
		return nil
	}
}

func (d *Div) Render() string {
	baseStyle := d.style

	var renderedChildren []string
	for _, child := range d.children {
		if divChild, ok := child.(*Div); ok && divChild.hidden {
			continue
		}
		renderedChildren = append(renderedChildren, child.Render())
	}

	var joined string
	if d.direction == Row {
		joined = lipgloss.JoinHorizontal(lipgloss.Top, renderedChildren...)
	} else {
		joined = lipgloss.JoinVertical(lipgloss.Left, renderedChildren...)
	}

	marginX := d.margin.Left + d.margin.Right
	marginY := d.margin.Top + d.margin.Bottom

	borderX, borderY := 0, 0
	if d.border {
		borderX = 2
		borderY = 2
	}

	contentW := d.width - marginX - borderX
	if contentW < 0 {
		contentW = 0
	}

	contentH := d.height - marginY - borderY
	if contentH < 0 {
		contentH = 0
	}

	style := baseStyle

	if d.border {
		style = style.
			Border(lipgloss.NormalBorder())
	}

	style = style.
		Width(contentW).
		Height(contentH).
		Padding(
			d.padding.Top, d.padding.Right, d.padding.Bottom, d.padding.Left,
		).
		Margin(
			d.margin.Top, d.margin.Right, d.margin.Bottom, d.margin.Left,
		)

	return style.Render(joined)
}

func (d *Div) Resize(w, h int) {
	d.width, d.height = w, h

	borderX, borderY := 0, 0
	if d.border {
		borderX = 2 // 1 char left + 1 char right
		borderY = 2 // 1 row top  + 1 row bottom
	}

	marginX := d.margin.Left + d.margin.Right
	marginY := d.margin.Top + d.margin.Bottom

	padX := d.padding.Left + d.padding.Right
	padY := d.padding.Top + d.padding.Bottom

	innerW := w - marginX - borderX - padX
	innerH := h - marginY - borderY - padY

	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}

	// distribute innerW/innerH among children (depending on row vs.
	// column).
	n := len(d.children)

	if n == 0 {
		return
	}

	if d.direction == Row {
		totalFixed := 0
		flexCount := 0

		for _, child := range d.children {
			if divChild, ok := child.(*Div); ok {
				if divChild.hidden {
					continue
				}

				if divChild.fixedW {
					totalFixed += divChild.width
				} else {
					flexCount++
				}
			}
		}

		remainingW := innerW - totalFixed
		if remainingW < 0 {
			remainingW = 0
		}

		baseW := 0
		remW := 0

		if flexCount > 0 {
			baseW = remainingW / flexCount
			remW = remainingW % flexCount
		}

		usedRem := 0
		for _, child := range d.children {
			// skip hidden Divs
			if divChild, ok := child.(*Div); ok && divChild.hidden {
				continue
			}

			if divChild, ok := child.(*Div); ok && divChild.fixedW {
				child.Resize(divChild.width, innerH)
			} else {
				wi := baseW

				if usedRem < remW {
					wi++
					usedRem++
				}

				child.Resize(wi, innerH)
			}
		}
	} else {
		totalFixed := 0
		flexCount := 0

		for _, child := range d.children {
			if divChild, ok := child.(*Div); ok && divChild.hidden {
				continue
			}

			if divChild, ok := child.(*Div); ok && divChild.fixedH {
				totalFixed += divChild.height
			} else {
				flexCount++
			}
		}

		remainingH := innerH - totalFixed
		if remainingH < 0 {
			remainingH = 0
		}

		baseH := 0
		remH := 0

		if flexCount > 0 {
			baseH = remainingH / flexCount
			remH = remainingH % flexCount
		}

		usedRem := 0
		for _, child := range d.children {
			if divChild, ok := child.(*Div); ok && divChild.hidden {
				continue
			}

			if divChild, ok := child.(*Div); ok && divChild.fixedH {
				child.Resize(innerW, divChild.height)
			} else {
				hi := baseH

				if usedRem < remH {
					hi++
					usedRem++
				}

				child.Resize(innerW, hi)
			}
		}
	}
}

func (d *Div) Width() int {
	return d.width
}

func (d *Div) Height() int {
	return d.height
}

func (d *Div) Children() []Node {
	return d.children
}

func (d *Div) AddChild(child Node) {
	d.children = append(d.children, child)
}

func (d *Div) ClearChildren() {
	d.children = make([]Node, 0)
}

func (d *Div) Hide() {
	d.hidden = true
}

func (d *Div) Show() {
	d.hidden = false
}

func (d *Div) Find(name string) *Div {
	if d.name == name {
		return d
	}

	for _, child := range d.children {
		if sub, ok := child.(*Div); ok {
			if found := sub.Find(name); found != nil {
				return found
			}
		}
	}
	return nil
}

type TextNode struct {
	Body string
}

func (t *TextNode) Render() string {
	return t.Body
}

func (t *TextNode) Resize(w, h int)      {}
func (t *TextNode) Children() []Node { return nil }
func (t *TextNode) AddChild(node Node) {}
