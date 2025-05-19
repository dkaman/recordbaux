package statemachine

import (
	tea "github.com/charmbracelet/bubbletea"
)

type stateTransitionMsg struct {
	nextState StateType
}

func WithNextState(target StateType) tea.Cmd {
	return func() tea.Msg {
		return stateTransitionMsg{
			nextState:  target,
		}
	}
}
