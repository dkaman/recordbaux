package statemachine

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

type StateType int

const (
	// states
	MainMenu StateType = iota
	CreateShelf
	LoadedShelf
	Quit
)

var (
	StateNotFoundErr = errors.New("state not found in state map")
)

type State interface {
	tea.Model
	Next(tea.Msg) (*StateType, error)
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

	currType := m.currentStateType

	stateModel, stateCmds := m.currentState.Update(msg)
	if s, ok := stateModel.(State); ok {
		m.allStates[currType] = s
		m.currentState = s
	}

	cmds = append(cmds, stateCmds)

	nextStateType, err := m.currentState.Next(msg)
	if err != nil {
		return m, tea.Quit
	}

	if *nextStateType != currType {
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

func (m Model) View() string {
	return m.currentState.View()
}
