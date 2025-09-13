package statemachine

import (
	"errors"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/util"

	discogs "github.com/dkaman/discogs-golang"
	cps "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/createplaylist"
	lbs "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadedbin"
	lps "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadedplaylist"
	lss "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadedshelf"
	mms "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/mainmenu"
)

var (
	StateNotFoundErr = errors.New("state not found in state map")
)

const (
	ConfDiscogsKey  = "discogs.key"
	ConfDiscogsUser = "discogs.username"
)

type Model struct {
	logger   *slog.Logger
	handlers *handlers.Registry

	currentState     states.State
	currentStateType states.StateType
	allStates        map[states.StateType]states.State

	width, height int
}

func New(svcs *services.AllServices, c *config.Config, log *slog.Logger) (Model, error) {
	logGroup := log.WithGroup("statemachine")

	m := Model{
		logger:   logGroup,
		handlers: getHandlers(),
	}

	discogsAPIKey := c.String(ConfDiscogsKey)
	discogsUsername := c.String(ConfDiscogsUser)
	discogsClient, err := discogs.New(
		discogs.WithToken(discogsAPIKey),
	)
	if err != nil {
		return m, err
	}

	m.allStates = map[states.StateType]states.State{
		states.MainMenu:         mms.New(svcs, log),
		states.LoadedShelf:      lss.New(svcs, log, discogsClient, discogsUsername),
		states.LoadedBin:        lbs.New(svcs, log),
		states.CreatePlaylist:   cps.New(svcs, log),
		states.LoadedPlaylist:   lps.New(svcs, log),
	}

	m.currentState = m.allStates[states.MainMenu]

	return m, nil
}

func (m Model) Init() tea.Cmd {
	m.logger.Debug("statemachine init",
		slog.String("currentState", m.currentStateType.String()),
	)

	return m.currentState.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if handler, ok := m.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(m, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		m = model.(Model)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	var stateCmds tea.Cmd
	m.currentState, stateCmds = util.UpdateModel(m.currentState, msg)

	return m, tea.Batch(stateCmds)
}

func (m Model) View() string {
	return m.renderModel()
}

func (m Model) Help() string {
	return "statemachine: " + m.currentState.Help()
}

func (m Model) State(t states.StateType) states.State {
	return m.allStates[t]
}

func (m Model) CurrentState() states.State {
	return m.currentState
}
func (m Model) CurrentStateType() states.StateType {
	return m.currentStateType
}
