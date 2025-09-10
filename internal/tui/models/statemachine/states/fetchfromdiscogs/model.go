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
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/style"

	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

type loadNextMsg struct{}

type FetchFromDiscogsState struct {
	svcs          *services.AllServices
	logger        *slog.Logger
	discogsClient *discogs.Client
	handlers      *handlers.Registry

	spinner  spinner.Model
	progress progress.Model

	releases      []*record.Entity
	shelf         *shelf.Entity
	currentIndex  int
	totalReleases int
	fetching      bool
	pct           float64
	currentTitle  string

	width, height int
}

func New(svcs *services.AllServices, log *slog.Logger, d *discogs.Client) FetchFromDiscogsState {
	logGroup := log.WithGroup("fetchfromdiscogsstate")

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(style.LightMagenta)

	prg := progress.New(progress.WithDefaultGradient())

	return FetchFromDiscogsState{
		svcs:          svcs,
		logger:        logGroup,
		discogsClient: d,
		handlers:      getHandlers(),

		spinner:  sp,
		progress: prg,
		fetching: true,
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

	if handler, ok := s.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(s, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		s = model.(FetchFromDiscogsState)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s FetchFromDiscogsState) View() string {
	return s.renderModel()
}

func (s FetchFromDiscogsState) Help() string {
	return "fetching collection from discogs"
}
