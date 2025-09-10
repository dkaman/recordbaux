package loadedplaylist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/table"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/models/playlist"

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
	return s, nil, nil
}

func handleTeaKeyPressMsg(s LoadedPlaylistState, msg tea.KeyPressMsg) (tea.Model, tea.Cmd, tea.Msg) {
	switch {
	case key.Matches(msg, s.keys.Back):
		return s, tcmds.Transition(states.MainMenu, nil, nil), nil
		// case key.Matches(msg, s.keys.Checkout):
		// 	s.logger.Info("checking out playlist")
		// 	playlist := s.playlistService.CurrentPlaylist
		// 	if playlist != nil && len(playlist.Tracks) > 0 {
		// 		return s, tcmds.CheckoutPlaylistCmd(s.recordService.Records, playlist, s.logger)
		// 	}
	}

	return s, nil, msg
}

func handleLoadPlaylistMsg(s LoadedPlaylistState, msg playlist.LoadPlaylistMsg) (tea.Model, tea.Cmd, tea.Msg) {
	playlist := msg.Phy

	var rows []table.Row
	if playlist != nil {
		for _, t := range playlist.Tracks {
			rows = append(rows, table.Row{t.Position, t.Title, t.Duration, t.Key, fmt.Sprintf("%d", t.BPM)})
		}
	}
	s.trackTable.SetRows(rows)
	return s, nil, nil
}
