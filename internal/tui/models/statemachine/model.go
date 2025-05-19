package statemachine

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/style/layouts"
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
	case stateTransitionMsg:
		nextStateType := msg.nextState

		nextState, ok := m.allStates[nextStateType]
		if !ok {
			return m, tea.Quit
		}

		m.currentState = nextState
		m.currentStateType = nextStateType
		m.layout.WithSection(layouts.Overlay, "")

		cmds = append(cmds, m.currentState.Init())

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
func (m Model) CurrentStateType() string {
	switch m.currentStateType {
	case MainMenu:
		return "main menu"
	case CreateShelf:
		return "create shelf"
	case LoadedShelf:
		return "loaded shelf"
	}

	return "undefined"
}
