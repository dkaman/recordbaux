package loadedplaylist

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/table"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"
)

type LoadedPlaylistState struct {
	svcs *services.AllServices
	keys            keyMap
	logger          *slog.Logger
	handlers *handlers.Registry

	trackTable      table.Model

	width, height   int
}

func New(svcs *services.AllServices, log *slog.Logger) LoadedPlaylistState {
	columns := []table.Column{
		{Title: "Position", Width: 10},
		{Title: "Title", Width: 50},
		{Title: "Duration", Width: 10},
		{Title: "Key", Width: 8},
		{Title: "BPM", Width: 8},
	}
	tbl := table.New(table.WithColumns(columns), table.WithFocused(true))
	tbl.SetStyles(style.DefaultTableStyles())

	return LoadedPlaylistState{
		svcs: svcs,
		keys:            defaultKeybinds(),
		logger:          log.WithGroup("playlistloaded"),
		handlers: getHandlers(),

		trackTable:      tbl,
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

	var tableCmd tea.Cmd
	s.trackTable, tableCmd = s.trackTable.Update(msg)
	cmds = append(cmds, tableCmd)

	return s, tea.Batch(cmds...)
}

func (s LoadedPlaylistState) View() string {
	return s.renderModel()
}

func (s LoadedPlaylistState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}
