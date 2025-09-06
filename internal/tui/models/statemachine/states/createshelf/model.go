package createshelf

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type refreshMsg struct{}

type CreateShelfState struct {
	svcs            *services.AllServices
	createShelfForm *form
	logger          *slog.Logger

	width, height int
}

func New(svcs *services.AllServices, log *slog.Logger) CreateShelfState {
	f := newShelfCreateForm()
	logger := log.WithGroup("createshelfstate")

	return CreateShelfState{
		svcs:            svcs,
		createShelfForm: f,
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

	switch msg := msg.(type) {
	case refreshMsg:
		return s, tea.Batch(s.createShelfForm.Init())

	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height
		return s, tea.Batch(cmds...)
	}

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
			newShelf, _ = shelf.New(s.createShelfForm.Name(), size)
		}

		s.logger.Debug("new shelf", slog.Any("shelf", newShelf))

		s.createShelfForm = newShelfCreateForm()

		return s, tcmds.WithNextState(
			states.MainMenu,
			[]tea.Cmd{s.svcs.SaveShelfCmd(newShelf)},
			nil,
		)
	}

	return s, tea.Batch(cmds...)
}

func (s CreateShelfState) View() string {
	return s.renderModel()
}

func (s CreateShelfState) Help() string {
	return "enter shelf details..."
}
