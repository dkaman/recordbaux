package tui

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/help"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"
)

var (
	LoggerIsNilErr = errors.New("supplied slog logger is nil")
)

type MessageHandler func(tea.Model, tea.Msg) (tea.Model, tea.Cmd)

type Model struct {
	// global application config/state
	cfg      *config.Config
	keys     keyMap
	logger   *slog.Logger
	handlers *handlers.Registry

	ready         bool
	stateMachine  statemachine.Model
	topBarText    string
	statusBarText string
	helpVisible   bool

	width, height int
}

func New(c *config.Config, log *slog.Logger, svcs *services.AllServices) (Model, error) {
	var m Model

	if log == nil {
		return m, LoggerIsNilErr
	}

	h := help.New()
	h.Styles = style.DefaultHelpStyles()

	sm, err := statemachine.New(svcs, c, log)
	if err != nil {
		return m, fmt.Errorf("error creating state machine: %w", err)
	}

	m = Model{
		cfg:           c,
		keys:          defaultKeybinds(),
		handlers:      getHandlers(),
		helpVisible:   false,
		ready:         false,
		stateMachine:  sm,
		logger:        log.WithGroup("root"),
		topBarText:    "recordbaux - organize your record collection",
		statusBarText: fmt.Sprintf("current state: %s", sm.CurrentStateType()),
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	m.logger.Debug("root tui model init called")
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logger.Info("event received",
		slog.Any("event", fmt.Sprintf("%#v", msg)),
	)

	var cmds []tea.Cmd

	if handler, ok := m.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(m, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		m = model.(Model)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	var stateMachineCmd tea.Cmd
	m.stateMachine, stateMachineCmd = util.UpdateModel(m.stateMachine, msg)
	m.statusBarText = fmt.Sprintf("current state: %s", m.stateMachine.CurrentStateType())
	cmds = append(cmds, stateMachineCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.renderModel()
}

func (m Model) Help() string {
	return "global[ " +
		util.FmtKeymap(m.keys.ShortHelp()) + "] " +
		m.stateMachine.Help()
}
