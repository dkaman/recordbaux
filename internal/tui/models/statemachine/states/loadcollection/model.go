package loadcollection

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"

	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/util"

	discogs "github.com/dkaman/discogs-golang"
	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	tshelf "github.com/dkaman/recordbaux/internal/tui/models/shelf"
)

// LoadCollectionFromDiscogsState holds the shelf model and renders it.
type LoadCollectionState struct {
	svcs            *services.AllServices
	discogsClient   *discogs.Client
	discogsUsername string
	logger          *slog.Logger
	handlers        *handlers.Registry

	shelf            *shelf.Entity
	selectFolderForm *form

	width, height int
}

func New(svcs *services.AllServices, log *slog.Logger, client *discogs.Client, username string) LoadCollectionState {
	logger := log.WithGroup("loadcollectionstate")

	f := newFolderSelectForm(client, username)

	return LoadCollectionState{
		svcs:            svcs,
		discogsClient:   client,
		discogsUsername: username,
		logger:          logger,
		handlers:        getHandlers(),

		selectFolderForm: f,
	}
}

// Init satisfies tea.Model.
func (s LoadCollectionState) Init() tea.Cmd {
	s.logger.Debug("loadcollection state init")
	return s.selectFolderForm.Init()
}

// Update handles incoming LoadCollectionMsg and updates the shelf model.
func (s LoadCollectionState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if handler, ok := s.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(s, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		s = model.(LoadCollectionState)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	var formCmd tea.Cmd
	s.selectFolderForm, formCmd = util.UpdateModel(s.selectFolderForm, msg)

	if s.selectFolderForm.Form.State == huh.StateCompleted {
		folder := s.selectFolderForm.Folder()

		s.logger.Debug("folder selected with form",
			slog.Any("folder", folder),
		)

		return s, tcmds.Transition(
			states.FetchFromDiscogs,
			nil,
			[]tea.Cmd{
				tshelf.WithPhysicalShelf(s.shelf),
				tcmds.RetrieveDiscogsCollection(s.discogsClient, s.discogsUsername, folder, s.logger),
			},
		)
	}

	cmds = append(cmds, formCmd)

	return s, tea.Batch(cmds...)
}

// View renders the shelf view into the TopWindow section.
func (s LoadCollectionState) View() string {
	return s.renderModel()
}

func (s LoadCollectionState) Help() string {
	return "select a discogs folder to load into this shelf..."
}
