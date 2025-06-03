package createshelf

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

var (
	viewportStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Margin(0)

	formStyle = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
)

func newCreateShelfLayout(base *div.Div, f *form) (*div.Div, error) {
	base.ClearChildren()

	base.AddChild(&div.TextNode{
		Body: f.View(),
	})

	return base, nil
}
