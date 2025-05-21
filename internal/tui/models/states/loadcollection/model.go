package loadcollection

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"

	discogs "github.com/dkaman/discogs-golang"
)

// LoadCollectionFromDiscogsState holds the shelf model and renders it.
type LoadCollectionState struct {
	collection       shelf.Model
	layout           *layouts.TallLayout
	discogsClient    *discogs.Client
	discogsUsername  string
	selectFolderForm *form
}

// New constructs the LoadCollectionFromDiscogs state with an empty shelf model.
func New(l *layouts.TallLayout, client *discogs.Client, username string) LoadCollectionState {
	c := shelf.New(nil, style.ActiveTextStyle)

	f := newFolderSelectForm(client, username)

	return LoadCollectionState{
		collection:       c,
		layout:           l,
		discogsClient:    client,
		discogsUsername:  username,
		selectFolderForm: f,
	}
}

// Init satisfies tea.Model.
func (s LoadCollectionState) Init() tea.Cmd {
	return tea.Batch(s.collection.Init(), s.selectFolderForm.Init())
}

// Update handles incoming LoadCollectionMsg and updates the shelf model.
func (s LoadCollectionState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case statemachine.LoadShelfMsg:
		sh := msg.Shelf
		s.collection = shelf.New(sh, style.ActiveTextStyle)
	case NewDiscogsCollectionMsg:
		sh := s.collection.PhysicalShelf()

		for _, rel := range msg.Releases {
			sh.Insert(rel)
		}

		cmds = append(cmds,
			statemachine.WithNextState(statemachine.LoadedShelf),
		)
	default:
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
	}

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
	s.layout.WithSection(layouts.Overlay, view)
	return view
}
