package createplaylist

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	ttrack "github.com/dkaman/recordbaux/internal/tui/models/track"
)

func getHandlers() *handlers.Registry {
	r := handlers.NewRegistry()

	handlers.Register(r, handleTeaWindowSizeMsg)
	handlers.Register(r, handleAllTracksLoadedMsg)
	handlers.Register(r, handleTeaKeyPressMsg)

	return r
}

func handleTeaWindowSizeMsg(s CreatePlaylistState, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.width, s.height = msg.Width, msg.Height
	return s, nil, nil
}

func handleAllTracksLoadedMsg(s CreatePlaylistState, msg services.AllTracksLoadedMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.logger.Debug("refreshing tracks from service")
	tracks := msg.Tracks
	items := make([]list.Item, len(tracks))

	for i, t := range tracks {
		items[i] = ttrack.New(t)
	}

	s.list.SetItems(items)
	return s, nil, nil
}

func handleTeaKeyPressMsg(s CreatePlaylistState, msg tea.KeyPressMsg) (tea.Model, tea.Cmd, tea.Msg) {
	switch {
	case key.Matches(msg, s.keys.Back):
		return s, tcmds.Transition(states.MainMenu, nil, nil), nil

	case key.Matches(msg, s.keys.Select):
		if i, ok := s.list.SelectedItem().(ttrack.Model); ok {
			s.logger.Debug("track selected", slog.Any("track", i))
			i.Selected = !i.Selected
			cmd := s.list.SetItem(s.list.Index(), i)
			return s, cmd, nil
		}
		return s, nil, nil

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
			return s, s.nameForm.Init(), nil
		}
	}

	return s, nil, msg
}
