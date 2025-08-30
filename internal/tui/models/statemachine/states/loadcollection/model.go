package loadcollection

import (
	"log/slog"

	huh "github.com/charmbracelet/huh/v2"
	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

	discogs "github.com/dkaman/discogs-golang"
	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type refreshShelfMsg struct{}

// LoadCollectionFromDiscogsState holds the shelf model and renders it.
type LoadCollectionState struct {
	shelfService    *services.ShelfService
	nextState       states.StateType
	discogsClient   *discogs.Client
	discogsUsername string
	logger          *slog.Logger

	selectFolderForm *form
	width, height    int
}

func New(s *services.ShelfService, log *slog.Logger, client *discogs.Client, username string) LoadCollectionState {
	logger := log.WithGroup("loadcollectionstate")

	f := newFolderSelectForm(client, username)

	return LoadCollectionState{
		shelfService:     s,
		nextState:        states.Undefined,
		discogsClient:    client,
		discogsUsername:  username,
		selectFolderForm: f,
		logger:           logger,
	}
}

// Init satisfies tea.Model.
func (s LoadCollectionState) Init() tea.Cmd {
	s.logger.Debug("loadcollection state init")
	return s.refresh()
}

func (s LoadCollectionState) refresh() tea.Cmd {
	return func() tea.Msg {
		return refreshShelfMsg{}
	}
}

// Update handles incoming LoadCollectionMsg and updates the shelf model.
func (s LoadCollectionState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case refreshShelfMsg:
		s.selectFolderForm = newFolderSelectForm(s.discogsClient, s.discogsUsername)
		cmds = append(cmds, s.selectFolderForm.Init())
		return s, tea.Batch(cmds...)

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

		// Kick off the Discogs fetch
		cmds = append(cmds,
			tcmds.RetrieveDiscogsCollection(s.discogsClient, s.discogsUsername, folder, s.logger),
		)

		s.nextState = states.FetchFromDiscogs

		return s, tea.Batch(cmds...)
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

func (s LoadCollectionState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}

	return states.Undefined, false
}

func (s LoadCollectionState) Transition() states.State {
	s.nextState = states.Undefined
	return s
}
