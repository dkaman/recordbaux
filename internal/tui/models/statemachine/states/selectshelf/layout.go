package selectshelf

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

var (
	listStyle = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
)

func newSelectShelfLayout(base *div.Div, m list.Model) (*div.Div, error) {
	base.ClearChildren()
	base.AddChild(&div.TextNode{
		Body: m.View(),
	})
	return base, nil
}
