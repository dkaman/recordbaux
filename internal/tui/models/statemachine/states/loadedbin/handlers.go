package loadedbin

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/table"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/record"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

func getHandlers() *handlers.Registry {
	r := handlers.NewRegistry()

	handlers.Register(r, handleTeaWindowSizeMsg)
	handlers.Register(r, handleTeaKeyPressMsg)
	handlers.Register(r, handleLoadBinMsg)

	return r
}

func handleTeaWindowSizeMsg(s LoadedBinState, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.width, s.height = msg.Width, msg.Height
	msg = tea.WindowSizeMsg{Width: s.width / 2, Height: s.height}
	return s, nil, msg
}

func handleTeaKeyPressMsg(s LoadedBinState, msg tea.KeyPressMsg) (tea.Model, tea.Cmd, tea.Msg) {
	switch {
	case key.Matches(msg, s.keys.Back):
		return s, tcmds.Transition(states.LoadedShelf, nil, nil), nil
	}
	return s, nil, msg
}

func handleLoadBinMsg(s LoadedBinState, msg bin.LoadBinMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.bin = bin.New(msg.Phy, bin.Style{})

	columns := []table.Column{
		{Title: "catalog no.", Width: 15},
		{Title: "release name", Width: 50},
		{Title: "artist", Width: 30},
	}

	var rows []table.Row

	for _, r := range s.bin.PhysicalBin().Records {
		catno := r.CatalogNumber
		name := r.Title
		artist := r.Artists[0]
		row := table.Row{catno, name, artist}
		rows = append(rows, row)
	}

	s.records = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithStyles(style.DefaultTableStyles()),
	)

	s.cursorIndex = 0
	if len(s.bin.PhysicalBin().Records) > 0 {
		// Create the initial record model for the first item
		initialRecord := s.bin.PhysicalBin().Records[s.cursorIndex]
		s.selectedRecord = record.New(initialRecord)
	}

	return s, nil, nil
}
