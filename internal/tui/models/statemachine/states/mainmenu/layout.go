package mainmenu

import (
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)


func newMainMenuLayout(base *div.Div) (*div.Div, error) {
	base.ClearChildren()

	base.AddChild(&div.TextNode{
		Body: "no shelves defined, 'o' to create...",
	})

	return base, nil
}
