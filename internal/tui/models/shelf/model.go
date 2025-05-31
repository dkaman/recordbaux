package shelf

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

var (
	boxSize      = 3
	baseBinStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Width(boxSize * 3).
			Height(boxSize)
)

type Model struct {
	selectedBin   int
	physicalShelf *physical.Shelf
	sty           lipgloss.Style
}

func New(p *physical.Shelf, s lipgloss.Style) Model {
	return Model{
		selectedBin:   0,
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
		renderedBins := renderBins(m.physicalShelf.Bins, m.selectedBin)
		out = renderShelf(renderedBins, m.physicalShelf.Shape)
	}

	return m.sty.Render(out)
}

func (m Model) SelectBin(b int) Model {
	numBins := len(m.physicalShelf.Bins)
	m.selectedBin = b % numBins
	return m
}

func (m Model) SelectNextBin() Model {
	numBins := len(m.physicalShelf.Bins)
	m.selectedBin = (m.selectedBin + 1) % numBins
	return m
}

func (m Model) SelectPrevBin() Model {
	numBins := len(m.physicalShelf.Bins)
	m.selectedBin = (m.selectedBin - 1) % numBins
	return m
}

func (m Model) GetSelectedBin() bin.Model {
	b := m.physicalShelf.Bins[m.selectedBin]
	return bin.New(b, style.ActiveTextStyle)
}


func renderBins(bins []*physical.Bin, selected int) []string {
	// Render each bin into its own square box string
	var boxes []string
	for i, b := range bins {
		count := len(b.Records)

		// ID on first line, capacity on second
		label := fmt.Sprintf("%s\n%d/%d", b.ID, count, b.Size)

		styleBox := baseBinStyle

		if count > 0 {
			styleBox = styleBox.BorderBackground(lipgloss.Color("62"))
		}

		if i == selected {
			styleBox = styleBox.BorderForeground(lipgloss.Color("5"))
		}

		boxes = append(boxes, styleBox.Render(label))
	}

	return boxes
}

func renderShelf(renderedBins []string, shape physical.Shape) string {
	// Decide number of columns based on shape
	var cols int
	if rect, ok := shape.(*physical.Rectangular); ok {
		cols = rect.X
	} else {
		cols = len(renderedBins)
	}

	// Chunk into rows
	var rows []string
	for i := 0; i < len(renderedBins); i += cols {
		end := i + cols

		if end > len(renderedBins) {
			end = len(renderedBins)
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top, renderedBins[i:end]...)

		rows = append(rows, row)
	}

	// Stack rows vertically
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m Model) Title() string {
	return m.physicalShelf.Name
}

func (m Model) FilterValue() string {
	return m.physicalShelf.Name
}

func (m Model) Description() string {
	bins := len(m.physicalShelf.Bins)
	cap := bins * m.physicalShelf.Shape.BinSize()
	return fmt.Sprintf("%d bins, capacity %d", bins, cap)
}

func (m Model) PhysicalShelf() *physical.Shelf {
	return m.physicalShelf
}
