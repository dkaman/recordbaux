package loadedshelf

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/tui/style"
)

func (s LoadedShelfState) renderModel() tea.View {
	// 1. Modal State: Load Collection Form
	if s.loading {
		formView := s.loadCollectionForm.View().Content
		modal := style.ModalStyle.Render(formView)

		// Place the rendered modal in the absolute center of the terminal window
		return tea.NewView(lipgloss.Place(s.width, s.height, lipgloss.Center, lipgloss.Center, modal))
	}

	// 2. Modal State: Fetching Progress
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

		modal := style.ModalStyle.Render(content)

		// Place the rendered progress box in the absolute center
		return tea.NewView(lipgloss.Place(s.width, s.height, lipgloss.Center, lipgloss.Center, modal))
	}

	// 3. Base State: Standard Shelf View
	// Anchor the shelf view to the top-left, guaranteeing it claims the full terminal dimensions
	shelfContent := s.shelf.View().Content
	return tea.NewView(lipgloss.Place(s.width, s.height, lipgloss.Left, lipgloss.Top, shelfContent))
}

// func (s LoadedShelfState) renderModel() string {
// 	var comp *lipgloss.Compositor

// 	shelfView := s.shelf.View()
// 	baseLayer := lipgloss.NewLayer(shelfView.Content)

// 	comp = lipgloss.NewCompositor(baseLayer)

// 	// Modal layer: load collection form
// 	if s.loading {
// 		formView := s.loadCollectionForm.View().Content

// 		formW := lipgloss.Width(formView)
// 		formH := lipgloss.Height(formView)

// 		modal := style.ModalStyle.
// 			Width(formW).
// 			Height(formH).
// 			Render(formView)

// 		formLayer := lipgloss.NewLayer(modal)

// 		formX := (s.width - lipgloss.Width(formView)) / 2
// 		formY := (s.height - lipgloss.Height(formView)) / 2

// 		comp.AddLayers(formLayer.X(formX).Y(formY).Z(1))
// 	}

// 	// Modal layer: fetching progress
// 	if s.fetching {
// 		var content string
// 		if s.totalReleases == 0 {
// 			content = fmt.Sprintf("%s fetching collection from Discogs...", s.spin.View())
// 		} else {
// 			header := "✔ collection loaded, enriching records...\n\n"
// 			title := fmt.Sprintf("loading: %s\n\n", s.currentTitle)
// 			progress := s.prog.ViewAs(s.pct)
// 			percent := fmt.Sprintf(" %d/%d", s.currentIndex, s.totalReleases)
// 			content = lipgloss.JoinVertical(lipgloss.Left, header, title, progress+percent)
// 		}

// 		formW := lipgloss.Width(content)
// 		formH := lipgloss.Height(content)

// 		modal := style.ModalStyle.
// 			Width(formW).
// 			Height(formH).
// 			Render(content)

// 		progLayer := lipgloss.NewLayer(modal)

// 		progX := (s.width - lipgloss.Width(content)) / 2
// 		progY := (s.height - lipgloss.Height(content)) / 2

// 		comp.AddLayers(progLayer.X(progX).Y(progY).Z(1))
// 	}

// 	return comp.Render()
// }
