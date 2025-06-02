package loadedbin

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

var (
	baseTableStyle = lipgloss.NewStyle().
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center)
)

func newLoadedBinLayout(base *div.Div, m table.Model) (*div.Div, error) {
	return base, nil
}
