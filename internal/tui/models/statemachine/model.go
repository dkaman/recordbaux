package statemachine

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	teaCmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"
)

var (
	StateNotFoundErr = errors.New("state not found in state map")
)

type State interface {
	tea.Model
	Next() (StateType, bool)
	Transition()
}

type Model struct {
	currentState     State
	currentStateType StateType
	allStates        map[StateType]State
}

func New(initialState StateType, states map[StateType]State) (Model, error) {
	s, ok := states[initialState]
	if !ok {
		return Model{}, StateNotFoundErr

	}

	return Model{
		currentState:     s,
		currentStateType: initialState,
		allStates:        states,
	}, nil
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
				teaCmds.WithLayoutUpdate(layouts.StatusBar, fmt.Sprintf("state: %s", m.currentStateType)),
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

func (m Model) State(t StateType) State {
	return m.allStates[t]
}

func (m Model) CurrentState() State {
	return m.currentState
}
func (m Model) CurrentStateType() StateType {
	return m.currentStateType
}
