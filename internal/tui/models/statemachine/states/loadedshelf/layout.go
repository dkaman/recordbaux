package loadedshelf

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

func (s LoadedShelfState) renderModel() string {
	canvas := lipgloss.NewCanvas()

	// Base layer: the shelf view
	shelfView := s.shelf.View()
	baseLayer := lipgloss.NewLayer(shelfView)
	canvas.AddLayers(baseLayer.X(0).Y(0))

	// Modal layer: load collection form
	if s.loading {
		formView := s.loadCollectionForm.View()
		formX := (s.width - lipgloss.Width(formView)) / 2
		formY := (s.height - lipgloss.Height(formView)) / 2
		formLayer := lipgloss.NewLayer(formView)
		canvas.AddLayers(formLayer.X(formX).Y(formY).Z(1))
	}

	// Modal layer: fetching progress
	if s.fetching {
		var content string
		if s.totalReleases == 0 {
			content = fmt.Sprintf("%s fetching collection from Discogs...", s.spin.View())
		} else {
			header := "✔ collection loaded, enriching records...\n\n"
			title := fmt.Sprintf("loading: %s\n\n", s.currentTitle)
			progress := s.prog.ViewAs(s.pct)
			percent := fmt.Sprintf(" %d/%d", s.currentIndex, s.totalReleases)
			content = lipgloss.JoinVertical(lipgloss.Left, header, title, progress+percent)
		}

		progX := (s.width - lipgloss.Width(content)) / 2
		progY := (s.height - lipgloss.Height(content)) / 2
		progLayer := lipgloss.NewLayer(content)
		canvas.AddLayers(progLayer.X(progX).Y(progY).Z(1))
	}

	return canvas.Render()
}
