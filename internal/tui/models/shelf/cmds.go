package shelf

import (
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db/shelf"
)

type LoadShelfMsg struct {
	Phy *shelf.Entity
}

func WithPhysicalShelf(p *shelf.Entity) tea.Cmd {
	return func() tea.Msg {
		return LoadShelfMsg{
			Phy: p,
		}
	}
}
