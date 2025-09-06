package statemachine

import (
	"errors"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

	discogs "github.com/dkaman/discogs-golang"
	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	cps "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/createplaylist"
	css "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/createshelf"
	ffd "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/fetchfromdiscogs"
	lcs "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadcollection"
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
	currentState     states.State
	currentStateType states.StateType
	allStates        map[states.StateType]states.State

	logger        *slog.Logger
	width, height int
}

func New(svcs *services.AllServices, c *config.Config, log *slog.Logger) (Model, error) {
	logGroup := log.WithGroup("statemachine")

	m := Model{
		logger: logGroup,
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
		states.CreateShelf:      css.New(svcs, log),
		states.LoadedShelf:      lss.New(svcs, log),
		states.LoadCollection:   lcs.New(svcs, log, discogsClient, discogsUsername),
		states.LoadedBin:        lbs.New(svcs, log),
		states.FetchFromDiscogs: ffd.New(svcs, log, discogsClient),
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

	// transition
	switch t := msg.(type) {
	case tcmds.StateTransitionMsg:
		var pre tea.Cmd = nil
		if len(t.PreCmds) > 0 {
			pre = tea.Batch(t.PreCmds...)
		}

		// 2) A continuation Msg that performs the swap & AfterInit.
		type doSwapMsg struct{ tcmds.StateTransitionMsg }
		doSwap := func() tea.Msg { return tcmds.StateTransitionPostMsg{t} }

		// Execute: (PreInit...) -> doSwapMsg
		if pre != nil {
			return m, tea.Sequence(pre, doSwap)
		}
		return m, doSwap

	case tcmds.StateTransitionPostMsg:
		// Swap to the next state.
		next := t.Next

		m.logger.Info("state transition (envelope)",
			slog.String("from", m.currentStateType.String()),
			slog.String("to", next.String()),
		)

		// park the old state instance
		if s, ok := m.currentState.(states.State); ok {
			m.allStates[m.currentStateType] = s
		}

		m.currentState = m.allStates[next]
		m.currentStateType = next

		return m, tea.Sequence(
			m.currentState.Init(),
			tea.Batch(t.PostCmds...),
		)

	case tea.WindowSizeMsg:
		m.width, m.height = t.Width, t.Height
		m.logger.Debug("dimensions at statemachine, broadcasting to all states",
			slog.Any("msg", msg),
		)

		// Create a slice for all the commands that the child updates might return.
		allCmds := make([]tea.Cmd, 0, len(m.allStates))

		// park current state so it will be updated like the rest
		m.allStates[m.currentStateType] = m.currentState

		// Iterate over every state and send it the size message.
		for stateType, state := range m.allStates {
			// Call the state's Update method.
			updatedStateModel, cmd := state.Update(t)

			// The Update method returns a new model, so we must replace the old one in our map.
			if updatedState, ok := updatedStateModel.(states.State); ok {
				m.allStates[stateType] = updatedState
			}

			if cmd != nil {
				allCmds = append(allCmds, cmd)
			}
		}

		// re-read now updated current state
		m.currentState = m.allStates[m.currentStateType]

		return m, tea.Batch(allCmds...)
	}

	stateModel, stateCmds := m.currentState.Update(msg)
	if stateCmds != nil {
		cmds = append(cmds, stateCmds)
	}

	if s, ok := stateModel.(states.State); ok {
		m.currentState = s
	}

	return m, tea.Batch(cmds...)
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
