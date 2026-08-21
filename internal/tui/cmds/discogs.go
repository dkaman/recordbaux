package cmds

import (
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/discogs-golang"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/db/shelf"
)

// LoadCollectionMsg carries a physical.Shelf pointer to the state.
type LoadCollectionMsg struct {
	Shelf *shelf.Entity
}

type NewDiscogsCollectionIntentMsg struct {
	Folder string
}

type NewDiscogsCollectionMsg struct {
	Releases []*record.Entity
	Err      error
}

type NewDiscogsEnrichRecordIntentMsg struct {
	Record *record.Entity
}

type NewDiscogsEnrichRecordMsg struct {
	Record *record.Entity
	Err    error
}

type GetDiscogsFoldersIntentMsg struct {}

type GetDiscogsFoldersMsg struct {
	Folders []discogs.Folder
}

// WithCollection constructs a Tea command that sends a LoadCollectionMsg.
func RetrieveDiscogsCollectionCmd(folder string) tea.Cmd {
	return func() tea.Msg {
		return NewDiscogsCollectionIntentMsg{
			Folder: folder,
		}
	}
}

func EnrichReleaseInstanceCmd(rec *record.Entity) tea.Cmd {
	return func() tea.Msg {
		return NewDiscogsEnrichRecordIntentMsg{
			Record: rec,
		}
	}
}

func ListDiscogsFoldersCmd() tea.Cmd {
	return func() tea.Msg {
		return GetDiscogsFoldersIntentMsg{}
	}
}
