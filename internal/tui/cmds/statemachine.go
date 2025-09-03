package cmds

import (
	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
)

type StateTransitionMsg struct {
	Next     states.StateType
	PreCmds  []tea.Cmd
	PostCmds []tea.Cmd
}

type StateTransitionPostMsg struct {
	StateTransitionMsg
}

func WithNextState(t states.StateType, before []tea.Cmd, after []tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return StateTransitionMsg{
			Next:     t,
			PreCmds:  before,
			PostCmds: after,
		}
	}
}
