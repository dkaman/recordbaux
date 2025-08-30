package statemachine

import (
	tea "github.com/charmbracelet/bubbletea/v2"

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
