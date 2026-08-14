package loadedbin

import (
	"log/slog"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"

	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/record"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type LoadedBinState struct {
	svcs     *services.AllServices
	keys     keyMap
	logger   *slog.Logger

	bin            bin.Model
	records        table.Model
	selectedRecord record.Model
	cursorIndex    int

	width, height int
}

// New constructs a LoadedBinState ready to receive a LoadShelfMsg
func New(svcs *services.AllServices, log *slog.Logger) LoadedBinState {
	h := help.New()
	h.Styles = style.DefaultHelpStyles()

	t := table.New()

	return LoadedBinState{
		svcs:     svcs,
		keys:     defaultKeybinds(),
		logger:   log.WithGroup("loadedbin"),

		records: t,
	}
}

func (s LoadedBinState) Init() tea.Cmd {
	s.logger.Debug("loadedbin state init")
	return nil
}

func (s LoadedBinState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var passthru tea.Msg

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height
		passthru = tea.WindowSizeMsg{Width: s.width / 2, Height: s.height}

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			return s, tcmds.Transition(states.LoadedShelf, nil, nil)
		}

		passthru = msg

	case bin.LoadBinMsg:
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

			if r.CheckedOut {
				name = "[OUT] " + name
			}

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

		return s, nil
	}

	oldIdx := s.records.Cursor()

	var tableUpdateCmd tea.Cmd
	s.records, tableUpdateCmd = s.records.Update(passthru)
	cmds = append(cmds, tableUpdateCmd)

	idx := s.records.Cursor()
	if idx != oldIdx {
		selectedPhysicalRecord := s.bin.PhysicalBin().Records[idx]
		s.selectedRecord = record.New(selectedPhysicalRecord)
		s.cursorIndex = idx
	}

	var recordCmd tea.Cmd
	s.selectedRecord, recordCmd = util.UpdateModel(s.selectedRecord, passthru)
	cmds = append(cmds, recordCmd)

	return s, tea.Batch(cmds...)
}

func (s LoadedBinState) View() tea.View {
	return tea.NewView(s.renderModel())
}

func (s LoadedBinState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}
