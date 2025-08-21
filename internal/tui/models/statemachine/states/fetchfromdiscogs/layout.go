package fetchfromdiscogs

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

func newFetchFromDiscogslayout(base *layout.Div, pb progress.Model, sp spinner.Model, pct float64, fetching bool, title string) (*layout.Div, error) {
	base.ClearChildren()

	var lines []string

	// 1) spinner vs. checkmark
	if fetching {
		// spinner spinning, “fetching…”
		spinnerLine := lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Render(fmt.Sprintf("%s fetching collection from Discogs…", sp.View()))
		lines = append(lines, spinnerLine)
	} else {
		// checkmark when fetching→inserting transition is done
		check := lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true).
			Render("✔ fetching collection complete")

		enriching := fmt.Sprintf("enriching %s from discogs", title)
		lines = append(lines, check, enriching)
	}

	// 2) always show the progress bar underneath
	barLine := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Render(pb.ViewAs(pct))
	lines = append(lines, barLine)

	// join the two lines vertically, then wrap in a single TextNode
	joined := lipgloss.JoinVertical(lipgloss.Center, lines...)
	base.AddChild(&layout.TextNode{
		Body: style.ProgressStyle.Render(joined),
	})

	return base, nil
}
