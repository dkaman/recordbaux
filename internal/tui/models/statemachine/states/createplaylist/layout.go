package createplaylist

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

func newCreatePlaylistLayout(base *layout.Div, tracks table.Model) (*layout.Div, error) {
	base.ClearChildren()

	if len(tracks.Rows()) == 0 {
		base.AddChild(&layout.TextNode{
			Body: "no tracks defined, load some into a shelf...",
		})
		return base, nil
	}

	tracks.SetHeight(base.Height()-2)
	tracks.SetWidth(base.Width())

	base.AddChild(&layout.TextNode{
		Body: tracks.View(),
	})

	return base, nil
}
