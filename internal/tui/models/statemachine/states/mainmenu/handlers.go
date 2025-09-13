package mainmenu

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/playlist"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

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
	s.width, s.height = msg.Width, msg.Height
	return s, nil, nil
}

func handleTeaKeyPressMsg(s MainMenuState, msg tea.KeyPressMsg) (tea.Model, tea.Cmd, tea.Msg) {
	switch {
	case key.Matches(msg, s.keys.Quit):
		return s, nil, nil

	case key.Matches(msg, s.keys.SwitchFocus):
		if s.focus == shelvesView {
			s.focus = playlistsView
		} else {
			s.focus = shelvesView
		}

		return s, nil, nil

	case key.Matches(msg, s.keys.NewShelf):
		if s.focus == shelvesView {
			s.logger.Debug("create shelf selected")
			s.creating = true
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
