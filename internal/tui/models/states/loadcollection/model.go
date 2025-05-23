package loadcollection

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"

	discogs "github.com/dkaman/discogs-golang"
	teaCmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type refreshShelfMsg struct{}

type loadNextMsg struct{}

// LoadCollectionFromDiscogsState holds the shelf model and renders it.
type LoadCollectionState struct {
	app              *app.App
	nextState        statemachine.StateType
	collection       shelf.Model
	discogsClient    *discogs.Client
	discogsUsername  string
	selectFolderForm *form

	progressBar   progress.Model
	releases      []*physical.Record
	currentIndex  int
	totalReleases int
}

// New constructs the LoadCollectionFromDiscogs state with an empty shelf model.
func New(a *app.App, client *discogs.Client, username string) LoadCollectionState {
	c := shelf.New(nil, style.ActiveTextStyle)

	f := newFolderSelectForm(client, username)

	return LoadCollectionState{
		app:              a,
		nextState:        statemachine.Undefined,
		collection:       c,
		discogsClient:    client,
		discogsUsername:  username,
		selectFolderForm: f,
	}
}

// Init satisfies tea.Model.
func (s LoadCollectionState) Init() tea.Cmd {
	return tea.Batch(
		s.collection.Init(),
		s.selectFolderForm.Init(),
		s.progressBar.Init(),
		func() tea.Msg {
			return refreshShelfMsg{}
		},
		func() tea.Msg {
			return loadNextMsg{}
		},
	)
}

// Update handles incoming LoadCollectionMsg and updates the shelf model.
func (s LoadCollectionState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case refreshShelfMsg:
		s.collection = s.app.CurrentShelf
		return s, teaCmds.WithLayoutUpdate(layouts.Overlay, s.selectFolderForm.View())

	case NewDiscogsCollectionMsg:
		s.releases = msg.Releases
		s.totalReleases = len(msg.Releases)
		s.currentIndex = 0

		s.progressBar = progress.New(progress.WithDefaultGradient())

		cmds = append(cmds,
			s.progressBar.Init(),
			s.progressBar.SetPercent(0),
			func() tea.Msg { return loadNextMsg{} },
		)

		return s, tea.Batch(cmds...)

	case progress.FrameMsg:
		barModel, barUpdateCmds := s.progressBar.Update(msg)
		if bar, ok := barModel.(progress.Model); ok {
			s.progressBar = bar
		}
		cmds = append(cmds, barUpdateCmds)
		return s, tea.Batch(cmds...)

	case loadNextMsg:
		if s.releases != nil {
			if s.currentIndex < s.totalReleases {
				phy := s.collection.PhysicalShelf()
				phy.Insert(s.releases[s.currentIndex])
				s.currentIndex++

				pct := float64(s.currentIndex) / float64(s.totalReleases)

				cmds = append(cmds,
					s.progressBar.SetPercent(pct),
					teaCmds.WithLayoutUpdate(layouts.Overlay, s.progressBar.ViewAs(pct)),
					tea.Cmd(func() tea.Msg { return loadNextMsg{}}),
				)

				return s, tea.Batch(cmds...)
			} else {
				cmds = append(cmds,
					teaCmds.WithLayoutUpdate(layouts.Overlay, ""),
				)

				s.releases = nil
				s.nextState = statemachine.LoadedShelf

				return s, tea.Batch(cmds...)
			}
		}
	}

	fModel, formUpdateCmds := s.selectFolderForm.Update(msg)
	if f, ok := fModel.(*form); ok {
		s.selectFolderForm = f
	}
	cmds = append(cmds, formUpdateCmds)

	if s.selectFolderForm.State == huh.StateCompleted {
		fol := s.selectFolderForm.Folder()
		s.selectFolderForm = newFolderSelectForm(s.discogsClient, s.discogsUsername)
		cmds = append(cmds,
			RetrieveDiscogsCollection(s.discogsClient, s.discogsUsername, fol),
		)
	}

	cmds = append(cmds,
		teaCmds.WithLayoutUpdate(layouts.Overlay, s.selectFolderForm.View()),
	)

	cModel, cCmds := s.collection.Update(msg)
	if c, ok := cModel.(shelf.Model); ok {
		s.collection = c
	}

	cmds = append(cmds, cCmds)

	return s, tea.Batch(cmds...)
}

// View renders the shelf view into the TopWindow section.
func (s LoadCollectionState) View() string {
	view := s.selectFolderForm.View()
	return view
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
