package loadedshelf

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
)

const (
	layoutViewport layout.Section = iota
)

var (
	baseBinStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)

func newSelectShelfLayout(base *layout.Node, sh shelf.Model) (*layout.Node, error) {
	r := &layout.TeaModelRenderer{
		Model: sh,
		Style: baseBinStyle,
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
