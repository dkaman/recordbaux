package createshelf

import (
	"github.com/charmbracelet/huh"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"

	mms "github.com/dkaman/recordbaux/internal/tui/statemachine/states/mainmenu"
)

type CreateShelfState struct {
	createShelfForm *form
	nextState       statemachine.StateType

	layout  *layouts.TallLayout
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

	s.nextState = statemachine.CreateShelf

	// once done
	if s.createShelfForm.State == huh.StateCompleted {
		x := s.createShelfForm.DimX()
		y := s.createShelfForm.DimY()
		size := s.createShelfForm.BinSize()

		var totalBins int
		if s.createShelfForm.Shape() == Rect {
			totalBins = x * y
		} else {
			totalBins = s.createShelfForm.NumBins()
		}

		newShelf := physical.NewShelf(s.createShelfForm.Name(), totalBins, size)

		cmds = append(cmds, mms.WithShelf(newShelf))

		s.createShelfForm = newShelfCreateForm()

		s.nextState = statemachine.MainMenu
	}

	return s, tea.Batch(cmds...)
}

func (s CreateShelfState) View() string {
	s.layout.WithSection(layouts.Overlay, s.createShelfForm.View())
	return s.createShelfForm.View()
}

func (s CreateShelfState) Next(msg tea.Msg) (*statemachine.StateType, error) {
	return &s.nextState, nil
}
