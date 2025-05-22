package statemachine

import (
	"fmt"
	"errors"

	tea "github.com/charmbracelet/bubbletea"

	teaCmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"
)

var (
	StateNotFoundErr = errors.New("state not found in state map")
)

type State interface {
	tea.Model
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

	switch msg := msg.(type) {
	case StateTransitionMsg:
		nextStateType := msg.NextState

		nextState, ok := m.allStates[nextStateType]
		if !ok {
			return m, tea.Quit
		}

		m.currentState = nextState
		m.currentStateType = nextStateType

		cmds = append(cmds,
			m.currentState.Init(),
			teaCmds.WithLayoutUpdate(layouts.StatusBar, fmt.Sprintf("state: %s", m.currentStateType)),
			teaCmds.WithLayoutUpdate(layouts.Overlay, ""),
		)

		return m, tea.Batch(cmds...)
	}

	stateModel, stateCmds := m.currentState.Update(msg)
	if s, ok := stateModel.(State); ok {
		m.allStates[m.currentStateType] = s
		m.currentState = s
	}

	cmds = append(cmds, stateCmds)

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
