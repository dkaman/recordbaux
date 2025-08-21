package createplaylist

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/db/playlist"
	"github.com/dkaman/recordbaux/internal/db/track"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type refreshMsg struct{}

func (s CreatePlaylistState) refresh() tea.Cmd {
	return func() tea.Msg { return refreshMsg{} }
}

type CreatePlaylistState struct {
	nextState states.StateType
	shelves   *services.ShelfService
	tracks    *services.TrackService
	layout    *layout.Div
	logger    *slog.Logger
	table     table.Model
	keys      keyMap

	playlists      *services.PlaylistService
	selectedTracks map[uint]*track.Entity
	namingPlaylist bool
	nameForm       *form
	playlistName   string
}

func New(s *services.ShelfService, t *services.TrackService, p *services.PlaylistService, l *layout.Div, log *slog.Logger) CreatePlaylistState {
	logger := log.WithGroup("createplayliststate")

	// Add a column for the selection indicator
	columns := []table.Column{
		{Title: "", Width: 3}, // For [x]
		{Title: "Position", Width: 10},
		{Title: "Title", Width: 50},
		{Title: "Duration", Width: 10},
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	tbl.SetStyles(style.DefaultTableStyles())
	return CreatePlaylistState{
		nextState:      states.Undefined,
		shelves:        s,
		tracks:         t,
		playlists:      p, // Store the playlist service
		layout:         l,
		logger:         logger,
		table:          tbl,
		keys:           defaultKeybinds(),
		selectedTracks: make(map[uint]*track.Entity), // Initialize the map
	}
}

func (s CreatePlaylistState) Init() tea.Cmd {
	return s.refresh()
}

func (s CreatePlaylistState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if s.namingPlaylist {
		// Update the form
		fModel, formUpdatesCmds := s.nameForm.Update(msg)
		if f, ok := fModel.(*form); ok {
			s.nameForm = f
		}
		cmds = append(cmds, formUpdatesCmds)

		// Re-render the form in the layout
		s.layout.ClearChildren()
		centeredForm := layout.CenteredBox(s.layout.Width(), s.layout.Height(), s.nameForm.View(), 0.5, 0.25)
		s.layout.AddChild(centeredForm)

		if s.nameForm.State == huh.StateCompleted {
			name := s.nameForm.Name()
			newPlaylist := &playlist.Entity{
				Name:   name,
				Tracks: make([]*track.Entity, 0, len(s.selectedTracks)),
			}

			for _, t := range s.selectedTracks {
				newPlaylist.Tracks = append(newPlaylist.Tracks, t)
			}

			// Dispatch save command and transition back to main menu
			cmds = append(cmds, tcmds.SavePlaylistCmd(s.playlists.Playlists, newPlaylist, s.logger))

			s.selectedTracks = make(map[uint]*track.Entity)
			s.namingPlaylist = false
			s.nameForm = newNameForm()
			s.nextState = states.MainMenu
		}
		return s, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.layout.Resize(msg.Width, msg.Height)

	case refreshMsg:
		s.logger.Debug("refreshing tracks from service")
		var rows []table.Row
		for _, t := range s.tracks.AllTracks {
			selectedMarker := "[ ]"
			if _, ok := s.selectedTracks[t.ID]; ok {
				selectedMarker = "[x]"
			}
			rows = append(rows, table.Row{selectedMarker, t.Position, t.Title, t.Duration})
		}
		s.table.SetRows(rows)
		s.layout, _ = newCreatePlaylistLayout(s.layout, s.table)
		return s, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			s.nextState = states.MainMenu
			return s, nil

		case key.Matches(msg, s.keys.Select):
			// Get the selected track
			if len(s.tracks.AllTracks) == 0 {
				return s, nil
			}
			selectedTrack := s.tracks.AllTracks[s.table.Cursor()]

			// Toggle selection
			if _, ok := s.selectedTracks[selectedTrack.ID]; ok {
				delete(s.selectedTracks, selectedTrack.ID)
			} else {
				s.selectedTracks[selectedTrack.ID] = selectedTrack
			}
			// Refresh the table view to show the change
			return s, s.refresh()

		case key.Matches(msg, s.keys.Create):
			if len(s.selectedTracks) > 0 {
				s.namingPlaylist = true
				s.nameForm = newNameForm()
				s.layout.ClearChildren()
				centeredForm := layout.CenteredBox(s.layout.Width(), s.layout.Height(), s.nameForm.View(), 0.5, 0.25)
				s.layout.AddChild(centeredForm)

				return s, s.nameForm.Init()
			}
		}
	}

	var tableCmd tea.Cmd
	s.table, tableCmd = s.table.Update(msg)
	cmds = append(cmds, tableCmd)

	s.layout, _ = newCreatePlaylistLayout(s.layout, s.table)

	return s, tea.Batch(cmds...)
}

func (s CreatePlaylistState) View() string {
	return s.layout.Render()
}

func (s CreatePlaylistState) Title() string {
	return "create playlist"
}

func (s CreatePlaylistState) State() states.StateType {
	return states.CreatePlaylist
}

func (s CreatePlaylistState) Help() string {
	return "create a playlist"
}

func (s CreatePlaylistState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}

	return states.Undefined, false
}

func (s CreatePlaylistState) Transition() states.State {
	s.nextState = states.Undefined
	return s
}
