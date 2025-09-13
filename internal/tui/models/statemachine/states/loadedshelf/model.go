package loadedshelf

import (
	"log/slog"


	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/util"
)

type LoadedShelfState struct {
	svcs     *services.AllServices
	keys     keyMap
	logger   *slog.Logger
	handlers *handlers.Registry

	shelf       shelf.Model
	selectedBin int

	width, height int
}

// New constructs a LoadedShelfState ready to receive a LoadShelfMsg
func New(svcs *services.AllServices, log *slog.Logger) LoadedShelfState {
	logGroup := log.WithGroup(states.LoadedShelf.String())

	return LoadedShelfState{
		svcs:   svcs,
		keys:   defaultKeybinds(),
		logger: logGroup,
		handlers: getHandlers(),
	}
}

func (s LoadedShelfState) Init() tea.Cmd {
	s.logger.Debug("loadedshelf state init")
	return nil
}

func (s LoadedShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if handler, ok := s.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(s, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		s = model.(LoadedShelfState)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	var shelfCmd tea.Cmd
	s.shelf, shelfCmd = util.UpdateModel(s.shelf, msg)
	cmds = append(cmds, shelfCmd)

	return s, tea.Batch(cmds...)
}

func (s LoadedShelfState) View() string {
	return s.renderModel()
}

func (s LoadedShelfState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}
