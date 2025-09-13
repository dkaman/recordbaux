package loadedbin

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/table"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/record"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"
)

type LoadedBinState struct {
	svcs     *services.AllServices
	keys     keyMap
	logger   *slog.Logger
	handlers *handlers.Registry

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
		handlers: getHandlers(),

		records: t,
	}
}

func (s LoadedBinState) Init() tea.Cmd {
	s.logger.Debug("loadedbin state init")
	return nil
}

func (s LoadedBinState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if handler, ok := s.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(s, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		s = model.(LoadedBinState)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	oldIdx := s.records.Cursor()

	var tableUpdateCmd tea.Cmd
	s.records, tableUpdateCmd = s.records.Update(msg)
	cmds = append(cmds, tableUpdateCmd)

	idx := s.records.Cursor()
	if idx != oldIdx {
		selectedPhysicalRecord := s.bin.PhysicalBin().Records[idx]
		s.selectedRecord = record.New(selectedPhysicalRecord)
		s.cursorIndex = idx
	}

	var recordCmd tea.Cmd
	s.selectedRecord, recordCmd = util.UpdateModel(s.selectedRecord, msg)
	cmds = append(cmds, recordCmd)

	return s, tea.Batch(cmds...)
}

func (s LoadedBinState) View() string {
	return s.renderModel()
}

func (s LoadedBinState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}
