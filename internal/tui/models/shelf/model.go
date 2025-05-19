package shelf

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/physical"
)

type Model struct {
	physicalShelf *physical.Shelf
	sty           lipgloss.Style
}

func New(p *physical.Shelf, s lipgloss.Style) Model {
	return Model{
		physicalShelf: p,
		sty:           s,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case LoadShelfMsg:
		m.physicalShelf = msg.phy
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var out string

	if m.physicalShelf == nil {
		out = "no shelves loaded"

	} else {
		name := m.physicalShelf.Name
		nBins := len(m.physicalShelf.Bins)
		sz := m.physicalShelf.BinSize
		capacity := nBins * sz

		out = fmt.Sprintf(
			"shelf name: %s\nnum bins: %d\nbin size: %d\n\nshelf %s has a capacity of %d records!",
			name, nBins, sz, name, capacity,
		)
	}

	return m.sty.Render(out)
}

func (m Model) FilterValue() string {
	return m.physicalShelf.Name
}

func (m Model) PhysicalShelf() *physical.Shelf {
	return m.physicalShelf
}
