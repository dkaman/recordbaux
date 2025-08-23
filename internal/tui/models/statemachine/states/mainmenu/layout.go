package mainmenu

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
)

func newMainMenuLayout(base *layout.Div, shelves, playlists list.Model, focus focusedView) (*layout.Div, error) {
	base.ClearChildren()

	if len(shelves.Items()) == 0 && len(playlists.Items()) == 0 {
		base.AddChild(&layout.TextNode{
			Body: "no shelves or playlists defined, 'o' to create shelf...",
		})
		return base, nil
	}

	// Create styles for focused and blurred list containers
	focusedStyle := lipgloss.NewStyle().BorderForeground(style.LightGreen)
	blurredStyle := lipgloss.NewStyle().BorderForeground(style.DarkWhite)

	var shelfBoxStyle, playlistBoxStyle lipgloss.Style
	if focus == shelvesView {
		shelfBoxStyle = focusedStyle
		playlistBoxStyle = blurredStyle
	} else {
		shelfBoxStyle = blurredStyle
		playlistBoxStyle = focusedStyle
	}

	// Create Divs for each list. The content will be added after we know the dimensions.
	shelfBox, _ := layout.New(layout.Column, shelfBoxStyle,
		layout.WithBorder(true),
	)

	playlistBox, _ := layout.New(layout.Column, playlistBoxStyle,
		layout.WithBorder(true),
	)

	// Use a Row Div to contain them side-by-side.
	listsContainer, _ := layout.New(layout.Row, lipgloss.NewStyle(),
		layout.WithChild(shelfBox),
		layout.WithChild(playlistBox),
	)

	base.AddChild(listsContainer)

	base.Resize(base.Width(), base.Height())

	// Now that the containers know their size, calculate the inner area for the lists.
	// We subtract 2 to account for the container's top/bottom and left/right borders.
	shelfListW := shelfBox.Width()
	shelfListH := shelfBox.Height()
	shelves.SetSize(shelfListW, shelfListH-2)

	playlistListW := playlistBox.Width()
	playlistListH := playlistBox.Height()
	playlists.SetSize(playlistListW, playlistListH-2)

	// With the lists correctly sized, render their views and place them in their containers.
	shelfBox.ClearChildren()
	shelfBox.AddChild(&layout.TextNode{Body: shelves.View()})

	playlistBox.ClearChildren()
	playlistBox.AddChild(&layout.TextNode{Body: playlists.View()})

	return base, nil
}
