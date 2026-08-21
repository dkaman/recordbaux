package cmds

import (
	tea "charm.land/bubbletea/v2"
)

type RefreshWindowSizeMsg struct {}

func RefreshWindowSizeCmd() tea.Cmd {
	return func() tea.Msg {
		return RefreshWindowSizeMsg{}
	}
}
