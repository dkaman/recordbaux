package fetchfromdiscogs

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/progress"
	"github.com/charmbracelet/bubbles/v2/spinner"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/discogs-golang"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"

	lipgloss "github.com/charmbracelet/lipgloss/v2"
	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	tshelf "github.com/dkaman/recordbaux/internal/tui/models/shelf"
)

type loadNextMsg struct{}

type FetchFromDiscogsState struct {
	shelfService *services.ShelfService
	nextState    states.StateType
	logger       *slog.Logger

	spinner  spinner.Model
	progress progress.Model

	releases      []*record.Entity
	shelf         *shelf.Entity
	currentIndex  int
	totalReleases int
	fetching      bool

	discogsClient *discogs.Client
	width, height int
	pct           float64
	currentTitle  string
}

func New(s *services.ShelfService, log *slog.Logger, d *discogs.Client) FetchFromDiscogsState {
	logGroup := log.WithGroup("fetchfromdiscogsstate")

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(style.LightMagenta)

	prg := progress.New(progress.WithDefaultGradient())

	return FetchFromDiscogsState{
		shelfService:  s,
		nextState:     states.Undefined,
		logger:        logGroup,
		spinner:       sp,
		progress:      prg,
		fetching:      true,
		discogsClient: d,
	}
}

func (s FetchFromDiscogsState) Init() tea.Cmd {
	return nil
}

func (s FetchFromDiscogsState) loadNextRecord() tea.Cmd {
	return func() tea.Msg {
		return loadNextMsg{}
	}
}

func (s FetchFromDiscogsState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tshelf.LoadShelfMsg:
		s.shelf = msg.Phy
		return s, nil

	case tcmds.NewDiscogsCollectionMsg:
		s.releases = msg.Releases
		s.totalReleases = len(msg.Releases)
		s.currentIndex = 0
		s.fetching = false
		s.pct = 0.0

		s.logger.Debug("new discogs collection to process",
			slog.Int("totalReleases", s.totalReleases),
			slog.Int("currentIndex", s.currentIndex),
		)

		cmds = append(cmds,
			s.progress.Init(),
			s.loadNextRecord(),
		)

		return s, tea.Batch(cmds...)

	// Step 2: Triggered internally to process the current record.
	case loadNextMsg:
		if s.currentIndex < s.totalReleases {
			s.logger.Debug("enriching release",
				slog.Any("release", s.releases[s.currentIndex].Title),
				slog.Int("index", s.currentIndex),
			)
			// Dispatch a command to fetch detailed track info.
			rec := s.releases[s.currentIndex]
			return s, tcmds.EnrichReleaseInstance(s.discogsClient, rec)
		}

		// Loop is finished. Finalize and prepare to transition state.
		s.logger.Debug("finished processing all releases")
		s.releases = nil

		return s, tcmds.WithNextState(
			states.LoadedShelf,
			nil,
			[]tea.Cmd{tcmds.GetShelfCmd(s.shelfService.Shelves, s.shelf.ID, s.logger)},
		)

	// Step 3: Receives the fully hydrated record from the enrichment command.
	case tcmds.NewDiscogsEnrichRecordMsg:
		if msg.Err != nil {
			s.logger.Error("failed to enrich record, skipping", slog.String("error", msg.Err.Error()))
			// Skip this record and move to the next one.
			s.currentIndex++
			return s, s.loadNextRecord()
		}

		if s.currentIndex < len(s.releases) {
			s.currentTitle = msg.Record.Title
		}

		// Insert the hydrated record into the in-memory shelf object.
		s.shelf.Insert(msg.Record)
		s.logger.Debug("received record and inserted", slog.Any("record", msg.Record))

		// Dispatch a command to save the updated shelf to the database.
		return s, tea.Batch(
			tcmds.SaveShelfCmd(s.shelfService.Shelves, s.shelf, s.logger),
			tshelf.WithPhysicalShelf(s.shelf),
		)

	// Step 4: Receives confirmation that the shelf was saved.
	case tcmds.ShelfSavedMsg:
		if msg.Err != nil {
			s.logger.Error("failed to save shelf", slog.String("error", msg.Err.Error()))
			// Depending on desired behavior, you could stop or just log and continue.
		}

		// The save was successful, so now we can update the progress and process the next record.
		s.currentIndex++
		s.pct = float64(s.currentIndex) / float64(s.totalReleases)

		cmds = append(cmds,
			s.progress.SetPercent(s.pct),
			s.loadNextRecord(),
		)
		return s, tea.Batch(cmds...)

	case spinner.TickMsg:
		if s.fetching {
			var spinnerCmds tea.Cmd

			s.spinner, spinnerCmds = s.spinner.Update(msg)
			cmds = append(cmds, spinnerCmds)

			return s, tea.Batch(cmds...)
		}
	}

	return s, tea.Batch(cmds...)
}

func (s FetchFromDiscogsState) View() string {
	return s.renderModel()
}

func (s FetchFromDiscogsState) Help() string {
	return "fetching collection from discogs"
}

func (s FetchFromDiscogsState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}

	return states.Undefined, false
}

func (s FetchFromDiscogsState) Transition() states.State {
	s.nextState = states.Undefined
	return s
}
