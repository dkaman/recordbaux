package playlist

import (
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db/playlist"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

type Model struct {
	width, height int
	entity        *playlist.Entity
	trackTable    table.Model
}

func New() Model {
	columns := []table.Column{
		{Title: "Position", Width: 8},
		{Title: "Title", Width: 50},
		{Title: "Duration", Width: 8},
		{Title: "Key", Width: 3},
		{Title: "BPM", Width: 3},
	}

	tbl := table.New(table.WithColumns(columns), table.WithFocused(true))
	tbl.SetStyles(style.DefaultTableStyles())

	return Model{
		trackTable: tbl,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = sizeMsg.Width, sizeMsg.Height
		return m, nil
	}

	var tableCmd tea.Cmd
	m.trackTable, tableCmd = m.trackTable.Update(msg)
	cmds = append(cmds, tableCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	if m.entity == nil {
		return tea.NewView("no playlist is loaded yet...")
	}

	m.trackTable.SetWidth(m.width)
	m.trackTable.SetHeight(m.height)
	tableContent := m.trackTable.View()

	return tea.NewView(tableContent)
}

// FilterValue implements the list.Item interface for filtering.
func (m Model) FilterValue() string {
	if m.entity == nil {
		return ""
	}
	return m.entity.Name
}

func (m Model) Title() string {
	return m.FilterValue()
}

func (m Model) Description() string {
	if m.entity == nil {
		return ""
	}
	return fmt.Sprintf("%d tracks", len(m.entity.Tracks))
}

func (m Model) PhysicalPlaylist() *playlist.Entity {
	return m.entity
}

// this is ugly i want to get rid of this eventually
func (m *Model) SetEntity(p *playlist.Entity) {
	m.entity = p

	var rows []table.Row

	if p != nil {
		for _, t := range p.Tracks {
			rows = append(rows, table.Row{
				t.Position,
				t.Title,
				t.Duration,
				t.Key,
				fmt.Sprintf("%d", t.BPM),
			})
		}
	}

	m.trackTable.SetRows(rows)
}
