package loadcollection

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	discogs "github.com/dkaman/discogs-golang"
	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type refreshShelfMsg struct{}

type doneFetchingMsg struct{}

type loadNextMsg struct{}

// LoadCollectionFromDiscogsState holds the shelf model and renders it.
type LoadCollectionState struct {
	shelfService    *services.ShelfService
	nextState       states.StateType
	collection      shelf.Model
	discogsClient   *discogs.Client
	discogsUsername string
	layout          *layout.Div
	logger          *slog.Logger

	selectFolderForm *form

	spinner   spinner.Model
	fetching  bool
	inserting bool

	progressBar   progress.Model
	releases      []*record.Entity
	currentIndex  int
	totalReleases int
}

func New(s *services.ShelfService, l *layout.Div, log *slog.Logger, client *discogs.Client, username string) LoadCollectionState {
	logger := log.WithGroup("loadcollectionstate")

	f := newFolderSelectForm(client, username)

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return LoadCollectionState{
		shelfService:     s,
		nextState:        states.Undefined,
		discogsClient:    client,
		discogsUsername:  username,
		selectFolderForm: f,
		spinner:          sp,
		fetching:         false,
		inserting:        false,
		layout:           l,
		logger:           logger,
	}
}

// Init satisfies tea.Model.
func (s LoadCollectionState) Init() tea.Cmd {
	s.logger.Info("loadcollection state init called")
	return tea.Batch(
		s.collection.Init(),
		s.selectFolderForm.Init(),
		func() tea.Msg {
			return refreshShelfMsg{}
		},
	)
}

// Update handles incoming LoadCollectionMsg and updates the shelf model.
func (s LoadCollectionState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case refreshShelfMsg:
		s.logger.Info("loaded collection current shelf",
			slog.Any("shelf", s.shelfService.CurrentShelf),
		)
		s.collection = shelf.New(s.shelfService.CurrentShelf, s.logger)

		if s.selectFolderForm.State == huh.StateCompleted {
			s.selectFolderForm = newFolderSelectForm(s.discogsClient, s.discogsUsername)
		}

		return s, nil

	case tea.WindowSizeMsg:
		s.layout = newLoadedCollectionFormLayout(s.layout, s.selectFolderForm)

	case NewDiscogsCollectionMsg:
		s.logger.Info("new discogs collection to process",
			slog.Any("releases", msg.Releases),
		)

		s.releases = msg.Releases
		s.totalReleases = len(msg.Releases)
		s.currentIndex = 0

		s.fetching = true
		s.inserting = true

		cmds = append(cmds,
			func() tea.Msg { return loadNextMsg{} },
		)

		s.layout, _ = newLoadedCollectionProgressLayout(
			s.layout, s.progressBar, s.spinner, s.fetching, s.inserting, 0.0,
		)

		return s, tea.Batch(cmds...)

	case loadNextMsg:
		if s.releases != nil && s.currentIndex < s.totalReleases {
			s.logger.Info("inserting release",
				slog.Any("release", s.releases[s.currentIndex]),
			)
			phy := s.collection.PhysicalShelf()
			phy.Insert(s.releases[s.currentIndex])
			s.currentIndex++

			_ = s.shelfService.Shelves.Save(phy)
			s.shelfService.CurrentShelf = phy
			pct := float64(s.currentIndex) / float64(s.totalReleases)

			cmds = append(cmds,
				s.progressBar.SetPercent(pct),
				tea.Cmd(func() tea.Msg { return loadNextMsg{} }),
			)

			s.layout, _ = newLoadedCollectionProgressLayout(
				s.layout, s.progressBar, s.spinner, s.fetching, s.inserting, pct,
			)

			return s, tea.Batch(cmds...)
		}

		s.inserting = false
		s.fetching = false

		s.layout, _ = newLoadedCollectionProgressLayout(
			s.layout, s.progressBar, s.spinner, s.fetching, s.inserting, 1.0,
		)

		cmds = append(cmds, func() tea.Msg {
			return doneFetchingMsg{}
		})

		return s, tea.Batch(cmds...)

	case spinner.TickMsg:
		if s.fetching {
			var spinnerCmds tea.Cmd
			s.spinner, spinnerCmds = s.spinner.Update(msg)
			cmds = append(cmds, spinnerCmds)

			s.layout, _ = newLoadedCollectionProgressLayout(
				s.layout, s.progressBar, s.spinner, s.fetching, s.inserting, 0,
			)

			return s, tea.Batch(cmds...)
		}

	case doneFetchingMsg:
		cmds = append(cmds,
			tcmds.GetShelfCmd(s.shelfService.Shelves, s.collection.ID(), s.logger),
		)

		s.releases = nil
		s.nextState = states.LoadedShelf

		return s, tea.Batch(cmds...)

	default:
		// If we're not yet fetching, pass input to the folder‐select form
		if !s.fetching {
			fModel, formCmds := s.selectFolderForm.Update(msg)
			if f, ok := fModel.(*form); ok {
				s.selectFolderForm = f
			}
			cmds = append(cmds, formCmds)

			if s.selectFolderForm.State == huh.StateCompleted {
				folder := s.selectFolderForm.Folder()

				// Enter “fetching” mode:
				s.fetching = true

				s.progressBar = progress.New(progress.WithDefaultGradient())

				// Kick off the Discogs fetch
				cmds = append(cmds,
					s.progressBar.Init(),
					s.progressBar.SetPercent(0),
					s.spinner.Tick,
					RetrieveDiscogsCollection(s.discogsClient, s.discogsUsername, folder, s.logger),
				)

				s.layout, _ = newLoadedCollectionProgressLayout(
					s.layout, s.progressBar, s.spinner, s.fetching, s.inserting, 0,
				)

				s.selectFolderForm = newFolderSelectForm(s.discogsClient, s.discogsUsername)

				return s, tea.Batch(cmds...)
			}

			s.layout = newLoadedCollectionFormLayout(s.layout, s.selectFolderForm)

			// No further key handling while form is present
			return s, tea.Batch(cmds...)
		}
	}

	return s, tea.Batch(cmds...)
}

// View renders the shelf view into the TopWindow section.
func (s LoadCollectionState) View() string {
	return s.layout.Render()
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
