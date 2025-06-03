package loadedshelf

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

var (
	baseBinStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)

func newSelectShelfLayout(base *div.Div, sh shelf.Model) (*div.Div, error) {
	base.ClearChildren()
	base.AddChild(&div.TextNode{
		Body: sh.View(),
	})
	return base, nil
}
