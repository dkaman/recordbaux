package shelf

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"
)

type LoadShelfMsg struct {
	phy *physical.Shelf
}

func WithPhysicalShelf(p *physical.Shelf) tea.Cmd {
	return func() tea.Msg {
		return LoadShelfMsg{
			phy: p,
		}
	}
}
