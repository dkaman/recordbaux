package layout

import "github.com/charmbracelet/lipgloss"

func CenteredBox(w, h int, content string, fractionW, fractionH float64) *Div {
	boxW := int(float64(w) * fractionW)
	boxH := int(float64(h) * fractionH)

	mh := (w - boxW) / 2
	mv := (h - boxH) / 2

	box, _ := New(Column, lipgloss.NewStyle(),
		WithName("centerbox"),
		WithBorder(true),
		WithFixedWidth(boxW),
		WithFixedHeight(boxH),
		WithMargin(mv, mh, mv, mh),
	)

	box.AddChild(&TextNode{
		Body: content,
	})

	return box
}
