package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

var (
	viewportStyle = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
)


func newTUILayout() (*layout.Div, error) {
	topBar, err := layout.New(layout.Row, style.BarStyle,
		layout.WithName("topbar"),
		layout.WithBorder(false),
		layout.WithFixedHeight(1),
	)
	if err != nil {
		return nil, err
	}

	helpBar, err := layout.New(layout.Row, style.HelpBarStyle,
		layout.WithName("helpbar"),
		layout.WithBorder(false),
		layout.WithFixedHeight(1),
		layout.WithHidden(true),
	)
	if err != nil {
		return nil, err
	}

	statusBar, err := layout.New(layout.Row, style.BarStyle,
		layout.WithName("statusbar"),
		layout.WithBorder(false),
		layout.WithFixedHeight(1),
	)
	if err != nil {
		return nil, err
	}

	viewport, err := layout.New(layout.Row, viewportStyle,
		layout.WithName("viewport"),
		layout.WithBorder(true),
	)
	if err != nil {
		return nil, err
	}

	root, err := layout.New(layout.Column, lipgloss.NewStyle())
	if err != nil {
		return nil, err
	}

	root.AddChild(topBar)
	root.AddChild(viewport)
	root.AddChild(helpBar)

	root.AddChild(statusBar)

	return root, nil
}

func addTopBarText(d *layout.Div, body string) error {
	if tb := d.Find("topbar"); tb != nil {
		tb.ClearChildren()
		tb.AddChild(&layout.TextNode{
			Body: body,
		})
		return nil
	}
	return fmt.Errorf("top bar not found in layout")
}

func addStatusBarText(d *layout.Div, body string) error {
	if sb := d.Find("statusbar"); sb != nil {
		sb.ClearChildren()
		sb.AddChild(&layout.TextNode{
			Body: body,
		})
		return nil
	}
	return fmt.Errorf("status bar not found in layout")
}

func addHelpBarText(d *layout.Div, body string) error {
	if hb := d.Find("helpbar"); hb != nil {
		hb.ClearChildren()
		hb.AddChild(&layout.TextNode{
			Body: body,
		})
		return nil
	}
	return fmt.Errorf("status bar not found in layout")
}
