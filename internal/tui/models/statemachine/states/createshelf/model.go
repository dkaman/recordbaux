package createshelf

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type CreateShelfState struct {
	svcs     *services.AllServices
	logger   *slog.Logger
	handlers *handlers.Registry

	createShelfForm *form

	width, height int
}

func New(svcs *services.AllServices, log *slog.Logger) CreateShelfState {
	f := newShelfCreateForm()
	logger := log.WithGroup("createshelfstate")

	return CreateShelfState{
		svcs:            svcs,
		createShelfForm: f,
		logger:          logger,
		handlers:        getHandlers(),
	}
}

// tea.Model implementation

func (s CreateShelfState) Init() tea.Cmd {
	s.logger.Debug("createshelf state init called")
	return s.createShelfForm.Init()
}

func (s CreateShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if handler, ok := s.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(s, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		s = model.(CreateShelfState)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	var formUpdateCmd tea.Cmd
	s.createShelfForm, formUpdateCmd = util.UpdateModel(s.createShelfForm, msg)

	// once done
	if s.createShelfForm.Form.State == huh.StateCompleted {
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

		return s, tcmds.Transition(
			states.MainMenu,
			[]tea.Cmd{s.svcs.SaveShelfCmd(newShelf)},
			nil,
		)
	}

	cmds = append(cmds, formUpdateCmd)

	return s, tea.Batch(cmds...)
}

func (s CreateShelfState) View() string {
	return s.renderModel()
}

func (s CreateShelfState) Help() string {
	return "enter shelf details..."
}
