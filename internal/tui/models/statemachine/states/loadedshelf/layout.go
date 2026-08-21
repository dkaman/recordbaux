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
