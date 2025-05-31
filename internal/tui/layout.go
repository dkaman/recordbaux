package tui

import (
	"fmt"

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

	layoutTopBarBody
	layoutViewportBody
	layoutStatusBarBody
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

	topBar.AddSection(layoutTopBarBody, &layout.TextRenderer{
		Body:  "recordbaux - organize your record collection",
		Style: barStyle,
	})

	statusBar, err := layout.New(b.Width, 1)
	if err != nil {
		return nil, err
	}

	statusBar.AddSection(layoutStatusBarBody, &layout.TextRenderer{
		Body:  "status bar test",
		Style: barStyle,
	})

	usableHeight := b.Height - 2

	frameWidth, frameHeight := viewportStyle.GetFrameSize()

	vpWidth := b.Width - frameWidth
	vpHeight := usableHeight - frameHeight

	viewport, err := layout.New(vpWidth, vpHeight)
	if err != nil {
		return nil, err
	}

	viewport.AddSection(layoutViewportBody, &layout.TextRenderer{
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

func setStatusBarText(n *layout.Node, body string) error {
	statusBarRenderer, err := n.GetSection(layoutStatusBar)
	if err != nil {
		return err
	}

	if n, ok := statusBarRenderer.(*layout.Node); ok {
		n.AddSection(layoutStatusBarBody, &layout.TextRenderer{
			Body: body,
			Style: barStyle,
		})
		return nil
	} else {
		return fmt.Errorf("status bar section is not a *layout.Node")
	}
}
