package loadedbin

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"

	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

func newLoadedBinLayout(base *layout.Div, m table.Model, r *record.Entity) (*layout.Div, error) {
	base.ClearChildren()

	base.AddChild(&layout.TextNode{
		Body: style.BaseTableStyle.Render(m.View()),
	})

	id := r.Release.ID
	title := r.Release.BasicInfo.Title
	artists := r.Release.BasicInfo.Artists

	info := fmt.Sprintf("id: %d\ntitle: %s\nartists: %v\n", id, title, artists)
	base.AddChild(&layout.TextNode{
		Body: info,
	})

	return base, nil
}
