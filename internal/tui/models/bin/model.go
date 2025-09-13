package bin

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/db/bin"
)

type Style struct {
	EmptySelected   lipgloss.Style
	EmptyUnselected lipgloss.Style
	FullSelected    lipgloss.Style
	FullUnselected  lipgloss.Style
}

type Model struct {
	physicalBin *bin.Entity

	selected bool

	sty Style

	width, height   int
	border          bool
}

func New(p *bin.Entity, s Style) Model {
	m := Model{
		physicalBin: p,
		sty:         s,
	}

	return m
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

	var style lipgloss.Style

	if m.physicalBin != nil {
		name := m.physicalBin.Label
		cur := len(m.physicalBin.Records)
		sz := m.physicalBin.Size
		out = fmt.Sprintf("%s\n%d/%d", name, cur, sz)

		filled := cur > 0
		switch {
		case !filled && m.selected:
			style = m.sty.EmptySelected
		case !filled && !m.selected:
			style = m.sty.EmptyUnselected
		case filled && m.selected:
			style = m.sty.FullSelected
		case filled && !m.selected:
			style = m.sty.FullUnselected
		}
	} else {
		out = "no bin loaded"
	}

	return style.Width(m.width).Height(m.height).Render(out)
}

func (m Model) SetSize(w, h int) Model {
	m.width, m.height = w, h
	return m
}

func (m Model) Unselect() Model {
	m.selected = false
	return m
}

func (m Model) Select() Model {
	m.selected = true
	return m
}

func (m Model) PhysicalBin() *bin.Entity {
	return m.physicalBin
}
