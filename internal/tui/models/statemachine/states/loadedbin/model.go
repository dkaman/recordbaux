package loadedbin

import (
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"

	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/record"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type LoadedBinState struct {
	logger *slog.Logger
	keys   keyMap

	binID          uint
	bin            bin.Model
	records        table.Model
	selectedRecord record.Model
	cursorIndex    int

	width, height int
}

// New constructs a LoadedBinState ready to receive a LoadShelfMsg
func New(log *slog.Logger, binID uint) (LoadedBinState, error) {
	s := LoadedBinState{
		keys: defaultKeybinds(),
	}

	s.logger = log.WithGroup("loadedbin")
	s.records = table.New()
	s.binID = binID

	columns := []table.Column{
		{Title: "catalog no.", Width: 15},
		{Title: "release name", Width: 50},
		{Title: "artist", Width: 30},
	}

	s.records = table.New(
		table.WithColumns(columns),
		table.WithStyles(style.DefaultTableStyles()),
	)

	return s, nil
}

func (s LoadedBinState) Init() tea.Cmd {
	s.logger.Debug("loadedbin state init")
	return tea.Sequence(
		tcmds.GetBinCmd(s.binID),
		tcmds.RefreshWindowSizeCmd(),
	)
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
			return s, func() tea.Msg {
				return tcmds.TransitionToLoadedShelfMsg{
					ShelfID: s.bin.PhysicalBin().ShelfID,
				}
			}
		}

		passthru = msg

	case tcmds.BinsLoadedMsg:
		if msg.Err != nil {
			return s, nil
		}

		if len(msg.Bins) != 1 {
			return s, nil
		}

		entity := msg.Bins[0]
		s.logger.Debug("bin entity", slog.Any("e", entity))

		s.bin = bin.New(entity, bin.Style{})


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

		s.records.SetRows(rows)
		s.records.Focus()

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

func (s LoadedBinState) Type() states.StateType {
	return states.LoadedBin
}
