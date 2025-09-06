package loadedbin

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/table"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/record"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	keyFmt "github.com/dkaman/recordbaux/internal/tui/key"
)

type LoadedBinState struct {
	svcs      *services.AllServices
	nextState states.StateType
	bin       bin.Model
	keys      keyMap
	logger    *slog.Logger

	records        table.Model
	selectedRecord record.Model
	width, height  int
}

// New constructs a LoadedBinState ready to receive a LoadShelfMsg
func New(svcs *services.AllServices, log *slog.Logger) LoadedBinState {
	h := help.New()
	h.Styles = style.DefaultHelpStyles()

	t := table.New()

	return LoadedBinState{
		svcs:      svcs,
		nextState: states.Undefined,
		keys:      defaultKeybinds(),
		logger:    log.WithGroup("loadedbin"),
		records:   t,
	}
}

func (s LoadedBinState) Init() tea.Cmd {
	s.logger.Debug("loadedbin state init")
	return nil
}

func (s LoadedBinState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	processedMsg := msg

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w, h := msg.Width, msg.Height
		s.width = w
		s.height = h

		processedMsg = tea.WindowSizeMsg{Width: w / 2, Height: h}

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
			row := table.Row{catno, name, artist}
			rows = append(rows, row)
		}

		s.records = table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithStyles(style.DefaultTableStyles()),
		)

		return s, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			return s, tcmds.WithNextState(states.LoadedShelf, nil, nil)
		}
	}

	tableModel, tableUpdateCmds := s.records.Update(processedMsg)
	s.records = tableModel
	cmds = append(cmds,
		tableUpdateCmds,
	)

	idx := s.records.Cursor()

	// TODO fix this bug, not rendering
	if b := s.bin.PhysicalBin(); b != nil {
		r := s.bin.PhysicalBin().Records[idx]

		recordModel, recordCmds := record.New(r).Update(processedMsg)
		if rec, ok := recordModel.(record.Model); ok {
			s.selectedRecord = rec
		}
		cmds = append(cmds, recordCmds)
	}

	return s, tea.Batch(cmds...)
}

func (s LoadedBinState) View() string {
	return s.renderModel()
}

func (s LoadedBinState) Help() string {
	return keyFmt.FmtKeymap(s.keys.ShortHelp())
}

func (s LoadedBinState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}

	return states.Undefined, false
}

func (s LoadedBinState) Transition() states.State {
	s.nextState = states.Undefined
	return s
}
