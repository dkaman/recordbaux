package loadedshelf

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"

	teaCmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type binKey = key.Binding

type keyMap struct {
	Next binKey
	Prev binKey
	Back binKey
	Load binKey
}

func defaultKeys() keyMap {
	return keyMap{
		Next: key.NewBinding(key.WithKeys("n")),
		Prev: key.NewBinding(key.WithKeys("N")),
		Back: key.NewBinding(key.WithKeys("q")),
		Load: key.NewBinding(key.WithKeys("l")),
	}
}

type LoadedShelfState struct {
	shelf       *physical.Shelf
	selectedBin int
	keys        keyMap
}

// New constructs a LoadedShelfState ready to receive a LoadShelfMsg
func New() LoadedShelfState {
	return LoadedShelfState{
		keys:   defaultKeys(),
	}
}

func (s LoadedShelfState) Init() tea.Cmd {
	return teaCmds.WithLayoutUpdate(layouts.Viewport, s.View())
}

func (s LoadedShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case statemachine.BroadcastLoadShelfMsg:
		s.shelf = msg.Shelf
		s.selectedBin = 0

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Next):
			if s.shelf != nil {
				s.selectedBin = (s.selectedBin + 1) % len(s.shelf.Bins)
			}

		case key.Matches(msg, s.keys.Prev):
			if s.shelf != nil {
				s.selectedBin = (s.selectedBin - 1 + len(s.shelf.Bins)) % len(s.shelf.Bins)
			}

		case key.Matches(msg, s.keys.Back):
			cmds = append(cmds,
				statemachine.WithNextState(statemachine.MainMenu),
			)
			s.shelf = nil

		case key.Matches(msg, s.keys.Load):
			cmds = append(cmds,
				statemachine.WithLoadShelfBroadcast(s.shelf),
				statemachine.WithNextState(statemachine.LoadCollection),
			)

		case msg.String() == "enter":
			cmds = append(cmds,
				statemachine.WithNextState(statemachine.LoadedBin),
				statemachine.WithLoadBinBroadcast(s.shelf.Bins[s.selectedBin]),
			)

			return s, tea.Sequence(cmds...)
		}
	}

	cmds = append(cmds,
			teaCmds.WithLayoutUpdate(layouts.Viewport, s.View()),
	)

	return s, tea.Batch(cmds...)
}

func (s LoadedShelfState) View() string {
	if s.shelf == nil {
		return "no shelf loaded"
	}

	// Lipgloss box style
	base := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Padding(0, 1).
		Margin(0, 1).
		Align(lipgloss.Center).
		Width(12)

	// Render each bin into its own little box string
	var boxes []string
	for i, bin := range s.shelf.Bins {
		count := len(bin.Records)
		label := fmt.Sprintf("%s %d/%d", bin.ID, count, bin.Size)

		style := base
		if count > 0 {
			style = style.BorderBackground(lipgloss.Color("62"))
		}
		if i == s.selectedBin {
			label = "★ " + label
			style = style.BorderForeground(lipgloss.Color("5"))
		}

		boxes = append(boxes, style.Render(label))
	}

	// Decide how many columns per row
	cols := 4
	var rows []string

	// Chunk into rows of `cols`
	for i := 0; i < len(boxes); i += cols {
		end := i + cols
		if end > len(boxes) {
			end = len(boxes)
		}
		// join boxes[i:end] horizontally, aligning tops
		row := lipgloss.JoinHorizontal(lipgloss.Top, boxes[i:end]...)
		rows = append(rows, row)
	}

	// stack all rows vertically, left-aligned
	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return grid
}
