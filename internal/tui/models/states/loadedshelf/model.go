package loadedshelf

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/physical"
)

type refreshLoadedShelfMsg struct{}

type LoadedShelfState struct {
	app       *app.App
	keys      keyMap
	help      help.Model
	nextState statemachine.StateType

	shelf       shelf.Model
	selectedBin int
}

// New constructs a LoadedShelfState ready to receive a LoadShelfMsg
func New(a *app.App) LoadedShelfState {
	return LoadedShelfState{
		app:       a,
		keys:      defaultKeybinds(),
		help:      help.New(),
		nextState: statemachine.Undefined,
	}
}

func (s LoadedShelfState) Init() tea.Cmd {
	return func() tea.Msg {
		return refreshLoadedShelfMsg{}
	}
}

func (s LoadedShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case refreshLoadedShelfMsg:
		s.shelf = s.app.CurrentShelf
		s.selectedBin = 0

		return s, tea.Batch(cmds...)

	case tea.KeyMsg:
		sh := s.shelf.PhysicalShelf()

		switch {
		case key.Matches(msg, s.keys.Next):
			if sh != nil {
				s.selectedBin = (s.selectedBin + 1) % len(sh.Bins)
			}

		case key.Matches(msg, s.keys.Prev):
			if sh != nil {
				s.selectedBin = (s.selectedBin - 1 + len(sh.Bins)) % len(sh.Bins)
			}

		case key.Matches(msg, s.keys.Back):
			s.nextState = statemachine.MainMenu

		case key.Matches(msg, s.keys.Load):
			s.nextState = statemachine.LoadCollection

		case msg.String() == "enter":
			b := bin.New(sh.Bins[s.selectedBin], style.ActiveTextStyle)
			s.app.CurrentBin = b
			s.nextState = statemachine.LoadedBin
		}
	}

	return s, tea.Batch(cmds...)
}

func (s LoadedShelfState) View() string {
	sh := s.shelf.PhysicalShelf()

	if sh == nil {
		return "no shelf loaded"
	}

	// Define square box size
	boxSize := 8

	// Lipgloss square box style
	base := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(boxSize*3).
		Height(boxSize)

	// Render each bin into its own square box string
	var boxes []string
	for i, b := range sh.Bins {
		count := len(b.Records)

		// ID on first line, capacity on second
		label := fmt.Sprintf("%s\n%d/%d", b.ID, count, b.Size)

		styleBox := base

		if count > 0 {
			styleBox = styleBox.BorderBackground(lipgloss.Color("62"))
		}

		if i == s.selectedBin {
			styleBox = styleBox.BorderForeground(lipgloss.Color("5"))
		}

		boxes = append(boxes, styleBox.Render(label))
	}

	// Decide number of columns based on shape
	var cols int
	if rect, ok := sh.Shape.(*physical.Rectangular); ok {
		cols = rect.X
	} else {
		cols = len(boxes)
	}

	// Chunk into rows
	var rows []string
	for i := 0; i < len(boxes); i += cols {
		end := i + cols

		if end > len(boxes) {
			end = len(boxes)
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top, boxes[i:end]...)

		rows = append(rows, row)
	}

	// Stack rows vertically
	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return grid
}

func (s LoadedShelfState) Help() string {
	return s.help.View(s.keys)
}

func (s LoadedShelfState) Next() (statemachine.StateType, bool) {
	if s.nextState != statemachine.Undefined {
		return s.nextState, true
	}

	return statemachine.Undefined, false
}

func (s LoadedShelfState) Transition() {
	s.nextState = statemachine.Undefined
}
