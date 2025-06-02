package statemachine

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
)

type StateTransitionMsg struct {
	NextState states.StateType
}

func WithNextState(t states.StateType) tea.Cmd {
	return func() tea.Msg {
		return StateTransitionMsg{
			NextState: t,
		}
	}
}
