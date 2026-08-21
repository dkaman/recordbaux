package createplaylist

import (
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"

	"github.com/dkaman/recordbaux/internal/db/playlist"
	"github.com/dkaman/recordbaux/internal/db/track"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	ttrack "github.com/dkaman/recordbaux/internal/tui/models/track"
)

type CreatePlaylistState struct {
	logger *slog.Logger
	keys   keyMap

	shelfIDs       []uint
	list           list.Model
	namingPlaylist bool
	nameForm       *form
	playlistName   string

	width, height int
}

func New(log *slog.Logger, shelfIDs []uint) (CreatePlaylistState, error) {
	s := CreatePlaylistState{
		keys:           defaultKeybinds(),
		namingPlaylist: false,
		shelfIDs:       shelfIDs,
	}

	s.logger = log.WithGroup("createplayliststate")

	delegate := trackDelegate{}
	trackList := list.New([]list.Item{}, delegate, 0, 0)
	trackList.Styles = style.DefaultListStyles()
	trackList.Title = "select tracks for new playlist"
	s.list = trackList

	return s, nil
}

func (s CreatePlaylistState) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, id := range s.shelfIDs {
		cmds = append(cmds, tcmds.GetAllTracksFromShelfCmd(id))
	}

	cmds = append(cmds, tcmds.RefreshWindowSizeCmd())

	return tea.Batch(cmds...)
}

func (s CreatePlaylistState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	passthru := msg

	if s.namingPlaylist {
		fModel, formUpdatesCmds := s.nameForm.Update(msg)
		if f, ok := fModel.(*form); ok {
			s.nameForm = f
		}
		cmds = append(cmds, formUpdatesCmds)

		if s.nameForm.State == huh.StateCompleted {
			s.logger.Debug("end naming playlist")
			name := s.nameForm.Name()
			newPlaylist := &playlist.Entity{
				Name:   name,
				Tracks: make([]*track.Entity, 0),
			}

			for _, item := range s.list.Items() {
				if trackModel, ok := item.(ttrack.Model); ok && trackModel.Selected {
					newPlaylist.Tracks = append(newPlaylist.Tracks, trackModel.PhysicalTrack())
				}
			}

			cmds = append(cmds, tcmds.SavePlaylistCmd(newPlaylist))

			items := s.list.Items()
			for i, item := range items {
				if trackModel, ok := item.(ttrack.Model); ok && trackModel.Selected {
					trackModel.Selected = false
					items[i] = trackModel
				}
			}

			s.list.SetItems(items)
			s.namingPlaylist = false

			cmds = append(cmds, func() tea.Msg {
				return tcmds.TransitionToMainMenuMsg{}
			})
		}

		return s, tea.Batch(cmds...)
	}


	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height

	case tcmds.ShelfAllTracksMsg:
		s.logger.Debug("refreshing tracks from service")
		tracks := msg.Tracks
		items := make([]list.Item, len(tracks))

		for i, t := range tracks {
			items[i] = ttrack.New(t)
		}

		s.list.SetItems(items)

		return s, nil

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			return s, func() tea.Msg {
				return tcmds.TransitionToMainMenuMsg{}
			}

		case key.Matches(msg, s.keys.Select):
			if i, ok := s.list.SelectedItem().(ttrack.Model); ok {
				s.logger.Debug("track selected", slog.Any("track", i))
				i.Selected = !i.Selected
				cmd := s.list.SetItem(s.list.Index(), i)
				return s, cmd
			}

			return s, nil

		case key.Matches(msg, s.keys.Create):
			var selectedCount int

			for _, item := range s.list.Items() {
				if trackModel, ok := item.(ttrack.Model); ok && trackModel.Selected {
					selectedCount++
				}
			}

			if selectedCount > 0 {
				s.namingPlaylist = true
				s.nameForm = newNameForm()
				return s, s.nameForm.Init()
			}
		}
	}

	var listCmd tea.Cmd
	s.list, listCmd = s.list.Update(passthru)
	cmds = append(cmds, listCmd)
	return s, tea.Batch(cmds...)
}

func (s CreatePlaylistState) View() tea.View {
	return tea.NewView(s.renderModel())
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

func (s CreatePlaylistState) Type() states.StateType {
	return states.CreatePlaylist
}
