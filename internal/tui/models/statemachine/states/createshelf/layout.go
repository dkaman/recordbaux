package createshelf

import (
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

// newCreateShelfLayout centers the shelf‐creation form in a bordered box
// that always takes up 1/3 of the viewport’s width and height.
func newCreateShelfLayout(base *layout.Div, f *form) (*layout.Div, error) {
	base.ClearChildren()

	// 1) Grab current viewport dimensions (set by t.WindowSizeMsg → base.Resize)
	w := base.Width()
	h := base.Height()

	cb := layout.CenteredBox(w, h, f.View(), 1.0/3, 1.0/3)

	// 6) Attach it as the lone child of the viewport
	base.AddChild(cb)

	return base, nil
}

func addViewportText(d *layout.Div, f *form) {
	if cb := d.Find("centerbox"); cb != nil {
		cb.ClearChildren()
		cb.AddChild(&layout.TextNode{
			Body: f.View(),
		})
	}
}
