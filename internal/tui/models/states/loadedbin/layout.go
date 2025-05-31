package loadedbin

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

const (
	layoutViewport layout.Section = iota
)

var (
	baseTableStyle = lipgloss.NewStyle().
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center)
)

func newLoadedBinLayout(base *layout.Node, m table.Model) (*layout.Node, error) {
	r := &layout.TextRenderer{
		Body:  m.View(),
		Style: baseTableStyle,
	}

	base.AddSection(layoutViewport, r)
	base.SetJoinFunc(joinFunc)

	return base, nil
}

func joinFunc(sec map[layout.Section]layout.Renderer) string {
	var sections []string

	if viewPort, ok := sec[layoutViewport]; ok {
		renderedViewport := viewPort.Render()
		sections = append(sections, renderedViewport)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
