package selectshelf

import (
	"github.com/charmbracelet/huh"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"

	teaCmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

// LoadCollectionFromDiscogsState holds the shelf model and renders it.
type SelectShelfState struct {
	selectShelfForm *form
	shelves         []shelf.Model

}

// New constructs the LoadCollectionFromDiscogs state with an empty shelf model.
func New() SelectShelfState {
	f := newShelfSelectForm([]shelf.Model{})

	return SelectShelfState{
		selectShelfForm: f,
	}
}

// Init satisfies tea.Model.
func (s SelectShelfState) Init() tea.Cmd {
	return tea.Batch(
		s.selectShelfForm.Init(),
	)
}

// Update handles incoming LoadCollectionMsg and updates the shelf model.
func (s SelectShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case statemachine.BroadcastNewShelfMsg:
		s.shelves = append(s.shelves, msg.Shelf)
		s.selectShelfForm = newShelfSelectForm(s.shelves)
	default:
		fModel, formUpdateCmds := s.selectShelfForm.Update(msg)
		if f, ok := fModel.(*form); ok {
			s.selectShelfForm = f
		}
		cmds = append(cmds, formUpdateCmds)

		if s.selectShelfForm.State == huh.StateCompleted {
			choice := s.selectShelfForm.Shelf()

			for _, sh := range s.shelves {
				if sh.PhysicalShelf().Name == choice {
					cmds = append(cmds,
						statemachine.WithLoadShelfBroadcast(sh.PhysicalShelf()),
						statemachine.WithNextState(statemachine.LoadedShelf),
					)
				}
			}

			s.selectShelfForm = newShelfSelectForm(s.shelves)
		} else {
			cmds = append(cmds,
				teaCmds.WithLayoutUpdate(layouts.Overlay, s.View()),
			)
		}
	}

	return s, tea.Batch(cmds...)
}

// View renders the shelf view into the TopWindow section.
func (s SelectShelfState) View() string {
	var view string

	if len(s.shelves) == 0 {
		view = "no shelves defined, press 'o' to create new shelf..."
	} else {
		view = s.selectShelfForm.View()
	}

	return view
}
