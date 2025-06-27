package shelf

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dkaman/recordbaux/internal/db/shelf"
)

type LoadShelfMsg struct {
	phy *shelf.Entity
}

func WithPhysicalShelf(p *shelf.Entity) tea.Cmd {
	return func() tea.Msg {
		return LoadShelfMsg{
			phy: p,
		}
	}
}
