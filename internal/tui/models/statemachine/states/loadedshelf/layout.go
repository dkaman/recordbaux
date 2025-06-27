package loadedshelf

import (
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

func newSelectShelfLayout(base *layout.Div, sh shelf.Model) (*layout.Div, error) {
	base.ClearChildren()

	base.AddChild(&layout.TextNode{
		Body: sh.View(),
	})

	return base, nil
}
