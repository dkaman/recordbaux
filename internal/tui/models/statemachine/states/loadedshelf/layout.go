package loadedshelf

import (
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

func newSelectShelfLayout(base *div.Div, sh shelf.Model) (*div.Div, error) {
	base.ClearChildren()

	base.AddChild(&div.TextNode{
		Body: sh.View(),
	})

	return base, nil
}
