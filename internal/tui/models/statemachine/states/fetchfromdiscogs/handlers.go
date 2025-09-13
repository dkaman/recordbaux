package fetchfromdiscogs

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/spinner"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	tshelf "github.com/dkaman/recordbaux/internal/tui/models/shelf"
)

func getHandlers() *handlers.Registry {
	r := handlers.NewRegistry()

	handlers.Register(r, handleTeaWindowSizeMsg)
	handlers.Register(r, handleLoadShelfMsg)
	handlers.Register(r, handleNewDiscogsCollectionMsg)
	handlers.Register(r, handleLoadNextMsg)
	handlers.Register(r, handleNewDiscogsEnrichRecordMsg)
	handlers.Register(r, handleShelfSavedMsg)
	handlers.Register(r, handleTickMsg)

	return r
}

func handleTeaWindowSizeMsg(s FetchFromDiscogsState, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.width, s.height = msg.Width, msg.Height
	return s, nil, nil
}

func handleLoadShelfMsg(s FetchFromDiscogsState, msg shelf.LoadShelfMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.shelf = msg.Phy
	return s, nil, nil
}

func handleNewDiscogsCollectionMsg(s FetchFromDiscogsState, msg tcmds.NewDiscogsCollectionMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.releases = msg.Releases
	s.totalReleases = len(msg.Releases)
	s.currentIndex = 0
	s.fetching = false
	s.pct = 0.0

	s.logger.Debug("new discogs collection to process",
		slog.Int("totalReleases", s.totalReleases),
		slog.Int("currentIndex", s.currentIndex),
	)

	return s, tea.Batch(s.progress.Init(), s.loadNextRecord()), nil
}

func handleLoadNextMsg(s FetchFromDiscogsState, msg loadNextMsg) (tea.Model, tea.Cmd, tea.Msg) {
	if s.currentIndex < s.totalReleases {
		s.logger.Debug("enriching release",
			slog.Any("release", s.releases[s.currentIndex].Title),
			slog.Int("index", s.currentIndex),
		)
		// Dispatch a command to fetch detailed track info.
		rec := s.releases[s.currentIndex]
		return s, tcmds.EnrichReleaseInstance(s.discogsClient, rec), nil
	}

	// Loop is finished. Finalize and prepare to transition state.
	s.logger.Debug("finished processing all releases")
	s.releases = nil

	return s, tcmds.Transition(
		states.LoadedShelf,
		nil,
		[]tea.Cmd{s.svcs.GetShelfCmd(s.shelf.ID)},
	), nil
}

func handleNewDiscogsEnrichRecordMsg(s FetchFromDiscogsState, msg tcmds.NewDiscogsEnrichRecordMsg) (tea.Model, tea.Cmd, tea.Msg) {
	if msg.Err != nil {
		s.logger.Error("failed to enrich record, skipping", slog.String("error", msg.Err.Error()))
		// Skip this record and move to the next one.
		s.currentIndex++
		return s, s.loadNextRecord(), nil
	}

	if s.currentIndex < len(s.releases) {
		s.currentTitle = msg.Record.Title
	}

	// Insert the hydrated record into the in-memory shelf object.
	s.shelf.Insert(msg.Record)
	s.logger.Debug("received record and inserted", slog.Any("record", msg.Record))

	// Dispatch a command to save the updated shelf to the database.
	return s, tea.Batch(
		s.svcs.SaveShelfCmd(s.shelf),
		tshelf.WithPhysicalShelf(s.shelf),
	), nil

}

func handleShelfSavedMsg(s FetchFromDiscogsState, msg services.ShelfSavedMsg) (tea.Model, tea.Cmd, tea.Msg) {
	if msg.Err != nil {
		s.logger.Error("failed to save shelf", slog.String("error", msg.Err.Error()))
		// Depending on desired behavior, you could stop or just log and continue.
	}

	// The save was successful, so now we can update the progress and process the next record.
	s.currentIndex++
	s.pct = float64(s.currentIndex) / float64(s.totalReleases)

	return s, tea.Batch(s.progress.SetPercent(s.pct), s.loadNextRecord()), nil
}

func handleTickMsg(s FetchFromDiscogsState, msg spinner.TickMsg) (tea.Model, tea.Cmd, tea.Msg) {
	if s.fetching {
		var spinnerCmd tea.Cmd
		s.spinner, spinnerCmd = s.spinner.Update(msg)
		return s, tea.Batch(spinnerCmd), nil
	}
	return s, nil, nil
}
