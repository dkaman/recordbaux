package loadedbin

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/tui/style"
)

func (s LoadedBinState) renderModel() string {
	canvas := lipgloss.NewCanvas()

	boxW := s.width/2
	boxH := s.height

	s.records.SetWidth(boxW)
	s.records.SetHeight(boxH)
	s.selectedRecord.SetSize(boxW, boxH)

	left := lipgloss.NewLayer(style.BaseTableStyle.Render(s.records.View()))

	right := lipgloss.NewLayer(s.selectedRecord.View())

	canvas.AddLayers(left.
		Width(boxW).
		Height(boxH).
		X(0).Y(0),
	)

	canvas.AddLayers(right.
		Width(boxW).
		Height(boxH).
		X(boxH).Y(0),
	)

	return canvas.Render()
}
