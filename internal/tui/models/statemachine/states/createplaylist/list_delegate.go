package createplaylist

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/v2/list"

	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/tui/models/track"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

var (
	boldStyle = lipgloss.NewStyle().
		Bold(true)

	dimStyle = lipgloss.NewStyle().
		Foreground(style.LightGrey)
)

type trackDelegate struct{}

func (d trackDelegate) Height() int                             { return 2 }
func (d trackDelegate) Spacing() int                            { return 1 }
func (d trackDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d trackDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	trackM, ok := listItem.(track.Model)
	if !ok {
		return
	}

	selectedMarker := "[ ] "
	if trackM.Selected {
		selectedMarker = "[" + boldStyle.Render("*") + "] "
	}

	pTrack := trackM.PhysicalTrack()

	titleLine := fmt.Sprintf("%s - %s", boldStyle.Render(pTrack.Title), pTrack.Duration)
	descLine := fmt.Sprintf("key: '%s' bpm: %d", pTrack.Key, pTrack.BPM)
	prefix := "  "
	if m.Index() == index {
		prefix = "> "
	}

	lines := lipgloss.JoinVertical(lipgloss.Left, titleLine, dimStyle.Render(descLine))
	display := lipgloss.JoinHorizontal(lipgloss.Top, prefix, selectedMarker, lines)

	sty := style.TextStyle
	if m.Index() == index {
		sty = style.ActiveTextStyle
	}

	fmt.Fprint(w, sty.Render(display))
}
