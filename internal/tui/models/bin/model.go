package bin

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/physical"
)

type Model struct {
	physicalBin *physical.Bin
	sty         lipgloss.Style
}

func New(p *physical.Bin, s lipgloss.Style) Model {
	return Model{
		physicalBin: p,
		sty:           s,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var out string

	if m.physicalBin == nil {
		out = "no bin loaded"

	} else {
		id := m.physicalBin.ID
		sz := m.physicalBin.Size


		out = fmt.Sprintf("bin id: %s\nsize: %s\n", id, sz)
	}

	return m.sty.Render(out)
}

func (m Model) PhysicalBin() *physical.Bin {
	return m.physicalBin
}
