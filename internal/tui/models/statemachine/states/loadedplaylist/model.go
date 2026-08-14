package loadedplaylist

import (
	"log/slog"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/playlist"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type LoadedPlaylistState struct {
	svcs   *services.AllServices
	keys   keyMap
	logger *slog.Logger

	playlist playlist.Model

	width, height int
}

func New(svcs *services.AllServices, log *slog.Logger) LoadedPlaylistState {
	return LoadedPlaylistState{
		svcs:     svcs,
		keys:     defaultKeybinds(),
		logger:   log.WithGroup("playlistloaded"),
		playlist: playlist.New(),
	}
}

func (s LoadedPlaylistState) Init() tea.Cmd {
	return nil
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
			return s, tcmds.Transition(states.MainMenu, nil, nil)
		case key.Matches(msg, s.keys.Checkout):
			return s, s.svcs.SetCheckoutCmd(s.playlist.PhysicalPlaylist(), true)
		case key.Matches(msg, s.keys.Checkin):
			return s, s.svcs.SetCheckoutCmd(s.playlist.PhysicalPlaylist(), false)
		}

		passthru = msg

	case playlist.LoadPlaylistMsg:
		s.playlist.SetEntity(msg.Phy)
		return s, nil

	case services.PlaylistCheckoutMsg:
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
