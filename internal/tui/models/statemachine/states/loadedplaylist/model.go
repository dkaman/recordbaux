package loadedplaylist

import (
	"log/slog"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/models/playlist"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type LoadedPlaylistState struct {
	width, height int

	logger *slog.Logger
	keys   keyMap

	playlistID uint
	playlist   playlist.Model
}

func New(log *slog.Logger, pID uint) (LoadedPlaylistState, error) {
	s := LoadedPlaylistState{
		keys: defaultKeybinds(),
	}

	s.logger = log.WithGroup("playlistloaded")
	s.playlistID = pID
	s.playlist = playlist.New()

	return s, nil
}

func (s LoadedPlaylistState) Init() tea.Cmd {
	return tea.Sequence(
		tcmds.GetPlaylistCmd(s.playlistID),
		tcmds.RefreshWindowSizeCmd(),
	)
}

func (s LoadedPlaylistState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var passthru tea.Msg

	passthru = msg

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height
		passthru = msg

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			return s, func() tea.Msg {
				return tcmds.TransitionToMainMenuMsg{}
			}
		case key.Matches(msg, s.keys.Checkout):
			return s, tcmds.SetCheckoutCmd(s.playlist.PhysicalPlaylist(), true)
		case key.Matches(msg, s.keys.Checkin):
			return s, tcmds.SetCheckoutCmd(s.playlist.PhysicalPlaylist(), false)
		}

		passthru = msg

	case tcmds.PlaylistsLoadedMsg:
		s.logger.Debug("loaded playlist", slog.Any("entity", msg.Playlists[0]))
		if msg.Err != nil {
			s.logger.Error("error loading playlist",
				slog.Any("err", msg.Err),
			)
			return s, nil
		}

		if len(msg.Playlists) == 1 {
			s.playlist.SetEntity(msg.Playlists[0])
			return s, nil
		} else {
			s.logger.Warn("more than one playlist was returned from the db, this shouldn't happen")
		}

		return s, nil

	case tcmds.PlaylistCheckoutMsg:
		if msg.Err != nil {
			s.logger.Error("failed to check out playlist",
				slog.String("error", msg.Err.Error()),
			)
		} else {
			if msg.Status {
				s.logger.Info("playlist checked out successfully")
			} else {

				s.logger.Info("playlist checked in successfully")
			}
		}

		return s, nil
	}

	var playlistCmd tea.Cmd
	s.playlist, playlistCmd = s.playlist.Update(passthru)
	cmds = append(cmds, playlistCmd)

	return s, tea.Batch(cmds...)
}

func (s LoadedPlaylistState) View() tea.View {
	return s.playlist.View()
}

func (s LoadedPlaylistState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}

func (s LoadedPlaylistState) Type() states.StateType {
	return states.LoadedPlaylist
}
