package loadcollection

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

const (
	layoutViewport layout.Section = iota
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

func newLoadedCollectionFormLayout(base *layout.Node, f *form) *layout.Node {
	r := &layout.TeaModelRenderer{
		Model: f,
		Style: formStyle,
	}

	base.AddSection(layoutViewport, r)
	base.SetJoinFunc(joinFunc)
	return base
}

func newLoadedCollectionProgressLayout(base *layout.Node, p progress.Model, s spinner.Model, fetching, doneFetching bool, pct float64) (*layout.Node, error) {
	var lines []string

	if fetching && !doneFetching {
		spinnerLine := lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Render(fmt.Sprintf("%s fetching collection from discogs...", s.View()))
		lines = append(lines, spinnerLine)
	} else {
		check := lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true).
			Render("✔ fetching collection complete")
		lines = append(lines, check)
	}

	// Always render progress bar below
	barLine := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Render(p.ViewAs(pct))

	lines = append(lines, barLine)

	r := &layout.TextRenderer{
		Body: lipgloss.JoinVertical(lipgloss.Center, lines...),
		Style: viewportStyle,
	}

	base.AddSection(layoutViewport, r)
	base.SetJoinFunc(joinFunc)

	return base, nil
}

func joinFunc(sec map[layout.Section]layout.Renderer) string {
	var sections []string

	if viewPort, ok := sec[layoutViewport]; ok {
		renderedViewport := viewPort.Render()
		sections = append(sections, renderedViewport)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
