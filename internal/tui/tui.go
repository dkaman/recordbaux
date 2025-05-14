package tui

import (
	"fmt"
	"io"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

// shelfDelegate implements Bubble Tea's list.ItemDelegate for shelves
type shelfDelegate struct{ focused bool }

func (d shelfDelegate) Height() int  { return 1 }
func (d shelfDelegate) Spacing() int { return 0 }
func (d shelfDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	sh, ok := listItem.(*physical.Shelf)
	if !ok {
		return
	}

	display := fmt.Sprintf("%s (%d bins × size %d)", sh.Name, len(sh.Bins), sh.BinSize)

	sty := style.TextStyle

	if m.Index() == index && d.focused {
		sty = style.ActiveTextStyle
	}

	prefix := "  "
	if m.Index() == index {
		prefix = "> "
	}

	w.Write([]byte(sty.Render(prefix + display)))
}
func (d shelfDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// State represents the TUI's modes
type State int

const (
	MainMenuState State = iota
	CreateShelfFormState
	QuitState
)

// Model holds the application state
type Model struct {
	cfg         *config.Config
	state       State
	nextState   State
	shelfForm   *form
	loadedShelf *physical.Shelf

	Shelves list.Model
}

// validateNum ensures a non-empty, numeric input
func validateNum(s string) error {
	if s == "" {
		return fmt.Errorf("required")
	}
	if _, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("must be a number")
	}
	return nil
}

// New initializes the TUI model
func New(c *config.Config) Model {
	shelves := list.New([]list.Item{}, shelfDelegate{focused: true}, 0, 10)
	shelves.DisableQuitKeybindings()
	shelves.SetShowTitle(true)
	shelves.Title = "shelves"
	shelves.Styles.Title = style.LabelStyle
	shelves.Styles.TitleBar = style.LabelStyle
	shelves.Styles.NoItems = style.PlaceholderStyle
	shelves.SetStatusBarItemName("shelf", "shelves")

	m := Model{
		cfg:     c,
		state:   MainMenuState,
		Shelves: shelves,
	}

	m.shelfForm = newShelfCreateForm()

	return m
}

// Init is the Bubble Tea initialization command
func (m Model) Init() tea.Cmd {
	return nil
}

// Update routes messages based on the current state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.state = m.nextState

	switch m.state {
	case MainMenuState:
		newShelf, listCmds := m.Shelves.Update(msg)
		m.Shelves = newShelf
		cmds = append(cmds, listCmds)

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "ctrl+c", "esc":
				m.nextState = QuitState
			case "o":
				m.nextState = CreateShelfFormState
				cmds = append(cmds, m.shelfForm.Init())
			case "enter":
				i, ok := m.Shelves.SelectedItem().(*physical.Shelf)
				if ok {
					m.loadedShelf = i
				}
			}
		}

	case CreateShelfFormState:
		// drive the form
		fModel, formUpdateCmds := m.shelfForm.Update(msg)
		if f, ok := fModel.(*form); ok {
			m.shelfForm = f
		}

		cmds = append(cmds, formUpdateCmds)

		// once done
		if m.shelfForm.State == huh.StateCompleted {
			x := m.shelfForm.DimX()
			y := m.shelfForm.DimY()
			size := m.shelfForm.BinSize()

			var totalBins int
			if m.shelfForm.Shape() == Rect {
				totalBins = x * y
			} else {
				totalBins = m.shelfForm.NumBins()
			}

			newShelf := physical.NewShelf(m.shelfForm.Name(), totalBins, size)
			insCmd := m.Shelves.InsertItem(0, newShelf)

			cmds = append(cmds, insCmd)

			m.shelfForm = newShelfCreateForm()

			m.nextState = MainMenuState
		}

	case QuitState:
		return m, tea.Quit
	}

	return m, tea.Batch(cmds...)
}

// View renders UI based on current state
func (m Model) View() string {
	switch m.state {
	case MainMenuState:
		if len(m.Shelves.Items()) == 0 {
			return "no currently defined shelves"
		}

		list := m.Shelves.View()

		if m.loadedShelf != nil {
			name := m.loadedShelf.Name
			nBins := len(m.loadedShelf.Bins)
			sz := m.loadedShelf.BinSize
			capacity := nBins * sz

			list = list + fmt.Sprintf(
				"\n\nshelf name: %s\nnum bins: %d\nbin size: %d\n\nshelf %s has a capacity of %d records!",
				name, nBins, sz, name, capacity,
			)
		}

		return list + "\n"

	case CreateShelfFormState:
		return m.shelfForm.View()

	case QuitState:
		return ""

	default:
		return m.Shelves.View()
	}
}
