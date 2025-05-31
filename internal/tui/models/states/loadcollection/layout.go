package loadcollection

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

const (
	layoutViewport layout.Section = iota
)

var (
	viewportStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Margin(0)

	formStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)

	progressStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)

func newLoadedCollectionFormLayout(base *layout.Node, f *form) (*layout.Node, error) {
	r := &layout.TeaModelRenderer{
		Model: f,
		Style: formStyle,
	}

	base.AddSection(layoutViewport, r)
	base.SetJoinFunc(joinFunc)
	return base, nil
}

func newLoadedCollectionProgressLayout(base *layout.Node, m progress.Model) (*layout.Node, error) {
	r := &layout.TeaModelRenderer{
		Model: m,
		Style: progressStyle,
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
