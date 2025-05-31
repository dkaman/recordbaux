package loadcollection

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	discogs "github.com/dkaman/discogs-golang"
)

type refreshShelfMsg struct{}

type doneFetchingMsg struct{}

type loadNextMsg struct{}

// LoadCollectionFromDiscogsState holds the shelf model and renders it.
type LoadCollectionState struct {
	app             *app.App
	nextState       statemachine.StateType
	collection      shelf.Model
	discogsClient   *discogs.Client
	discogsUsername string

	selectFolderForm *form

	spinner   spinner.Model
	fetching  bool
	inserting bool

	progressBar   progress.Model
	releases      []*physical.Record
	currentIndex  int
	totalReleases int

	layout *layout.Node
}

// New constructs the LoadCollectionFromDiscogs state with an empty shelf model.
func New(a *app.App, l *layout.Node, client *discogs.Client, username string) LoadCollectionState {
	c := shelf.New(nil, style.ActiveTextStyle)

	f := newFolderSelectForm(client, username)

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return LoadCollectionState{
		app:              a,
		nextState:        statemachine.Undefined,
		collection:       c,
		discogsClient:    client,
		discogsUsername:  username,
		selectFolderForm: f,
		spinner:          sp,
		fetching:         false,
		inserting:        false,
		layout:           l,
	}
}

// Init satisfies tea.Model.
func (s LoadCollectionState) Init() tea.Cmd {
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
		s.collection = s.app.CurrentShelf

		if s.selectFolderForm.State == huh.StateCompleted {
			s.selectFolderForm = newFolderSelectForm(s.discogsClient, s.discogsUsername)
		}

		return s, nil

	case NewDiscogsCollectionMsg:

		s.releases = msg.Releases
		s.totalReleases = len(msg.Releases)
		s.currentIndex = 0

		s.fetching = false
		s.inserting = true

		cmds = append(cmds,
			func() tea.Msg { return loadNextMsg{} },
		)

		s.layout, _ = newLoadedCollectionProgressLayout(
			s.layout, s.progressBar, s.spinner, s.fetching, s.inserting, 0,
		)

		return s, tea.Batch(cmds...)

	case loadNextMsg:
		if s.releases != nil && s.currentIndex < s.totalReleases {
			phy := s.collection.PhysicalShelf()
			phy.Insert(s.releases[s.currentIndex])
			s.currentIndex++

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
		s.releases = nil
		s.nextState = statemachine.LoadedShelf

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
					RetrieveDiscogsCollection(s.discogsClient, s.discogsUsername, folder),
				)

				s.layout, _ = newLoadedCollectionProgressLayout(
					s.layout, s.progressBar, s.spinner, s.fetching, s.inserting, 0,
				)

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
	return ""
}

func (s LoadCollectionState) Next() (statemachine.StateType, bool) {
	if s.nextState != statemachine.Undefined {
		return s.nextState, true
	}

	return statemachine.Undefined, false
}

func (s LoadCollectionState) Transition() {
	s.nextState = statemachine.Undefined
}
