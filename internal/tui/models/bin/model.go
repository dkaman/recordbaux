package bin

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
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
	margin, padding layout.TopRightBottomLeft
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

	d, _ := layout.New(layout.Column, lipgloss.NewStyle(),
		layout.WithFixedWidth(m.width),
		layout.WithFixedHeight(m.height),
		layout.WithBorder(true),
	)

	if m.physicalBin != nil {
		name := m.physicalBin.Label
		cur := len(m.physicalBin.Records)
		sz := m.physicalBin.Size
		out = fmt.Sprintf("%s\n%d/%d", name, cur, sz)

		filled := cur > 0
		switch {
		case !filled && m.selected:
			d.ApplyOption(
				layout.WithStyle(m.sty.EmptySelected),
			)

		case !filled && !m.selected:
			d.ApplyOption(
				layout.WithStyle(m.sty.EmptyUnselected),
			)
		case filled && m.selected:
			d.ApplyOption(
				layout.WithStyle(m.sty.FullSelected),
			)
		case filled && !m.selected:
			d.ApplyOption(
				layout.WithStyle(m.sty.FullUnselected),
			)
		}
	} else {
		out = "no bin loaded"
	}


	d.AddChild(&layout.TextNode{
		Body: out,
	})

	return d.Render()
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
