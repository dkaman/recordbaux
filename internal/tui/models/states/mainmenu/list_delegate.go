package mainmenu

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

// shelfDelegate implements Bubble Tea's list.ItemDelegate for shelves
type shelfDelegate struct{ focused bool }

func (d shelfDelegate) Height() int  { return 1 }

func (d shelfDelegate) Spacing() int { return 0 }

func (d shelfDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	shM, ok := listItem.(shelf.Model)
	if !ok {
		return
	}

	shPhy := shM.PhysicalShelf()

	display := fmt.Sprintf("%s (%d bins × size %d)", shPhy.Name, len(shPhy.Bins), shPhy.Shape.BinSize())

	sty := style.TextStyle

	if m.Index() == index && d.focused {
		sty = style.ActiveTextStyle
	}

	prefix := "  "
	if m.Index() == index {
		prefix = "> "
	}

	w.Write([]byte(sty.Render(prefix + display)))
}

func (d shelfDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
