package mainmenu

import (
	"github.com/charmbracelet/bubbles/list"

	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)


func newMainMenuLayout(base *layout.Div, lst list.Model) (*layout.Div, error) {
	base.ClearChildren()
	if len(lst.Items()) == 0 {
		base.AddChild(&layout.TextNode{
			Body: "no shelves defined, 'o' to create...",
		})
	} else {
		w := base.Width()
		h := base.Height()

		// this is kinda dumb but it'll hold for now lol
		boxW := int(float64(w) * (1.0/3)) - 2
		boxH := int(float64(h) * (1.0/3)) - 4
		lst.SetSize(boxW, boxH)

		cb := layout.CenteredBox(w, h, lst.View(), 1.0/3, 1.0/3)

		base.AddChild(cb)
	}
	return base, nil
}
