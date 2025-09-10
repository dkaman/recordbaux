package cmds

import (
	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
)

type envelope struct {
	Next     states.StateType
	PostCmds []tea.Cmd
}

type StateTransitionMsg struct {
	Transition envelope
}

func Transition(t states.StateType, before []tea.Cmd, after []tea.Cmd) tea.Cmd {
	var pre tea.Cmd = nil

	if len(before) > 0 {
		pre = tea.Batch(before...)
	}

	// lambda command that emits the actual transition message so that we
	// can sequence against the precmds before the current state is parked
	swap := func() tea.Msg {
		e := envelope{
			Next: t,
			PostCmds: after,
		}

		return StateTransitionMsg{
			Transition:  e,
		}
	}

	if pre != nil {
		return tea.Sequence(pre, swap)
	}

	return swap
}
