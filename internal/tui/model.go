package tui

import (
	"errors"
	"fmt"
	"log/slog"

	"charm.land/bubbles/v2/key"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/dkaman/discogs-golang"
	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/services"
	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"
)

const (
	CONF_DISCOGS_KEY  = "discogs.key"
	CONF_DISCOGS_USER = "discogs.username"
)

var (
	LoggerIsNilErr = errors.New("supplied slog logger is nil")
)

// root tui model
type Model struct {
	width, height int
	cfg           *config.Config
	svcs          *services.AllServices
	keys          keyMap
	logger        *slog.Logger

	discogsClient   *discogs.Client
	discogsUsername string

	stateMachine  statemachine.Model
	ready         bool
	topBarText    string
	statusBarText string
	helpVisible   bool
}

func New(c *config.Config, log *slog.Logger, svcs *services.AllServices) (Model, error) {
	m := Model{
		cfg:         c,
		svcs:        svcs,
		keys:        defaultKeybinds(),
		ready:       false,
		helpVisible: false,
	}

	if log == nil {
		return m, LoggerIsNilErr
	}

	sm, err := statemachine.New(log)
	if err != nil {
		return m, fmt.Errorf("error creating state machine: %w", err)
	}

	discogsAPIKey := c.String(CONF_DISCOGS_KEY)
	discogsUsername := c.String(CONF_DISCOGS_USER)
	discogsClient, err := discogs.New(
		discogs.WithToken(discogsAPIKey),
	)
	if err != nil {
		return m, err
	}

	m.stateMachine = sm
	m.logger = log.WithGroup("root")
	m.topBarText = "recordbaux - organize your record collection"
	m.statusBarText = fmt.Sprintf("current state: %s", sm.CurrentState().Type().String())
	m.discogsClient = discogsClient
	m.discogsUsername = discogsUsername

	return m, nil
}

func (m Model) Init() tea.Cmd {
	m.logger.Debug("root tui model init called")
	return m.stateMachine.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var passthruMsg tea.Msg

	m.logger.Info("event received",
		slog.Any("event", fmt.Sprintf("%#v", msg)),
	)

	passthruMsg = msg

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// this code updates the root element's size, then constructs a
		// size msg that represents the whole size of the viewport to
		// pass thru to child model updates
		m.width, m.height = msg.Width, msg.Height

		if !m.ready {
			m.ready = true
		}

		numBars := 2
		if m.helpVisible {
			numBars = 3
		}

		// modify window size msg and pass thru
		passthruMsg = tea.WindowSizeMsg{
			Width:  m.width - 2,
			Height: m.height - numBars - 2,
		}

	case tea.KeyPressMsg:
		// this branch also passes through the key message if there is
		// no match so the child models can handle the keys
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.ToggleHelp):
			m.helpVisible = !m.helpVisible
			return m, func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
			}
		}

	case tcmds.RefreshWindowSizeMsg:
		return m, m.refreshWindowSize()

	// handle intent messages
	case tcmds.GetDiscogsFoldersIntentMsg:
		return m, m.getDiscogsFolders()

	case tcmds.NewDiscogsCollectionIntentMsg:
		return m, m.retrieveDiscogsCollection(msg.Folder)

	case tcmds.NewDiscogsEnrichRecordIntentMsg:
		return m, m.enrichReleaseInstance(msg.Record)

	case tcmds.BinsLoadedIntentMsg:
		return m, m.getBinCmd(msg.ID)

	case tcmds.ShelfLoadIntentMsg:
		return m, m.getShelfCmd(msg.ID)

	case tcmds.ShelvesLoadIntentMsg:
		return m, m.getAllShelvesCmd()

	case tcmds.ShelfSaveIntentMsg:
		return m, m.saveShelfCmd(msg.Shelf)

	case tcmds.ShelfDeleteIntentMsg:
		return m, m.deleteShelfCmd(msg.ID)

	case tcmds.ShelfAllTracksIntentMsg:
		return m, m.getAllTracksFromShelfCmd(msg.ID)

	case tcmds.PlaylistLoadIntentMsg:
		return m, m.getPlaylistCmd(msg.ID)

	case tcmds.PlaylistsLoadIntentMsg:
		return m, m.getAllPlaylistsCmd()

	case tcmds.PlaylistSaveIntentMsg:
		return m, m.savePlaylistCmd(msg.Entity)

	case tcmds.PlaylistDeleteIntentMsg:
		return m, m.deletePlaylistsCmd(msg.ID)

	case tcmds.PlaylistCheckoutIntentMsg:
		return m, m.setCheckoutCmd(msg.Playlist, msg.Status)
	}

	var stateMachineCmd tea.Cmd
	m.stateMachine, stateMachineCmd = util.UpdateModel(m.stateMachine, passthruMsg)
	cmds = append(cmds, stateMachineCmd)

	m.statusBarText = fmt.Sprintf("current state: %s", m.stateMachine.CurrentState().Type())

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("\n initializing...")
	}

	numBars := 2
	if m.helpVisible {
		numBars = 3
	}

	barStyle := style.BarStyle.
		Width(m.width).
		Height(1)

	helpStyle := style.HelpBarStyle.
		Width(m.width).
		Height(1)

	viewportStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width).
		Height(max(0, m.height-numBars))

	topBar := barStyle.Render(m.topBarText)
	statusBar := barStyle.Render(m.statusBarText)
	viewPort := viewportStyle.Render(m.stateMachine.View().Content)

	var content string

	if m.helpVisible {
		helpBar := helpStyle.Render(m.Help())
		content = lipgloss.JoinVertical(lipgloss.Left,
			topBar,
			viewPort,
			helpBar,
			statusBar,
		)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Left,
			topBar,
			viewPort,
			statusBar,
		)
	}

	return tea.NewView(content)
}

func (m Model) Help() string {
	return "global[ " +
		util.FmtKeymap(m.keys.ShortHelp()) + "] " +
		m.stateMachine.Help()
}
