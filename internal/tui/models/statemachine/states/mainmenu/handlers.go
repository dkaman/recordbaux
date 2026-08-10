package mainmenu

import (
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/playlist"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	tplaylist "github.com/dkaman/recordbaux/internal/tui/models/playlist"
)

func getHandlers() *handlers.Registry {
	r := handlers.NewRegistry()

	handlers.Register(r, handleTeaWindowSizeMsg)
	handlers.Register(r, handleTeaKeyPressMsg)
	handlers.Register(r, handleShelvesLoadedMsg)
	handlers.Register(r, handlePlaylistsLoadedMsg)

	return r
}

func handleTeaWindowSizeMsg(s MainMenuState, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd, tea.Msg) {
	var cmds []tea.Cmd

	s.width, s.height = msg.Width, msg.Height

	halfWidth := s.width / 2
	otherWidth := s.width - halfWidth

	var shelfUpdateCmds tea.Cmd

	shelfSizeMsg := tea.WindowSizeMsg{
		Width: halfWidth,
		Height: s.height,
	}

	s.shelves, shelfUpdateCmds = util.UpdateModel(s.shelves, shelfSizeMsg)
	cmds = append(cmds, shelfUpdateCmds)

	var playlistUpdateCmds tea.Cmd

	playlistSizeMsg := tea.WindowSizeMsg{
		Width: otherWidth,
		Height: s.height,
	}

	s.playlists, playlistUpdateCmds = util.UpdateModel(s.playlists, playlistSizeMsg)
	cmds = append(cmds, playlistUpdateCmds)

	s.logger.Debug("window sizes for child models",
		slog.Any("shelf width", s.shelves.Width()),
		slog.Any("shelf height", s.shelves.Height()),
		slog.Any("playlist width", s.playlists.Width()),
		slog.Any("playlist height", s.playlists.Height()),
	)

	return s, tea.Batch(cmds...), nil
}

func handleTeaKeyPressMsg(s MainMenuState, msg tea.KeyPressMsg) (tea.Model, tea.Cmd, tea.Msg) {
	if s.creating {
		return s, nil, msg
	}

	switch {
	case key.Matches(msg, s.keys.Quit):
		return s, nil, nil

	case key.Matches(msg, s.keys.SwitchFocus):
		if s.focus == shelvesView {
			s.focus = playlistsView
			s.playlists = s.playlists.Focus()
			s.shelves = s.shelves.Blur()
		} else {
			s.focus = shelvesView
			s.shelves = s.shelves.Focus()
			s.playlists = s.playlists.Blur()
		}

		return s, nil, nil

	case key.Matches(msg, s.keys.NewShelf):
		if s.focus == shelvesView {
			s.logger.Debug("create shelf selected")
			s.creating = true
			s.shelves = s.shelves.Blur()
			s.playlists = s.playlists.Blur()
			s.createShelfForm = newCreateShelfForm()
			return s, s.createShelfForm.Init(), nil
		} else {
			s.logger.Debug("create playlist selected")
			return s, tcmds.Transition(
				states.CreatePlaylist,
				nil,
				[]tea.Cmd{s.svcs.GetAllTracksCmd()},
			), nil
		}

	case key.Matches(msg, s.keys.Select):
		if s.focus == shelvesView {
			if sel, ok := s.shelves.SelectedItem().(shelf.Model); ok {
				return s, tcmds.Transition(
					states.LoadedShelf,
					nil,
					[]tea.Cmd{shelf.WithPhysicalShelf(sel.PhysicalShelf())},
				), nil
			}
		} else {
			if sel, ok := s.playlists.SelectedItem().(tplaylist.Model); ok {
				return s, tcmds.Transition(
					states.LoadedPlaylist,
					nil,
					[]tea.Cmd{playlist.WithPhysicalPlaylist(sel.PhysicalPlaylist())},
				), nil
			}
		}
	}

	// pass key message through if we aren't handling it
	return s, nil, msg
}

func handleShelvesLoadedMsg(s MainMenuState, msg services.ShelvesLoadedMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.logger.Debug("refreshing shelves from service")

	shlvs := msg.Shelves
	items := make([]list.Item, len(shlvs))

	for i, sh := range shlvs {
		items[i] = shelf.New(sh, s.logger)
	}

	s.shelves.SetItems(items)

	return s, nil, nil
}

func handlePlaylistsLoadedMsg(s MainMenuState, msg services.PlaylistsLoadedMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.logger.Debug("refreshing playlists from service")

	playlists := msg.Playlists
	playlistItems := make([]list.Item, len(playlists))

	for i, p := range playlists {
		playlistItems[i] = tplaylist.New(p)
	}

	s.playlists.SetItems(playlistItems)

	return s, nil, nil
}
