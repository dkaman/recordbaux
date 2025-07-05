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

	m.SetHeight(base.Height()-2)

	left := style.BaseTableStyle.Render(m.View())
	right := fmt.Sprintf("id: %d\ntitle: %s\nartists: %v\n", r.ID, r.Title, r.Artists)

	d, err := layout.Splitscreen(base.Width()-2, base.Height()-2, left, right)
	if err != nil {
		return nil, err
	}

	base.AddChild(d)

	return base, nil
}
