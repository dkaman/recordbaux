package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

var (
	barStyle = lipgloss.NewStyle().
		Background(style.DarkBlue).
		Foreground(style.DarkBlack).
		Bold(true)

	helpBarStyle = lipgloss.NewStyle().
		Background(style.LightGreen).
		Foreground(style.LightBlack)

	viewportStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
)


func newTUILayout() (*div.Div, error) {
	topBar, err := div.New(div.Row, barStyle,
		div.WithName("topbar"),
		div.WithBorder(false),
		div.WithFixedHeight(1),
	)
	if err != nil {
		return nil, err
	}

	helpBar, err := div.New(div.Row, helpBarStyle,
		div.WithName("helpbar"),
		div.WithBorder(false),
		div.WithFixedHeight(1),
		div.WithHidden(true),
	)
	if err != nil {
		return nil, err
	}

	statusBar, err := div.New(div.Row, barStyle,
		div.WithName("statusbar"),
		div.WithBorder(false),
		div.WithFixedHeight(1),
	)
	if err != nil {
		return nil, err
	}

	viewport, err := div.New(div.Row, viewportStyle,
		div.WithName("viewport"),
		div.WithBorder(true),
	)
	if err != nil {
		return nil, err
	}

	root, err := div.New(div.Column, lipgloss.NewStyle())
	if err != nil {
		return nil, err
	}

	root.AddChild(topBar)
	root.AddChild(viewport)
	root.AddChild(helpBar)

	root.AddChild(statusBar)

	return root, nil
}

func addTopBarText(d *div.Div, body string) error {
	if tb := d.Find("topbar"); tb != nil {
		tb.ClearChildren()

		tb.AddChild(&div.TextNode{
			Body: body,
		})
		return nil
	}
	return fmt.Errorf("top bar not found in layout")
}

func addStatusBarText(d *div.Div, body string) error {
	if sb := d.Find("statusbar"); sb != nil {
		sb.ClearChildren()

		sb.AddChild(&div.TextNode{
			Body: body,
		})
		return nil
	}
	return fmt.Errorf("status bar not found in layout")
}

func addHelpBarText(d *div.Div, body string) error {
	if hb := d.Find("helpbar"); hb != nil {
		hb.ClearChildren()
		hb.AddChild(&div.TextNode{
			Body: body,
		})
		return nil
	}
	return fmt.Errorf("status bar not found in layout")
}
