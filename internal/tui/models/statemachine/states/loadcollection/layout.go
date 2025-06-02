package loadcollection

import (
	"github.com/charmbracelet/lipgloss"

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
	return base
}

// func newLoadedCollectionProgressLayout(base *layout.Node, p progress.Model, s spinner.Model, fetching, doneFetching bool, pct float64) (*layout.Node, error) {
// 	var lines []string

// 	if fetching && !doneFetching {
// 		spinnerLine := lipgloss.NewStyle().
// 			AlignHorizontal(lipgloss.Center).
// 			Render(fmt.Sprintf("%s fetching collection from discogs...", s.View()))
// 		lines = append(lines, spinnerLine)
// 	} else {
// 		check := lipgloss.NewStyle().
// 			Foreground(lipgloss.Color("10")).
// 			Bold(true).
// 			Render("✔ fetching collection complete")
// 		lines = append(lines, check)
// 	}

// 	// Always render progress bar below
// 	barLine := lipgloss.NewStyle().
// 		AlignHorizontal(lipgloss.Center).
// 		Render(p.ViewAs(pct))

// 	lines = append(lines, barLine)

// 	r := &layout.TextRenderer{
// 		Body: lipgloss.JoinVertical(lipgloss.Center, lines...),
// 		Style: viewportStyle,
// 	}

// 	base.AddSection(layoutViewport, r)
// 	base.SetJoinFunc(joinFunc)

// 	return base, nil
// }
