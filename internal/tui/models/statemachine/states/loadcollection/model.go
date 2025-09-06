package loadcollection

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"

	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

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
	shelf           *shelf.Entity

	selectFolderForm *form
	width, height    int
}

func New(svcs *services.AllServices, log *slog.Logger, client *discogs.Client, username string) LoadCollectionState {
	logger := log.WithGroup("loadcollectionstate")

	f := newFolderSelectForm(client, username)

	return LoadCollectionState{
		svcs:             svcs,
		discogsClient:    client,
		discogsUsername:  username,
		selectFolderForm: f,
		logger:           logger,
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

	switch msg := msg.(type) {
	case tshelf.LoadShelfMsg:
		s.shelf = msg.Phy
		return s, nil

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil
	}

	fModel, formCmds := s.selectFolderForm.Update(msg)
	if f, ok := fModel.(*form); ok {
		s.selectFolderForm = f
	}
	cmds = append(cmds, formCmds)

	if s.selectFolderForm.State == huh.StateCompleted {
		folder := s.selectFolderForm.Folder()

		s.logger.Debug("folder selected with form",
			slog.Any("folder", folder),
		)

		return s, tcmds.WithNextState(
			states.FetchFromDiscogs,
			nil,
			[]tea.Cmd{
				tshelf.WithPhysicalShelf(s.shelf),
				tcmds.RetrieveDiscogsCollection(s.discogsClient, s.discogsUsername, folder, s.logger),
			},
		)
	}

	// No further key handling while form is present
	return s, tea.Batch(cmds...)
}

// View renders the shelf view into the TopWindow section.
func (s LoadCollectionState) View() string {
	return s.renderModel()
}

func (s LoadCollectionState) Help() string {
	return "select a discogs folder to load into this shelf..."
}
