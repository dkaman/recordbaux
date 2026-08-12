package loadedplaylist

import (
	"log/slog"

	"charm.land/bubbles/v2/key"

	tea "charm.land/bubbletea/v2"

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
		s.logger.Info("checking out playlist")
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
