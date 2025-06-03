package loadcollection

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"

	"github.com/dkaman/recordbaux/internal/tui/style/div"
)

var (
	viewportStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Margin(0)

	formStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)

	progressStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)

func newLoadedCollectionFormLayout(base *div.Div, f *form) *div.Div {
	base.ClearChildren()
	base.AddChild(&div.TextNode{
		Body: formStyle.Render(f.View()),
	})
	return base
}

func newLoadedCollectionProgressLayout(
	base *div.Div,
	pb progress.Model,
	sp spinner.Model,
	fetching bool,
	inserting bool,
	pct float64,
) (*div.Div, error) {
	base.ClearChildren()

	var lines []string

	// 1) spinner vs. checkmark
	if fetching && !inserting {
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
		lines = append(lines, check)
	}

	// 2) always show the progress bar underneath
	barLine := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Render(pb.ViewAs(pct))
	lines = append(lines, barLine)

	// join the two lines vertically, then wrap in a single TextNode
	joined := lipgloss.JoinVertical(lipgloss.Center, lines...)
	base.AddChild(&div.TextNode{
		Body: progressStyle.Render(joined),
	})

	return base, nil
}
