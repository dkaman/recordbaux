package tui

import (

	lipgloss "github.com/charmbracelet/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/tui/style"
)

func (m Model) renderModel() string {
	canvas := lipgloss.NewCanvas()

	numBars := 2

	if m.helpVisible {
		numBars = 3
	}

	barStyle := style.BarStyle.
		Width(m.width).
		Height(1)

	helpStyle := style.HelpBarStyle.
		Width(m.width).
		Height(1)

	viewportStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width).
		Height(m.height-numBars)

	topBar := lipgloss.NewLayer(barStyle.Render(m.topBarText))
	helpBar := lipgloss.NewLayer(helpStyle.Render(m.Help()))
	statusBar := lipgloss.NewLayer(barStyle.Render(m.statusBarText))
	viewPort := lipgloss.NewLayer(viewportStyle.Render(m.stateMachine.View()))

	canvas.AddLayers(topBar.
		X(0).Y(1),
	)

	canvas.AddLayers(viewPort.
		X(0).Y(2),
	)

	if m.helpVisible {
		canvas.AddLayers(helpBar.
			X(0).Y(m.height-1),
		)
	}

	canvas.AddLayers(statusBar.
		X(0).Y(m.height),
	)

	return canvas.Render()
}
