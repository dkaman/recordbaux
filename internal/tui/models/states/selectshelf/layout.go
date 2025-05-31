package selectshelf

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

const (
	layoutViewport layout.Section = iota
)

var (
	listStyle = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
)

func newSelectShelfLayout(base *layout.Node, m list.Model) (*layout.Node, error) {
	r := &layout.TextRenderer{
		Body:  m.View(),
		Style: listStyle,
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
