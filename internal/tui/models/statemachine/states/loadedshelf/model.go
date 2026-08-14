package loadedshelf

import (
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/dkaman/discogs-golang"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type LoadedShelfState struct {
	svcs   *services.AllServices
	keys   keyMap
	logger *slog.Logger

	// discogs info
	discogsClient   *discogs.Client
	discogsUsername string

	shelf       shelf.Model
	selectedBin int

	loading            bool
	fetching           bool
	loadCollectionForm *loadCollectionForm

	// Fields for fetching progress UI and state
	spin          spinner.Model
	prog          progress.Model
	releases      []*record.Entity
	currentIndex  int
	totalReleases int
	pct           float64
	currentTitle  string

	width, height int
}

type loadNextMsg struct{}

func (s LoadedShelfState) loadNextRecord() tea.Cmd {
	return func() tea.Msg {
		return loadNextMsg{}
	}
}

// New constructs a LoadedShelfState ready to receive a LoadShelfMsg
func New(svcs *services.AllServices, log *slog.Logger, c *discogs.Client, u string) LoadedShelfState {
	logGroup := log.WithGroup(states.LoadedShelf.String())

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(style.LightMagenta)

	prg := progress.New(progress.WithDefaultBlend())

	return LoadedShelfState{
		svcs:   svcs,
		keys:   defaultKeybinds(),
		logger: logGroup,

		discogsClient:   c,
		discogsUsername: u,

		fetching: false,
		loading:  false,

		spin: sp,
		prog: prg,
	}
}

func (s LoadedShelfState) Init() tea.Cmd {
	s.logger.Debug("loadedshelf state init")
	return nil
}

func (s LoadedShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var passthru tea.Msg

	passthru = msg

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height
		passthru = msg

	case tea.KeyPressMsg:
		sh := s.shelf.PhysicalShelf()
		if sh == nil {
			return s, nil
		}

		if s.loading || s.fetching {
			passthru = msg
		}

		switch {
		case key.Matches(msg, s.keys.Next):
			s.shelf = s.shelf.SelectNextBin()

		case key.Matches(msg, s.keys.Prev):
			s.shelf = s.shelf.SelectPrevBin()

		case key.Matches(msg, s.keys.Back):
			return s, tcmds.Transition(states.MainMenu, nil, nil)

		case key.Matches(msg, s.keys.Load):
			s.loading = true
			s.shelf.Blur()
			s.loadCollectionForm = newFolderSelectForm(s.discogsClient, s.discogsUsername)
			return s, s.loadCollectionForm.Init()

		case msg.String() == "enter":
			b := s.shelf.GetSelectedBin().PhysicalBin()
			return s, tcmds.Transition(
				states.LoadedBin,
				nil,
				[]tea.Cmd{bin.WithPhysicalBin(b)},
			)
		}

		passthru = msg


	case shelf.LoadShelfMsg:
		sh := msg.Phy

		s.shelf = shelf.New(sh, s.logger).
			SetSize(s.width, s.height).
			SelectBin(0)

		return s, nil

	case spinner.TickMsg:
		if s.fetching {
			var cmd tea.Cmd
			s.spin, cmd = s.spin.Update(msg)
			return s, cmd
		}
		passthru = msg

	case tcmds.NewDiscogsCollectionMsg:
		s.releases = msg.Releases
		s.totalReleases = len(msg.Releases)
		s.currentIndex = 0
		s.pct = 0.0
		s.logger.Debug("new discogs collection to process", slog.Int("count", s.totalReleases))
		return s, s.loadNextRecord()

	case loadNextMsg:
		if s.currentIndex >= s.totalReleases {
			s.logger.Debug("finished processing all releases")
			s.fetching = false
			s.shelf.Focus()
			// Reload the shelf from the DB to get all relationships correctly
			return s, s.svcs.GetShelfCmd(s.shelf.ID())
		}

		s.logger.Debug("enriching next release", slog.Int("index", s.currentIndex))
		rec := s.releases[s.currentIndex]
		return s, tcmds.EnrichReleaseInstance(s.discogsClient, rec)

	case tcmds.NewDiscogsEnrichRecordMsg:
		if msg.Err != nil {
			s.logger.Error("failed to enrich record, skipping", slog.String("error", msg.Err.Error()))
			s.currentIndex++
			return s, s.loadNextRecord()
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
		return s, s.svcs.SaveShelfCmd(s.shelf.PhysicalShelf())

	case services.ShelfSavedMsg:
		if msg.Err != nil {
			s.logger.Error("failed to save shelf", slog.String("error", msg.Err.Error()))
		}
		s.currentIndex++
		s.pct = float64(s.currentIndex) / float64(s.totalReleases)
		return s, tea.Batch(s.prog.SetPercent(s.pct), s.loadNextRecord())
	}

	// If a modal is active, it captures all updates first.
	if s.loading {
		var formCmd tea.Cmd
		s.loadCollectionForm, formCmd = util.UpdateModel(s.loadCollectionForm, passthru)
		cmds = append(cmds, formCmd)

		if s.loadCollectionForm.Form.State == huh.StateCompleted {
			folder := s.loadCollectionForm.Folder()
			s.logger.Debug("folder selected, starting fetch", slog.String("folder", folder))
			s.loading = false
			s.fetching = true
			cmds = append(cmds, s.spin.Tick, tcmds.RetrieveDiscogsCollection(s.discogsClient, s.discogsUsername, folder, s.logger))
		} else if s.loadCollectionForm.Form.State == huh.StateAborted {
			s.loading = false
			s.shelf.Focus()
		}
		return s, tea.Batch(cmds...)
	}

	var shelfCmd tea.Cmd
	s.shelf, shelfCmd = util.UpdateModel(s.shelf, passthru)
	cmds = append(cmds, shelfCmd)

	return s, tea.Batch(cmds...)
}

func (s LoadedShelfState) View() tea.View {
	return s.renderModel()
}

func (s LoadedShelfState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}
