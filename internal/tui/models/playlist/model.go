package playlist

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db/playlist"
)

// Model is a wrapper for a playlist entity for use in the TUI.
type Model struct {
	physicalPlaylist *playlist.Entity
}

// New creates a new playlist model.
func New(p *playlist.Entity) Model {
	return Model{
		physicalPlaylist: p,
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
	if m.physicalPlaylist == nil {
		return ""
	}
	return m.physicalPlaylist.Name
}

// Title returns the playlist's name.
func (m Model) Title() string {
	return m.FilterValue()
}

// Description returns a summary of the playlist's contents.
func (m Model) Description() string {
	if m.physicalPlaylist == nil {
		return ""
	}
	return fmt.Sprintf("%d tracks", len(m.physicalPlaylist.Tracks))
}

// PhysicalPlaylist returns the underlying playlist entity.
func (m Model) PhysicalPlaylist() *playlist.Entity {
	return m.physicalPlaylist
}
