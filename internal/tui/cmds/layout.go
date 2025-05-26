package cmds

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

type LayoutUpdateMsg struct {
	Section layout.Section
	Content string
}

func WithLayoutUpdate(s layout.Section, c string) tea.Cmd {
	return func() tea.Msg {
		return LayoutUpdateMsg{
			Section: s,
			Content: c,
		}
	}
}
