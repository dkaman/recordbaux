package statemachine

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

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

	layout *layouts.TallLayout
}

func New(initialState StateType, states map[StateType]State, l *layouts.TallLayout) (Model, error) {
	s, ok := states[initialState]
	if !ok {
		return Model{}, StateNotFoundErr

	}

	return Model{
		currentState:     s,
		currentStateType: initialState,
		allStates:        states,
		layout: l,
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
		m.layout.WithSection(layouts.Overlay, "")

		cmds = append(cmds, m.currentState.Init())

		return m, tea.Batch(cmds...)
	case BroadcastLoadShelfMsg:
		var broadcastCmds []tea.Cmd

		for t, state := range m.allStates {
			sModel, sCmds := state.Update(msg)
			if s, ok := sModel.(State); ok {
				m.allStates[t] = s
			}
			broadcastCmds = append(broadcastCmds, sCmds)
		}
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
	statusLine := fmt.Sprintf("state: %s", m.CurrentStateType())
	m.layout.WithSection(layouts.StatusLine, statusLine)
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
