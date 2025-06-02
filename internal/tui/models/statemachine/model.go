package statemachine

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

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

type State interface {
	tea.Model
	Next() (states.StateType, bool)
	Transition()
	Help() string
}

type Model struct {
	currentState     State
	currentStateType states.StateType
	allStates        map[states.StateType]State
	layout           *div.Div
}

func New(a *app.App, c *config.Config, d *div.Div) (Model, error) {
	m := Model{
		layout: newStateMachineLayout(d),
	}

	discogsAPIKey, _ := c.String("shelf.discogs.key")
	discogsUsername, _ := c.String("shelf.discogs.username")

	discogsClient, err := discogs.New(
		discogs.WithToken(discogsAPIKey),
	)
	if err != nil {
		log.Printf("error in discogs client creation %w", err)
	}

	m.allStates = map[states.StateType]State{
		states.MainMenu:       mms.New(a, d),
		states.CreateShelf:    css.New(a, d),
		states.LoadedShelf:    lss.New(a, d),
		states.LoadCollection: lcs.New(a, d, discogsClient, discogsUsername),
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

	if s, okay := stateModel.(State); okay {
		if next, wanted := s.Next(); wanted {
			m.allStates[m.currentStateType] = m.currentState

			s.Transition()

			m.currentState = m.allStates[next]
			m.currentStateType = next

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
	return m.currentState.View()
}

func (m Model) Help() string {
	return m.currentState.Help()
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
