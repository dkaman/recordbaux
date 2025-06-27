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
		base.AddChild(&layout.TextNode{
			Body: lst.View(),
		})
	}
	return base, nil
}
