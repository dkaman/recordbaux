package statemachine

import (
	"errors"
	"log/slog"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/services"
	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
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

	currentState     states.State
	currentStateType states.StateType
	allStates        map[states.StateType]states.State

	width, height int
}

func New(svcs *services.AllServices, c *config.Config, log *slog.Logger) (Model, error) {
	logGroup := log.WithGroup("statemachine")

	m := Model{
		logger:   logGroup,
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tcmds.StateTransitionMsg:
		next := msg.Transition.Next

		nextState := m.allStates[next]

		sizeMsg := tea.WindowSizeMsg{
			Width: m.width,
			Height: m.height,
		}

		resizedNextState, sizeUpdateCmd := util.UpdateModel(nextState, sizeMsg)

		m.logger.Info("state transition",
			slog.String("from", m.currentStateType.String()),
			slog.String("to", next.String()),
		)

		m.allStates[m.currentStateType] = m.currentState
		m.currentState = resizedNextState
		m.currentStateType = next
		m.allStates[next] = resizedNextState

		// new state will be initialized and then post-transition commands will
		// run
		return m, tea.Sequence(
			m.currentState.Init(),
			sizeUpdateCmd,
			tea.Batch(msg.Transition.PostCmds...),
		)
	}

	var stateCmds tea.Cmd
	m.currentState, stateCmds = util.UpdateModel(m.currentState, msg)

	return m, tea.Batch(stateCmds)
}

func (m Model) View() tea.View {
	currentState := m.currentState.View()

	viewportStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height)

	content := viewportStyle.Render(currentState.Content)

	return tea.NewView(content)
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
