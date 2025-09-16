package loadedshelf

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/spinner"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

func getHandlers() *handlers.Registry {
	r := handlers.NewRegistry()

	handlers.Register(r, handleTeaWindowSizeMsg)
	handlers.Register(r, handleTeaKeyPressMsg)
	handlers.Register(r, handleLoadShelfMsg)
	handlers.Register(r, handleNewDiscogsCollectionMsg)
	handlers.Register(r, handleLoadNextMsg)
	handlers.Register(r, handleNewDiscogsEnrichRecordMsg)
	handlers.Register(r, handleShelfSavedMsg)
	handlers.Register(r, handleTickMsg)
	return r
}

func handleTeaWindowSizeMsg(s LoadedShelfState, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.width, s.height = msg.Width, msg.Height
	return s, nil, nil
}

func handleTeaKeyPressMsg(s LoadedShelfState, msg tea.KeyPressMsg) (tea.Model, tea.Cmd, tea.Msg) {
	sh := s.shelf.PhysicalShelf()
	if sh == nil {
		return s, nil, nil
	}

	if s.loading || s.fetching {
		return s, nil, msg
	}

	switch {
	case key.Matches(msg, s.keys.Next):
		s.shelf = s.shelf.SelectNextBin()

	case key.Matches(msg, s.keys.Prev):
		s.shelf = s.shelf.SelectPrevBin()

	case key.Matches(msg, s.keys.Back):
		return s, tcmds.Transition(states.MainMenu, nil, nil), nil

	case key.Matches(msg, s.keys.Load):
		s.loading = true
		s.shelf.Blur()
		s.loadCollectionForm = newFolderSelectForm(s.discogsClient, s.discogsUsername)
		return s, s.loadCollectionForm.Init(), nil

	case msg.String() == "enter":
		b := s.shelf.GetSelectedBin().PhysicalBin()
		return s, tcmds.Transition(
			states.LoadedBin,
			nil,
			[]tea.Cmd{bin.WithPhysicalBin(b)},
		), nil
	}

	return s, nil, msg
}

func handleLoadShelfMsg(s LoadedShelfState, msg shelf.LoadShelfMsg) (tea.Model, tea.Cmd, tea.Msg) {
	sh := msg.Phy

	s.shelf = shelf.New(sh, s.logger).
		SetSize(s.width, s.height).
		SelectBin(0)

	return s, nil, nil
}

func handleTickMsg(s LoadedShelfState, msg spinner.TickMsg) (tea.Model, tea.Cmd, tea.Msg) {
	if s.fetching {
		var cmd tea.Cmd
		s.spin, cmd = s.spin.Update(msg)
		return s, cmd, nil
	}
	return s, nil, msg
}

func handleNewDiscogsCollectionMsg(s LoadedShelfState, msg tcmds.NewDiscogsCollectionMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.releases = msg.Releases
	s.totalReleases = len(msg.Releases)
	s.currentIndex = 0
	s.pct = 0.0
	s.logger.Debug("new discogs collection to process", slog.Int("count", s.totalReleases))
	return s, s.loadNextRecord(), nil
}

func handleLoadNextMsg(s LoadedShelfState, msg loadNextMsg) (tea.Model, tea.Cmd, tea.Msg) {
	if s.currentIndex >= s.totalReleases {
		s.logger.Debug("finished processing all releases")
		s.fetching = false
		s.shelf.Focus()
		// Reload the shelf from the DB to get all relationships correctly
		return s, s.svcs.GetShelfCmd(s.shelf.ID()), nil
	}

	s.logger.Debug("enriching next release", slog.Int("index", s.currentIndex))
	rec := s.releases[s.currentIndex]
	return s, tcmds.EnrichReleaseInstance(s.discogsClient, rec), nil
}

func handleNewDiscogsEnrichRecordMsg(s LoadedShelfState, msg tcmds.NewDiscogsEnrichRecordMsg) (tea.Model, tea.Cmd, tea.Msg) {
	if msg.Err != nil {
		s.logger.Error("failed to enrich record, skipping", slog.String("error", msg.Err.Error()))
		s.currentIndex++
		return s, s.loadNextRecord(), nil
	}

	if s.currentIndex < len(s.releases) {
		s.currentTitle = msg.Record.Title
	}

	// This updates the in-memory representation
	_, err := s.shelf.PhysicalShelf().Insert(msg.Record)
	if err != nil {
		s.logger.Error("failed to insert record into shelf model", slog.String("error", err.Error()))
	}

	// Persist the change to the database
	return s, s.svcs.SaveShelfCmd(s.shelf.PhysicalShelf()), nil
}

func handleShelfSavedMsg(s LoadedShelfState, msg services.ShelfSavedMsg) (tea.Model, tea.Cmd, tea.Msg) {
	if msg.Err != nil {
		s.logger.Error("failed to save shelf", slog.String("error", msg.Err.Error()))
	}
	s.currentIndex++
	s.pct = float64(s.currentIndex) / float64(s.totalReleases)
	return s, tea.Batch(s.prog.SetPercent(s.pct), s.loadNextRecord()), nil
}
