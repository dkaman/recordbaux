package mainmenu

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/progress"

	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/tui/models/flist"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"

	tplaylist "github.com/dkaman/recordbaux/internal/tui/models/playlist"
)

type shelfDelegateStyles struct {
	ItemStyle            lipgloss.Style
	ItemStyleBlurred     lipgloss.Style
	SelectedStyle        lipgloss.Style
	SelectedStyleBlurred lipgloss.Style
}

// shelfDelegate implements Bubble Tea's list.ItemDelegate for shelves
type shelfDelegate struct {
	focused bool
	prog    progress.Model

	Styles shelfDelegateStyles
}

func newShelfDelegate(s shelfDelegateStyles) shelfDelegate {
	prg := progress.New(progress.WithDefaultGradient())

	return shelfDelegate{
		focused: true,
		prog:    prg,
		Styles:  s,
	}
}

func (d shelfDelegate) Height() int {
	return 2
}

func (d shelfDelegate) Spacing() int {
	return 1
}

func (d shelfDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	shM, ok := listItem.(shelf.Model)
	if !ok {
		return
	}

	shPhy := shM.PhysicalShelf()

	pct := float64(shPhy.TotalRecords()) / float64((shPhy.BinSize * len(shPhy.Bins)))

	display := fmt.Sprintf("%s (%d bins × size %d)\n%s", shPhy.Name, len(shPhy.Bins), shPhy.BinSize, d.prog.ViewAs(pct))

	var sty lipgloss.Style
	if d.focused {
		sty = d.Styles.ItemStyle
	} else {
		sty = d.Styles.ItemStyleBlurred
	}

	if m.Index() == index {
		if d.focused {
			sty = d.Styles.SelectedStyle
		} else {
			sty = d.Styles.SelectedStyleBlurred
		}
	}

	prefix := "  "
	if m.Index() == index {
		prefix = "> "
	}

	w.Write([]byte(sty.Render(prefix + display)))
}
func (d shelfDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d shelfDelegate) Focus() flist.FocusableDelegate {
	d.focused = true
	return d
}

func (d shelfDelegate) Blur() flist.FocusableDelegate {
	d.focused = false
	return d
}

type playlistDelegateStyles struct {
	ItemStyle            lipgloss.Style
	ItemStyleBlurred     lipgloss.Style
	SelectedStyle        lipgloss.Style
	SelectedStyleBlurred lipgloss.Style
}

// playlistDelegate implements Bubble Tea's list.ItemDelegate for playlists
type playlistDelegate struct {
	focused bool
	Styles  playlistDelegateStyles
}

func newPlaylistDelegate(s playlistDelegateStyles) playlistDelegate {
	return playlistDelegate{
		focused: true,
		Styles:  s,
	}
}

func (d playlistDelegate) Height() int { return 1 }

func (d playlistDelegate) Spacing() int { return 0 }

func (d playlistDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	plM, ok := listItem.(tplaylist.Model)
	if !ok {
		return
	}

	display := fmt.Sprintf("%s (%s)", plM.Title(), plM.Description())

	var sty lipgloss.Style
	if d.focused {
		sty = d.Styles.ItemStyle
	} else {
		sty = d.Styles.ItemStyleBlurred
	}

	if m.Index() == index {
		if d.focused {
			sty = d.Styles.SelectedStyle
		} else {
			sty = d.Styles.SelectedStyleBlurred
		}
	}

	prefix := "  "
	if m.Index() == index {
		prefix = "> "
	}

	w.Write([]byte(sty.Render(prefix + display)))
}

func (d playlistDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d playlistDelegate) Focus() flist.FocusableDelegate {
	d.focused = true
	return d
}

func (d playlistDelegate) Blur() flist.FocusableDelegate {
	d.focused = false
	return d
}
