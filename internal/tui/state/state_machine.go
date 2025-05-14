package state

import (
	tea "github.com/charmbracelet/bubbletea"
)

type StateType int

const (
	Idle StateType = iota
	NewShelf
	LoadShelf
)

type State interface {
	tea.Model
	Next(tea.Msg) StateType
}

type Machine struct {
	currentState State
	allStates    map[StateType]State
}

func NewMachine(initialState State, states map[StateType]State) (Machine, error) {
	return Machine{
		currentState: initialState,
		allStates:    states,
	}, nil
}

func (m Machine) Init() tea.Cmd {
	return nil
}

func (m Machine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentState := m.currentState

	model, cmd := currentState.Update(msg)

	m.currentState = m.allStates[currentState.Next(msg)]

	return model, cmd
}

func (m Machine) View() string {
	return m.currentState.View()
}
