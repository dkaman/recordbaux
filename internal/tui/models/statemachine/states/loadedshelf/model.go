package loadedshelf

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/progress"
	"github.com/charmbracelet/bubbles/v2/spinner"

	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"
	huh "github.com/charmbracelet/huh/v2"

	"github.com/dkaman/discogs-golang"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type LoadedShelfState struct {
	svcs     *services.AllServices
	keys     keyMap
	logger   *slog.Logger
	handlers *handlers.Registry

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

	prg := progress.New(progress.WithDefaultGradient())

	return LoadedShelfState{
		svcs:     svcs,
		keys:     defaultKeybinds(),
		logger:   logGroup,
		handlers: getHandlers(),

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

	// If a modal is active, it captures all updates first.
	if s.loading {
		var formCmd tea.Cmd
		s.loadCollectionForm, formCmd = util.UpdateModel(s.loadCollectionForm, msg)
		cmds = append(cmds, formCmd)

		if s.loadCollectionForm.Form.State == huh.StateCompleted {
			folder := s.loadCollectionForm.Folder()
			s.logger.Debug("folder selected, starting fetch", slog.String("folder", folder))
			s.loading = false
			s.fetching = true
			cmds = append(cmds, s.spin.Tick, tcmds.RetrieveDiscogsCollection(s.discogsClient, s.discogsUsername, folder, s.logger))
		} else if s.loadCollectionForm.Form.State == huh.StateAborted {
			s.loading = false
		}
		return s, tea.Batch(cmds...)
	}

	if handler, ok := s.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(s, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		s = model.(LoadedShelfState)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	var shelfCmd tea.Cmd
	s.shelf, shelfCmd = util.UpdateModel(s.shelf, msg)
	cmds = append(cmds, shelfCmd)

	return s, tea.Batch(cmds...)
}

func (s LoadedShelfState) View() string {
	return s.renderModel()
}

func (s LoadedShelfState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}
