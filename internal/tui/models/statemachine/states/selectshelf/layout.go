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
	return base, nil
}
