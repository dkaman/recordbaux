package cmds

import (
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db/bin"
)

type BinsLoadedIntentMsg struct {
	ID uint
}

type BinsLoadedMsg struct {
	Bins []*bin.Entity
	Err  error
}

func GetBinCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		return BinsLoadedIntentMsg{ID: id}
	}
}
