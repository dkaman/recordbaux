package createshelf

import (
	"github.com/charmbracelet/huh"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"
)

type CreateShelfState struct {
	createShelfForm *form

	layout *layouts.TallLayout
}

func New(l *layouts.TallLayout) CreateShelfState {
	f := newShelfCreateForm()

	return CreateShelfState{
		createShelfForm: f,
		layout:          l,
	}
}

func (s CreateShelfState) Init() tea.Cmd {
	return s.createShelfForm.Init()
}

func (s CreateShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	fModel, formUpdateCmds := s.createShelfForm.Update(msg)
	if f, ok := fModel.(*form); ok {
		s.createShelfForm = f
	}
	cmds = append(cmds, formUpdateCmds)

	// once done
	if s.createShelfForm.State == huh.StateCompleted {
		x := s.createShelfForm.DimX()
		y := s.createShelfForm.DimY()
		size := s.createShelfForm.BinSize()

		var shape physical.Shape

		if s.createShelfForm.Shape() == Rect {
			shape = &physical.Rectangular{
				X:    x,
				Y:    y,
				Size: size,
			}
		} else {
			shape = &physical.Irregular{
				N:    s.createShelfForm.NumBins(),
				Size: size,
			}
		}

		newShelf, _ := physical.New(s.createShelfForm.Name(),
			physical.WithShelfSortFunc(physical.AlphaByArtist),
			physical.WithShape(shape),
		)

		s.createShelfForm = newShelfCreateForm()

		cmds = append(cmds,
			statemachine.WithNewShelf(newShelf),
			statemachine.WithNextState(statemachine.MainMenu),
		)
	}

	return s, tea.Batch(cmds...)
}

func (s CreateShelfState) View() string {
	view := s.createShelfForm.View()
	s.layout.WithSection(layouts.Overlay, view)
	return view
}
