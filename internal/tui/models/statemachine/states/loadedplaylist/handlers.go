package loadedplaylist

import (
	"log/slog"

	"charm.land/bubbles/v2/key"

	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/playlist"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

func getHandlers() *handlers.Registry {
	r := handlers.NewRegistry()

	handlers.Register(r, handleTeaWindowSizeMsg)
	handlers.Register(r, handleTeaKeyPressMsg)
	handlers.Register(r, handleLoadPlaylistMsg)
	handlers.Register(r, handlePlaylistCheckoutMsg)

	return r
}

func handleTeaWindowSizeMsg(s LoadedPlaylistState, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.width, s.height = msg.Width, msg.Height
	return s, nil, msg
}

func handleTeaKeyPressMsg(s LoadedPlaylistState, msg tea.KeyPressMsg) (tea.Model, tea.Cmd, tea.Msg) {
	switch {
	case key.Matches(msg, s.keys.Back):
		return s, tcmds.Transition(states.MainMenu, nil, nil), nil
	case key.Matches(msg, s.keys.Checkout):
		return s, s.svcs.SetCheckoutCmd(s.playlist.PhysicalPlaylist(), true), nil
	case key.Matches(msg, s.keys.Checkin):
		return s , s.svcs.SetCheckoutCmd(s.playlist.PhysicalPlaylist(), false), nil
	}

	return s, nil, msg
}

func handleLoadPlaylistMsg(s LoadedPlaylistState, msg playlist.LoadPlaylistMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.logger.Debug("entity message received",
		slog.Any("tracks", msg.Phy.Tracks),
	)
	s.playlist.SetEntity(msg.Phy)
	return s, nil, nil
}

func handlePlaylistCheckoutMsg(s LoadedPlaylistState, msg services.PlaylistCheckoutMsg) (tea.Model, tea.Cmd, tea.Msg) {
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

	return s, nil, nil
}
