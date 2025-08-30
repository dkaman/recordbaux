package loadedplaylist

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/table"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	keyFmt "github.com/dkaman/recordbaux/internal/tui/key"
)

type refreshMsg struct{}

type PlaylistLoadedState struct {
	playlistService *services.PlaylistService
	recordService   *services.RecordService
	nextState       states.StateType
	keys            keyMap
	logger          *slog.Logger
	trackTable      table.Model
	width, height   int
}

func New(p *services.PlaylistService, r *services.RecordService, log *slog.Logger) PlaylistLoadedState {
	columns := []table.Column{
		{Title: "Position", Width: 10},
		{Title: "Title", Width: 50},
		{Title: "Duration", Width: 10},
		{Title: "Key", Width: 8},
		{Title: "BPM", Width: 8},
	}
	tbl := table.New(table.WithColumns(columns), table.WithFocused(true))
	tbl.SetStyles(style.DefaultTableStyles())

	return PlaylistLoadedState{
		nextState: states.Undefined,
		playlistService: p,
		recordService:   r,
		keys:            defaultKeybinds(),
		logger:          log.WithGroup("playlistloaded"),
		trackTable:      tbl,
	}
}

func (s PlaylistLoadedState) Init() tea.Cmd {
	return func() tea.Msg { return refreshMsg{} }
}

func (s PlaylistLoadedState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var tableCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.trackTable.SetWidth(msg.Width - 2)
		s.trackTable.SetHeight(msg.Height - 2)

	case refreshMsg:
		playlist := s.playlistService.CurrentPlaylist
		var rows []table.Row
		if playlist != nil {
			for _, t := range playlist.Tracks {
				rows = append(rows, table.Row{t.Position, t.Title, t.Duration, t.Key, fmt.Sprintf("%d", t.BPM)})
			}
		}
		s.trackTable.SetRows(rows)
		return s, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			s.nextState = states.MainMenu
			return s, nil
		case key.Matches(msg, s.keys.Checkout):
			s.logger.Info("checking out playlist")
			playlist := s.playlistService.CurrentPlaylist
			if playlist != nil && len(playlist.Tracks) > 0 {
				return s, tcmds.CheckoutPlaylistCmd(s.recordService.Records, playlist, s.logger)
			}
		}
	}

	s.trackTable, tableCmd = s.trackTable.Update(msg)
	cmds = append(cmds, tableCmd)

	return s, tea.Batch(cmds...)
}

func (s PlaylistLoadedState) View() string {
	return s.renderModel()
}

func (s PlaylistLoadedState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}
	return states.Undefined, false
}

func (s PlaylistLoadedState) Transition() states.State {
	s.nextState = states.Undefined
	return s
}

func (s PlaylistLoadedState) Help() string {
	return keyFmt.FmtKeymap(s.keys.ShortHelp())
}
