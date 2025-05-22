package cmds

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/style/layouts"
)

type LayoutUpdateMsg struct {
	Section layouts.Section
	Content string
}

func WithLayoutUpdate(s layouts.Section, c string) tea.Cmd {
	return func() tea.Msg {
		return LayoutUpdateMsg{
			Section: s,
			Content: c,
		}
	}
}
