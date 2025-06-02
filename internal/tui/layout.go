package tui

import (
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
		AlignVertical(lipgloss.Center).
		Margin(0)
)


func newTUILayout() (*div.Div, error) {
	topBar, err := div.New(div.Row, barStyle,
		div.WithBorder(false),
		div.WithPadding(0, 0, 0, 0),
		div.WithFixedHeight(1),
	)
	if err != nil {
		return nil, err
	}

	helpBar, err := div.New(div.Row, helpBarStyle,
		div.WithBorder(false),
		div.WithPadding(0, 0, 0, 0),
		div.WithFixedHeight(1),
	)
	if err != nil {
		return nil, err
	}

	statusBar, err := div.New(div.Row, barStyle,
		div.WithBorder(false),
		div.WithPadding(0, 0, 0, 0),
		div.WithFixedHeight(1),
	)
	if err != nil {
		return nil, err
	}

	viewport, err := div.New(div.Row, viewportStyle,
		div.WithBorder(true),
		div.WithPadding(0, 0, 0, 0),
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
