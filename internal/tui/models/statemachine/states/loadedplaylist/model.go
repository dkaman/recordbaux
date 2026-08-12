package loadedplaylist

import (
	"log/slog"

	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/playlist"
	"github.com/dkaman/recordbaux/internal/tui/util"
)

type LoadedPlaylistState struct {
	svcs     *services.AllServices
	keys     keyMap
	logger   *slog.Logger
	handlers *handlers.Registry

	playlist   playlist.Model

	width, height int
}

func New(svcs *services.AllServices, log *slog.Logger) LoadedPlaylistState {
	return LoadedPlaylistState{
		svcs:     svcs,
		keys:     defaultKeybinds(),
		logger:   log.WithGroup("playlistloaded"),
		handlers: getHandlers(),
		playlist: playlist.New(),
	}
}

func (s LoadedPlaylistState) Init() tea.Cmd {
	return nil
}

func (s LoadedPlaylistState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if handler, ok := s.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(s, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		s = model.(LoadedPlaylistState)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	var playlistCmd tea.Cmd
	s.playlist, playlistCmd = s.playlist.Update(msg)
	cmds = append(cmds, playlistCmd)

	return s, tea.Batch(cmds...)
}

func (s LoadedPlaylistState) View() tea.View {
	return s.renderModel()
}

func (s LoadedPlaylistState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}
