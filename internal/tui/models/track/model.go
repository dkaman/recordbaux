package track

import (
	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db/track"
)

// Model is a wrapper for a track entity for use in the TUI.
type Model struct {
	physicalTrack *track.Entity
}

// New creates a new track model.
func New(t *track.Entity) Model {
	return Model{
		physicalTrack: t,
	}
}

// Init is a no-op.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is a no-op.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

// View is a no-op.
func (m Model) View() string {
	return ""
}

// FilterValue implements the list.Item interface for filtering.
func (m Model) FilterValue() string {
	if m.physicalTrack == nil {
		return ""
	}
	return m.physicalTrack.Title
}

// Title returns the track's name.
func (m Model) Title() string {
	return m.FilterValue()
}

// Description returns a summary of the track's details.
func (m Model) Description() string {
	if m.physicalTrack == nil {
		return ""
	}
	return m.physicalTrack.Title
}

// PhysicalTrack returns the underlying track entity.
func (m Model) PhysicalTrack() *track.Entity {
	return m.physicalTrack
}
