package mainmenu

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

const (
	layoutViewport layout.Section = iota
	layoutTestViewportText
)

var (
	viewportStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Margin(0)
)

func newMainMenuLayout(base *layout.Node) (*layout.Node, error) {
	base.AddSection(layoutTestViewportText, &layout.TextRenderer{
		Body:  "welcome to recordbaux",
		Style: viewportStyle,
	})

	base.SetJoinFunc(joinFunc)

	return base, nil
}

func joinFunc(sec map[layout.Section]layout.Renderer) string {
	var sections []string

	if viewPort, ok := sec[layoutTestViewportText]; ok {
		renderedViewport := viewPort.Render()
		sections = append(sections, renderedViewport)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
