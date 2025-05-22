package statemachine

import (
	tea "github.com/charmbracelet/bubbletea"
)

type StateTransitionMsg struct {
	NextState StateType
}

func WithNextState(t StateType) tea.Cmd {
	return func() tea.Msg {
		return StateTransitionMsg{
			NextState: t,
		}
	}
}
