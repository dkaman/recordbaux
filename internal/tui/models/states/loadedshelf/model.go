package loadedshelf

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
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

type refreshLoadedShelfMsg struct {}

type LoadedShelfState struct {
	app         *app.App
	nextState   statemachine.StateType
	shelf       shelf.Model
	selectedBin int
	keys        keyMap
}

// New constructs a LoadedShelfState ready to receive a LoadShelfMsg
func New(a *app.App) LoadedShelfState {
	return LoadedShelfState{
		app:       a,
		nextState: statemachine.Undefined,
		keys:      defaultKeys(),
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
		cmds = append(cmds,
			teaCmds.WithLayoutUpdate(layouts.Viewport, s.View()),
		)

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

	cmds = append(cmds,
		teaCmds.WithLayoutUpdate(layouts.Viewport, s.View()),
	)

	return s, tea.Batch(cmds...)
}

func (s LoadedShelfState) View() string {
	sh := s.shelf.PhysicalShelf()

	if sh == nil {
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
	for i, bin := range sh.Bins {
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

func (s LoadedShelfState) Next() (statemachine.StateType, bool) {
	if s.nextState != statemachine.Undefined {
		return s.nextState, true
	}

	return statemachine.Undefined, false
}

func (s LoadedShelfState) Transition() {
	s.nextState = statemachine.Undefined
}
