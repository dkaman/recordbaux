package createshelf

import (
	"log/slog"

	"github.com/charmbracelet/huh"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type refreshMsg struct{}

type CreateShelfState struct {
	shelfService    *services.ShelfService
	createShelfForm *form
	nextState       states.StateType
	layout          *layout.Div
	logger          *slog.Logger
}

func New(s *services.ShelfService, l *layout.Div, log *slog.Logger) CreateShelfState {
	f := newShelfCreateForm()
	lay, _ := newCreateShelfLayout(l, f)
	logger := log.WithGroup("createshelfstate")

	return CreateShelfState{
		shelfService:    s,
		nextState:       states.Undefined,
		createShelfForm: f,
		layout:          lay,
		logger:          logger,
	}
}

// tea.Model implementation

func (s CreateShelfState) Init() tea.Cmd {
	s.logger.Debug("createshelf state init called")
	return s.refresh()
}

func (s CreateShelfState) refresh() tea.Cmd {
	return func() tea.Msg {
		return refreshMsg{}
	}
}

func (s CreateShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.(type) {
	case refreshMsg:
		s.layout, _ = newCreateShelfLayout(s.layout, s.createShelfForm)
		return s, tea.Batch(s.createShelfForm.Init())

	case tea.WindowSizeMsg:
		s.layout, _ = newCreateShelfLayout(s.layout, s.createShelfForm)
		return s, tea.Batch(cmds...)
	}

	fModel, formUpdateCmds := s.createShelfForm.Update(msg)
	if f, ok := fModel.(*form); ok {
		s.createShelfForm = f
	}
	cmds = append(cmds, formUpdateCmds)

	addViewportText(s.layout, s.createShelfForm)

	// once done
	if s.createShelfForm.State == huh.StateCompleted {
		x := s.createShelfForm.DimX()
		y := s.createShelfForm.DimY()
		size := s.createShelfForm.BinSize()

		s.logger.Debug("form complete",
			slog.Int("x", x),
			slog.Int("y", y),
			slog.Int("size", size),
		)

		var newShelf *shelf.Entity

		if s.createShelfForm.Shape() == Rect {
			newShelf, _ = shelf.New(s.createShelfForm.Name(), size,
				shelf.WithShapeRect(x, y, size, bin.SortAlphaByArtist),
			)
		} else {
			newShelf, _ = shelf.New(s.createShelfForm.Name(), size,
			)
		}

		s.logger.Debug("new shelf", slog.Any("shelf", newShelf))

		cmds = append(cmds, tcmds.SaveShelfCmd(s.shelfService.Shelves, newShelf, s.logger))

		s.createShelfForm = newShelfCreateForm()

		s.nextState = states.MainMenu
	}

	return s, tea.Batch(cmds...)
}

func (s CreateShelfState) View() string {
	return ""
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

func (s CreateShelfState) Transition() states.State {
	s.nextState = states.Undefined
	return s
}
