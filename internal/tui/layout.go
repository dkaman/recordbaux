package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

const (
	layoutTopBar layout.Section = iota
	layoutViewport
	layoutHelpBar
	layoutStatusBar
	layoutOverlay
	layoutTestTopBarText
	layoutTestStatusBarText
	layoutTestViewportText
)

var (
	barStyle = lipgloss.NewStyle().
			Background(style.DarkBlue).
			Foreground(style.DarkBlack).
			Bold(true)

	viewportStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Margin(0)
)

func newTUILayout(base *layout.Node) (*layout.Node, error) {
	b := base

	topBar, err := layout.New(b.Width, 1)
	if err != nil {
		return nil, err
	}

	topBar.AddSection(layoutTestTopBarText, &layout.TextRenderer{
		Body:  "recordbaux - organize your record collection",
		Style: barStyle,
	})

	statusBar, err := layout.New(b.Width, 1)
	if err != nil {
		return nil, err
	}

	statusBar.AddSection(layoutTestStatusBarText, &layout.TextRenderer{
		Body:  "status bar test",
		Style: barStyle,
	})

	vpWidth := b.Width - 2
	vpHeight := b.Height - 2

	viewport, err := layout.New(vpWidth, vpHeight-2)
	if err != nil {
		return nil, err
	}

	viewport.AddSection(layoutTestViewportText, &layout.TextRenderer{
		Body:  "welcome to recordbaux",
		Style: viewportStyle,
	})

	b.AddSection(layoutTopBar, topBar)
	b.AddSection(layoutStatusBar, statusBar)
	b.AddSection(layoutViewport, viewport)
	b.SetJoinFunc(joinFunc)

	return b, nil
}

func joinFunc(sec map[layout.Section]layout.Renderer) string {
	var sections []string

	if topBar, ok := sec[layoutTopBar]; ok {
		renderedTopBar := topBar.Render()
		sections = append(sections, renderedTopBar)
	}

	if viewPort, ok := sec[layoutViewport]; ok {
		renderedViewport := viewPort.Render()
		sections = append(sections, renderedViewport)
	}

	if statusBar, ok := sec[layoutStatusBar]; ok {
		renderedStatusBar := statusBar.Render()
		sections = append(sections, renderedStatusBar)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
