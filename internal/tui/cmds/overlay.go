package cmds

import (
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/models/overlay"
)

type ShowModalMsg struct {
	Modal overlay.TeaActiveCloser
}

type ClearModalMsg struct{}

func ShowModalCmd(modal overlay.TeaActiveCloser) tea.Cmd {
	return func() tea.Msg {
		return ShowModalMsg{Modal: modal}
	}
}

func ClearModalCmd() tea.Cmd {
	return func() tea.Msg {
		return ClearModalMsg{}
	}
}
