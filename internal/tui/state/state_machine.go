package state

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

type StateType int

const (
	// states
	MainMenu StateType = iota
	CreateShelf
	Quit
)

var (
	StateNotFoundErr = errors.New("state not found in state map")
)

type State interface {
	tea.Model
	Next(tea.Msg) (*StateType, error)
}

type Machine struct {
	currentState     State
	currentStateType StateType
	allStates        map[StateType]State
}

func NewMachine(initialState StateType, states map[StateType]State) (Machine, error) {
	s, ok := states[initialState]
	if !ok {
		return Machine{}, StateNotFoundErr

	}

	return Machine{
		currentState:     s,
		currentStateType: initialState,
		allStates:        states,
	}, nil
}

func (m Machine) Init() tea.Cmd {
	return m.currentState.Init()
}

func (m Machine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	prevType := m.currentStateType

	stateModel, stateCmds := m.currentState.Update(msg)
	if s, ok := stateModel.(State); ok {
		m.allStates[prevType] = s
		m.currentState = s
	}

	cmds = append(cmds, stateCmds)

	nextStateType, err := m.currentState.Next(msg)
	if err != nil {
		return m, tea.Quit
	}

	if *nextStateType != prevType {
		nextState, ok := m.allStates[*nextStateType]
		if !ok {
			return m, tea.Quit
		}

		m.currentState = nextState
		m.currentStateType = *nextStateType

		cmds = append(cmds, m.currentState.Init())
	}

	return m, tea.Batch(cmds...)
}

func (m Machine) View() string {
	return m.currentState.View()
}
