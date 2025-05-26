package selectshelf

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

// LoadCollectionFromDiscogsState holds the shelf model and renders it.
type SelectShelfState struct {
	app       *app.App
	keys      keyMap
	nextState statemachine.StateType

	shelfList list.Model
}

// New constructs the LoadCollectionFromDiscogs state with an empty shelf model.
func New(a *app.App) SelectShelfState {
	// create an empty list; width/height can be adjusted
	lst := list.New([]list.Item{}, list.NewDefaultDelegate(), 1000, 20)
	lst.Title = "select a Shelf"
	lst.Styles = style.DefaultListStyles()

	items := make([]list.Item, len(a.Shelves))
	for i, sh := range a.Shelves {
		items[i] = sh
	}

	lst.SetItems(items)

	return SelectShelfState{
		app:       a,
		keys:      defaultKeybinds(),
		nextState: statemachine.Undefined,
		shelfList: lst,
	}
}

type refreshShelvesMsg struct{}

// Init satisfies tea.Model.
func (s SelectShelfState) Init() tea.Cmd {
	return func() tea.Msg {
		return refreshShelvesMsg{}
	}
}

func (s SelectShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case refreshShelvesMsg:
		items := make([]list.Item, len(s.app.Shelves))

		for i, sh := range s.app.Shelves {
			items[i] = sh
		}

		s.shelfList.SetItems(items)

		return s, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Select):
			if sel, ok := s.shelfList.SelectedItem().(shelf.Model); ok {
				s.app.CurrentShelf = sel

				s.nextState = statemachine.LoadedShelf
			}

			return s, tea.Batch(cmds...)

		case key.Matches(msg, s.keys.Back):
			s.nextState = statemachine.MainMenu

			return s, tea.Batch(cmds...)
		}
	}

	listModel, listCmds := s.shelfList.Update(msg)
	cmds = append(cmds,
		listCmds,
	)

	s.shelfList = listModel

	return s, tea.Batch(cmds...)
}

// View renders the shelf view into the TopWindow section.
func (s SelectShelfState) View() string {
	var view string
	return view
}

func (s SelectShelfState) Help() string {
	return "please select a shelf from the list..."
}

func (s SelectShelfState) Next() (statemachine.StateType, bool) {
	if s.nextState != statemachine.Undefined {
		return s.nextState, true
	}

	return statemachine.Undefined, false
}

func (s SelectShelfState) Transition() {
	s.nextState = statemachine.Undefined
}
