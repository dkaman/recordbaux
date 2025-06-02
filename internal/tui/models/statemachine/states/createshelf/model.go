package createshelf

import (
	"github.com/charmbracelet/huh"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

type resetFormMsg struct{}

type CreateShelfState struct {
	app             *app.App
	createShelfForm *form
	nextState       states.StateType
	layout          *div.Div
}

func New(a *app.App, l *div.Div) CreateShelfState {
	f := newShelfCreateForm()

	return CreateShelfState{
		app:             a,
		nextState:       states.Undefined,
		createShelfForm: f,
		layout:          l,
	}
}

func (s CreateShelfState) Init() tea.Cmd {
	return func() tea.Msg {
		return resetFormMsg{}
	}
}

func (s CreateShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case resetFormMsg:
		s.createShelfForm = newShelfCreateForm()

		s.layout, _ = newCreateShelfLayout(s.layout, s.createShelfForm)

		cmds = append(cmds,
			s.createShelfForm.Init(),
		)

	default:
		fModel, formUpdateCmds := s.createShelfForm.Update(msg)
		if f, ok := fModel.(*form); ok {
			s.createShelfForm = f
		}
		cmds = append(cmds, formUpdateCmds)

		// TODO replace with div version
		// r := &layout.TeaModelRenderer{
		// 	Model: s.createShelfForm,
		// 	Style: formStyle,
		// }

		// s.layout.AddSection(layoutViewport, r)

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

			ns := shelf.New(newShelf, style.ActiveTextStyle)
			s.app.Shelves = append(s.app.Shelves, ns)

			s.createShelfForm = newShelfCreateForm()
			s.nextState = states.MainMenu
		} else {
		}
	}

	return s, tea.Batch(cmds...)
}

func (s CreateShelfState) View() string {
	return s.layout.Render()
}

func (s CreateShelfState) Help() string {
	return "enter shelf details..."
}

func (s CreateShelfState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}

	return states.Undefined, false
}

func (s CreateShelfState) Transition() {
	s.nextState = states.Undefined
}
