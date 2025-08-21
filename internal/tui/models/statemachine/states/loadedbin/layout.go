package loadedbin

import (
	"github.com/charmbracelet/bubbles/table"

	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	tRecord "github.com/dkaman/recordbaux/internal/tui/models/record"
)

func newLoadedBinLayout(base *layout.Div, m table.Model, r *record.Entity) (*layout.Div, error) {
	base.ClearChildren()

	m.SetHeight(base.Height()-2)
	recordModel := tRecord.New(r)

	left := style.BaseTableStyle.Render(m.View())

	right := recordModel.SetSize((base.Width()-3)/2, base.Height()-2).View()

	d, err := layout.Splitscreen(base.Width()-2, base.Height()-2, left, right)
	if err != nil {
		return nil, err
	}

	base.AddChild(d)

	return base, nil
}
