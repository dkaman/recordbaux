package fetchfromdiscogs

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/tui/style"
)

func (s FetchFromDiscogsState) renderModel() string {
	var lines []string

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
	joined := lipgloss.JoinVertical(lipgloss.Center, lines...)
	content := style.ProgressStyle.Render(joined)

	// 3. Use lipgloss.Place to center the content on the screen.
	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
