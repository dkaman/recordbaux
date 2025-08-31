package fetchfromdiscogs

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

func (s FetchFromDiscogsState) renderModel() string {
	var lines []string

	canvas := lipgloss.NewCanvas()

	// 1. Build the content lines, same as the old layout function.
	if s.fetching {
		spinnerLine := fmt.Sprintf("%s fetching collection from Discogs…", s.spinner.View())
		lines = append(lines, spinnerLine)
	} else {
		check := lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true).
			Render("✔ fetching collection complete")
		enriching := fmt.Sprintf("enriching %s from discogs", s.currentTitle)
		lines = append(lines, check, enriching)
	}

	// Always show the progress bar underneath.
	barLine := s.progress.ViewAs(s.pct)
	lines = append(lines, barLine)

	// 2. Join the lines and apply the final style.
	content := lipgloss.JoinVertical(lipgloss.Center, lines...)

	contentW := lipgloss.Width(content)
	contentH := lipgloss.Height(content)

	borderedContent := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(contentW).
		Height(contentH).
		Render(content)

	contentLayer := lipgloss.NewLayer(borderedContent)
	contentX := (s.width - contentW) / 2
	contentY := (s.height - contentH) / 2

	canvas.AddLayers(contentLayer.
		X(contentX).Y(contentY),
	)

	return canvas.Render()
}
