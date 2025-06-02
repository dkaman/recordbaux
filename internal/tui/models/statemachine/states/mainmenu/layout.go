package mainmenu

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

var (
	viewportStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Margin(0)
)

func newMainMenuLayout(base *div.Div) (*div.Div, error) {
	return base, nil
}
