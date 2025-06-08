package selectshelf

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

var (
	listStyle = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
)

// newCreateShelfLayout centers the shelf‐creation form in a bordered box
// that always takes up 1/3 of the viewport’s width and height.
func newSelectShelfLayout(base *div.Div, l list.Model) (*div.Div, error) {
	base.ClearChildren()

	// 1) Grab current viewport dimensions (set by t.WindowSizeMsg → base.Resize)
	w := base.Width()
	h := base.Height()

	// 2) Compute box size as 1/3 of viewport
	boxW := w / 3
	boxH := h / 3

	// 3) Compute margins so the box is centered
	//    leftover width/height divided evenly on both sides
	marginV := (h - boxH) / 2
	marginH := (w - boxW) / 2

	// 4) Create the centered container with border, fixed size, and margins
	centerBox, err := div.New(div.Column, lipgloss.NewStyle(),
		div.WithName("centerbox"),
		div.WithBorder(true),
		div.WithFixedWidth(boxW),
		div.WithFixedHeight(boxH),
		div.WithMargin(marginV, marginH, marginV, marginH),
	)
	if err != nil {
		return base, err
	}

	// 5) Render the Huh form inside
	centerBox.AddChild(&div.TextNode{
		Body: l.View(),
	})

	// 6) Attach it as the lone child of the viewport
	base.AddChild(centerBox)

	return base, nil
}

func addViewportText(d *div.Div, l list.Model) {
	if cb := d.Find("centerbox"); cb != nil {
		cb.ClearChildren()
		cb.AddChild(&div.TextNode{
			Body: l.View(),
		})
	}
}
