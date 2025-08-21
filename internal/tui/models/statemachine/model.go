package statemachine

import (
	"errors"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	discogs "github.com/dkaman/discogs-golang"
	css "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/createshelf"
	ffd "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/fetchfromdiscogs"
	lcs "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadcollection"
	lbs "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadedbin"
	lss "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadedshelf"
	mms "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/mainmenu"
	cps "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/createplaylist"
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
	layout           *layout.Div

	logger *slog.Logger
}

func New(s *services.ShelfService, t *services.TrackService, p *services.PlaylistService, c *config.Config, d *layout.Div, log *slog.Logger) (Model, error) {
	logGroup := log.WithGroup("statemachine")

	m := Model{
		layout: newStateMachineLayout(d),
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
		states.MainMenu:         mms.New(s, t, p, d, log),
		states.CreateShelf:      css.New(s, d, log),
		states.LoadedShelf:      lss.New(s, d, log),
		states.LoadCollection:   lcs.New(s, d, log, discogsClient, discogsUsername),
		states.LoadedBin:        lbs.New(s, d, log),
		states.FetchFromDiscogs: ffd.New(s, d, log, discogsClient),
		states.CreatePlaylist:   cps.New(s, t, p, d, log),
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

	stateModel, stateCmds := m.currentState.Update(msg)
	if stateCmds != nil {
		cmds = append(cmds, stateCmds)
	}

	if s, ok := stateModel.(states.State); ok {
		if next, wanted := s.Next(); wanted {
			m.logger.Info("state transition",
				slog.String("from", m.currentStateType.String()),
				slog.String("to", next.String()),
			)

			s = s.Transition()

			m.allStates[m.currentStateType] = s
			m.currentState = m.allStates[next]
			m.currentStateType = next

			vp := m.layout.Find("viewport")
			vp.ClearChildren()

			cmds = append(cmds,
				m.currentState.Init(),
			)

			return m, tea.Sequence(cmds...)
		}

		m.currentState = s
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return ""
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
