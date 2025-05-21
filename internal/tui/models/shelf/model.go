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
		sz := m.physicalShelf.Shape.BinSize()
		capacity := nBins * sz

		var shape string
		switch m.physicalShelf.Shape.(type) {
		case *physical.Rectangular:
			shape = "rectangular"
		case *physical.Irregular:
			shape = "irregular"
		}

		out = fmt.Sprintf(
			"shelf name: %s\nshape: %s\nnum bins: %d\nbin size: %d\n\nshelf %s has a capacity of %d records!",
			name, shape, nBins, sz, name, capacity,
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
