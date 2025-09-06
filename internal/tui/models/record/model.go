package record

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/v2/table"

	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

// Model represents the record view, displaying detailed information about a single record.
type Model struct {
	physicalRecord *record.Entity
	width, height  int
	tracklistTable table.Model // The table for the tracklist
	trackIndex     int         // The currently selected track index
}

// New creates a new Model for a specific record.
func New(r *record.Entity) Model {
	columns := []table.Column{
		{Title: "position", Width: 10},
		{Title: "title", Width: 40},
		{Title: "duration", Width: 10},
	}

	// Create the rows from the record's tracklist.
	rows := make([]table.Row, len(r.Tracklist))
	for i, t := range r.Tracklist {
		rows[i] = table.Row{t.Position, t.Title, t.Duration}
	}

	// Create a new table model and populate it with the rows.
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows), // Set the rows here!
	)

	// Set the table styles
	t.SetStyles(style.DefaultTableStyles())

	return Model{
		physicalRecord: r,
		tracklistTable: t,
		trackIndex:     0, // Initialize the selected track to the first one
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages for the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			// Move the track selection up
			if m.tracklistTable.Cursor() > 0 {
				m.tracklistTable.MoveUp(1)
			}
		case "down":
			// Move the track selection down
			if m.tracklistTable.Cursor() < len(m.physicalRecord.Tracklist)-1 {
				m.tracklistTable.MoveDown(1)
			}
		}
		// Sync the track index with the table cursor
		m.trackIndex = m.tracklistTable.Cursor()
	}

	// Pass the message to the table model
	m.tracklistTable, cmd = m.tracklistTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the model's UI.
func (m Model) View() string {
	if m.physicalRecord == nil {
		return "No record loaded"
	}

	canvas := lipgloss.NewCanvas()

	// Calculate height for each section to divide the screen
	sectionHeight := m.height / 3

	// --- Render Record Info Card ---
	recordInfoCard := m.renderRecordInfoCard(m.width, sectionHeight)

	// --- Render Tracklist Table ---
	m.tracklistTable.SetWidth(m.width)
	m.tracklistTable.SetHeight(sectionHeight)
	tracklistTable := m.tracklistTable.View()

	// // --- Render Track Info Card ---
	// trackInfoCard := m.renderTrackInfoCard(m.width, sectionHeight)

	recordCard := lipgloss.NewLayer(recordInfoCard)
	tracklist := lipgloss.NewLayer(tracklistTable)
	// trackInfo := lipgloss.NewLayer(trackInfoCard)

	canvas.AddLayers(
		recordCard.
			X(0).Y(0),
		tracklist.
			X(0).Y(sectionHeight),
		// trackInfo.
		// 	X(0).Y(2*sectionHeight),
	)

	return canvas.Render()
}

// SetSize updates the model's dimensions.
func (m Model) SetSize(w, h int) Model {
	m.width = w
	m.height = h
	return m
}

// renderRecordInfoCard formats the record data into a string for the UI.
func (m Model) renderRecordInfoCard(w, h int) string {
	var s strings.Builder
	s.WriteString("Record Card\n\n") // Add a title
	s.WriteString(fmt.Sprintf("Title: %s\n", m.physicalRecord.Title))
	s.WriteString(fmt.Sprintf("Artists: %s\n", strings.Join(m.physicalRecord.Artists, ", ")))
	s.WriteString(fmt.Sprintf("Catalog Number: %s\n", m.physicalRecord.CatalogNumber))
	s.WriteString(fmt.Sprintf("physical location: %s\n", m.physicalRecord.Coordinate))

	sty := lipgloss.NewStyle().
		Width(w).
		Height(h).
		Align(lipgloss.Center)

	return sty.Render(s.String())
}

// renderTrackInfoCard formats a single track's data for the UI.
func (m Model) renderTrackInfoCard(w, h int) string {
	if len(m.physicalRecord.Tracklist) == 0 {
		return "No track info available."
	}
	t := m.physicalRecord.Tracklist[m.tracklistTable.Cursor()]

	card := fmt.Sprintf("Track Info Card\n\nTitle: %s\nArtist: %s\nDuration: %s", t.Title, m.physicalRecord.Artists[0], t.Duration)

	sty := lipgloss.NewStyle().
		Width(w).
		Height(h).
		Align(lipgloss.Center)

	return sty.Render(card)
}
