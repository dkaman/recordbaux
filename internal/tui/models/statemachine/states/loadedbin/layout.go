package loadedbin

import (
	"github.com/charmbracelet/bubbles/table"

	"github.com/dkaman/recordbaux/internal/tui/style/div"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

func newLoadedBinLayout(base *div.Div, m table.Model) (*div.Div, error) {
	base.ClearChildren()
	base.AddChild(&div.TextNode{
		Body: style.BaseTableStyle.Render(m.View()),
	})
	return base, nil
}
