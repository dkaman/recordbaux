package loadcollection

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"

	discogs "github.com/dkaman/discogs-golang"
)

// LoadCollectionMsg carries a physical.Shelf pointer to the state.
type LoadCollectionMsg struct {
	Shelf *physical.Shelf
}

type NewDiscogsCollectionMsg struct {
	Releases []*physical.Record
}

// WithCollection constructs a Tea command that sends a LoadCollectionMsg.
func RetrieveDiscogsCollection(c *discogs.Client, username string, folder string) tea.Cmd {
	return func() tea.Msg {
		var rs []*physical.Record

		releases, _ := c.Collection.GetReleasesByFolder(context.TODO(), username, 0)
		for _, rel := range releases {
			rs = append(rs, &physical.Record{Release: rel})
		}

		return NewDiscogsCollectionMsg{
			Releases: rs,
		}
	}
}
