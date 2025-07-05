package layout

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

func Splitscreen(width, height int, lContent, rContent string) (*Div, error) {
	dividerWidth := 1

	availableWidth := width - dividerWidth

	if availableWidth%2 != 0 {
		dividerWidth += 1
		availableWidth -= 1
	}

	panelWidth := availableWidth / 2

	leftPanel, err := New(Column, lipgloss.NewStyle(),
		WithName("splitscreen.leftPanel"),
		WithFixedWidth(panelWidth),
		WithFixedHeight(height),
		WithChild(&TextNode{Body: lContent}),
	)
	if err != nil {
		return nil, err
	}

	rightPanel, err := New(Column, style.Centered,
		WithName("splitscreen.rightPanel"),
		WithFixedWidth(panelWidth),
		WithFixedHeight(height),
		WithChild(&TextNode{Body: rContent}),
	)
	if err != nil {
		return nil, err
	}

	dividerStyle := lipgloss.NewStyle().
		Background(style.LightBlack).
		Width(dividerWidth).
		Height(height)

	dividerPanel, err := New(Column, dividerStyle,
		WithName("splitscreen.divider"),
		WithFixedWidth(dividerWidth),
		WithFixedHeight(height),
	)
	if err != nil {
		return nil, err
	}

	splitscreenContainer, err := New(Row, lipgloss.NewStyle(),
		WithName("splitscreen"),
		WithFixedWidth(width),
		WithFixedHeight(height),
		WithChild(leftPanel),
		WithChild(dividerPanel),
		WithChild(rightPanel),
	)
	if err != nil {
		return nil, err
	}

	return splitscreenContainer, nil
}
