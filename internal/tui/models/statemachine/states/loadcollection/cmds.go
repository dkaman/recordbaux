package loadcollection

import (
	"context"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"

	discogs "github.com/dkaman/discogs-golang"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/db/shelf"
)

// LoadCollectionMsg carries a physical.Shelf pointer to the state.
type LoadCollectionMsg struct {
	Shelf *shelf.Entity
}

type NewDiscogsCollectionMsg struct {
	Releases []*record.Entity
}

// WithCollection constructs a Tea command that sends a LoadCollectionMsg.
func RetrieveDiscogsCollection(c *discogs.Client, username string, folder string, log *slog.Logger) tea.Cmd {
	l := log.WithGroup("retreivediscogscmd")
	return func() tea.Msg {
		var rs []*record.Entity

		releases, err := c.Collection.GetReleasesByFolder(context.TODO(), username, 0)
		if err != nil {
			l.Error("error getting releases from discogs",
				slog.String("error", err.Error()),
			)
		}

		l.Info("releases from discogs", releases)

		for _, rel := range releases {
			r, err := record.New(rel)
			if err != nil {
				l.Error("error in release to record entity conversion",
					slog.Any("error", err),
				)
				continue

			}
			rs = append(rs, r)
		}

		return NewDiscogsCollectionMsg{
			Releases: rs,
		}
	}
}
