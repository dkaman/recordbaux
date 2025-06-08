package statemachine

import (
	"errors"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style/div"

	discogs "github.com/dkaman/discogs-golang"
	css "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/createshelf"
	lcs "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadcollection"
	lbs "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadedbin"
	lss "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadedshelf"
	mms "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/mainmenu"
	sss "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/selectshelf"
)

var (
	StateNotFoundErr = errors.New("state not found in state map")
)

type helper interface{
	Help() string
}

type State interface {
	tea.Model
	helper
	Next() (states.StateType, bool)
	Transition()
}

type Model struct {
	currentState     State
	currentStateType states.StateType
	allStates        map[states.StateType]State
	layout           *div.Div

	logger *slog.Logger
}

func New(a *app.App, c *config.Config, d *div.Div, log *slog.Logger) (Model, error) {
	m := Model{
		layout: newStateMachineLayout(d),
		logger: log,
	}

	discogsAPIKey, _ := c.String("shelf.discogs.key")
	discogsUsername, _ := c.String("shelf.discogs.username")
	discogsClient, err := discogs.New(
		discogs.WithToken(discogsAPIKey),
	)
	if err != nil {
		m.logger.Info("error in discogs client", slog.Any("errorMsg", err))
	}

	m.allStates = map[states.StateType]State{
		states.MainMenu:       mms.New(a, d),
		states.CreateShelf:    css.New(a, d, log),
		states.LoadedShelf:    lss.New(a, d, log),
		states.LoadCollection: lcs.New(a, d, log, discogsClient, discogsUsername),
		states.LoadedBin:      lbs.New(a, d),
		states.SelectShelf:    sss.New(a, d),
	}

	m.currentState = m.allStates[states.MainMenu]

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return m.currentState.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	stateModel, stateCmds := m.currentState.Update(msg)
	cmds = append(cmds, stateCmds)

	if s, ok := stateModel.(State); ok {
		if next, wanted := s.Next(); wanted {
			m.allStates[m.currentStateType] = m.currentState

			s.Transition()

			m.currentState = m.allStates[next]
			m.currentStateType = next

			vp :=  m.layout.Find("viewport")
			vp.ClearChildren()

			cmds = append(cmds,
				m.currentState.Init(),
			)

			return m, tea.Batch(cmds...)
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

func (m Model) State(t states.StateType) State {
	return m.allStates[t]
}

func (m Model) CurrentState() State {
	return m.currentState
}
func (m Model) CurrentStateType() states.StateType {
	return m.currentStateType
}
