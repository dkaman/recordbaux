package mainmenu

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/v2/list"

	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/style"

	tplaylist "github.com/dkaman/recordbaux/internal/tui/models/playlist"
)

// shelfDelegate implements Bubble Tea's list.ItemDelegate for shelves
type shelfDelegate struct {
	dim bool
}

func (d shelfDelegate) Height() int  { return 1 }
func (d shelfDelegate) Spacing() int { return 0 }

func (d shelfDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	shM, ok := listItem.(shelf.Model)
	if !ok {
		return
	}

	shPhy := shM.PhysicalShelf()

	display := fmt.Sprintf("%s (%d bins × size %d)", shPhy.Name, len(shPhy.Bins), shPhy.BinSize)

	var sty lipgloss.Style
	if d.dim {
		sty = style.TextStyleDimmed
	} else {
		sty = style.TextStyle
	}

	if m.Index() == index {
		if d.dim {
			sty = style.ActiveTextStyleDimmed
		} else {
			sty = style.ActiveTextStyle
		}
	}

	prefix := "  "
	if m.Index() == index {
		prefix = "> "
	}

	w.Write([]byte(sty.Render(prefix + display)))
}

func (d shelfDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// playlistDelegate implements Bubble Tea's list.ItemDelegate for playlists
type playlistDelegate struct{
	dim bool
}

func (d playlistDelegate) Height() int  { return 1 }
func (d playlistDelegate) Spacing() int { return 0 }

func (d playlistDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	plM, ok := listItem.(tplaylist.Model)
	if !ok {
		return
	}

	display := fmt.Sprintf("%s (%s)", plM.Title(), plM.Description())

	var sty lipgloss.Style
	if d.dim {
		sty = style.TextStyleDimmed
	} else {
		sty = style.TextStyle
	}

	if m.Index() == index {
		if d.dim {
			sty = style.ActiveTextStyleDimmed
		} else {
			sty = style.ActiveTextStyle
		}
	}

	prefix := "  "
	if m.Index() == index {
		prefix = "> "
	}

	w.Write([]byte(sty.Render(prefix + display)))
}

func (d playlistDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
