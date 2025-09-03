package bin

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/dkaman/recordbaux/internal/db/bin"
)

type LoadBinMsg struct {
	Phy *bin.Entity
}

func WithPhysicalBin(p *bin.Entity) tea.Cmd {
	return func() tea.Msg { return LoadBinMsg{Phy: p} }
}
