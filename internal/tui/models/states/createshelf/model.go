package createshelf

import (
	"github.com/charmbracelet/huh"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"
	teaCmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"
)

type CreateShelfState struct {
	createShelfForm *form
}

func New() CreateShelfState {
	f := newShelfCreateForm()

	return CreateShelfState{
		createShelfForm: f,
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

		ns := shelf.New(newShelf, style.ActiveTextStyle)

		cmds = append(cmds,
			statemachine.WithNextState(statemachine.MainMenu),
		)
	} else {
		cmds = append(cmds,
			teaCmds.WithLayoutUpdate(layouts.Overlay, s.createShelfForm.View()),
		)
	}


	return s, tea.Batch(cmds...)
}

func (s CreateShelfState) View() string {
	view := s.createShelfForm.View()
	return view
}
