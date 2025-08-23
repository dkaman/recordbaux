package loadedplaylist

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

func newPlaylistLoadedLayout(base *layout.Div, tracks table.Model) (*layout.Div, error) {
	base.ClearChildren()
	base.AddChild(&layout.TextNode{
		Body: tracks.View(),
	})
	return base, nil
}
