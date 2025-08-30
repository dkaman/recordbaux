package cmds

import (
	"context"
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"

	discogs "github.com/dkaman/discogs-golang"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/db/track"
)

// LoadCollectionMsg carries a physical.Shelf pointer to the state.
type LoadCollectionMsg struct {
	Shelf *shelf.Entity
}

type NewDiscogsCollectionMsg struct {
	Releases []*record.Entity
	Err      error
}

type NewDiscogsEnrichRecordMsg struct {
	Record *record.Entity
	Err    error
}

// WithCollection constructs a Tea command that sends a LoadCollectionMsg.
func RetrieveDiscogsCollection(c *discogs.Client, username string, folder string, log *slog.Logger) tea.Cmd {
	l := log.WithGroup("retreivediscogscmd")
	return func() tea.Msg {
		var recs []*record.Entity
		var msg NewDiscogsCollectionMsg

		releaseInstances, err := c.Collection.GetReleasesByFolder(context.TODO(), username, 0)
		if err != nil {
			msg.Err = fmt.Errorf("error getting releases from discogs: %w" , err)
			return msg
		}

		l.Debug("got releases from discogs",
			slog.Int("count", len(releaseInstances)),
		)

		for _, ri := range releaseInstances {
			l.Debug("processing release",
				slog.Int("id", ri.ID),
			)

			r, err := record.New(ri)
			if err != nil {
				msg.Err = fmt.Errorf("error constructing record entity: %w", err)
				return msg
			}

			recs = append(recs, r)
		}

		msg.Releases = recs
		msg.Err = nil

		return msg
	}
}

func EnrichReleaseInstance(c *discogs.Client, rec *record.Entity) tea.Cmd {
	return func() tea.Msg {
		var msg NewDiscogsEnrichRecordMsg

		rel, err := c.Database.GetRelease(context.TODO(), rec.ReleaseID)
		if err != nil {
			msg.Err = err
			return msg
		}

		var tracks []*track.Entity
		for _, trk := range rel.Tracklist {
			t, err := track.New(trk)
			if err != nil {
				msg.Err = err
				return msg
			}

			tracks = append(tracks, t)
		}

		rec.Tracklist = tracks

		return NewDiscogsEnrichRecordMsg{
			Record: rec,
			Err:    nil,
		}
	}
}
