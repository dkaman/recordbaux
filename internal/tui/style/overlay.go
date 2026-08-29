package style

import (
	lipgloss "charm.land/lipgloss/v2"
)

var (
	ModalContainerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(LightBlue).
		Background(DarkBlack).
		Padding(1, 2)

	ModalTitleStyle = lipgloss.NewStyle().
		Foreground(LightYellow).
		Bold(true).
		MarginBottom(1)
)

func RenderOverlay(baseView, modalView string, w, h int) string {
	if modalView == "" {
		return baseView
	}

	baseLayer := lipgloss.NewLayer(baseView).X(0).Y(0)

	styledModal := ModalContainerStyle.Render(modalView)
	modalWidth := lipgloss.Width(styledModal)
	modalHeight := lipgloss.Height(styledModal)

	x :=  max(0, (w - modalWidth) / 2)
	y := max(0, (h - modalHeight) / 2)

	modalLayer := lipgloss.NewLayer(styledModal).X(x).Y(y)

	comp := lipgloss.NewCompositor(baseLayer, modalLayer)
	return comp.Render()
}
