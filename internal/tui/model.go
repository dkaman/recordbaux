package tui

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"

	kfmt "github.com/dkaman/recordbaux/internal/tui/key"
)

var (
	LoggerIsNilErr = errors.New("supplied slog logger is nil")
)

type Model struct {
	// global application config/state
	cfg    *config.Config
	svcs   *services.AllServices
	keys   keyMap
	logger *slog.Logger

	// tea models
	stateMachine statemachine.Model

	// styling/layout
	helpVisible   bool
	topBarText    string
	statusBarText string

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
		svcs:          svcs,
		keys:          defaultKeybinds(),
		helpVisible:   false,
		stateMachine:  sm,
		logger:        log.WithGroup("root"),
		topBarText:    "recordbaux - organize your record collection",
		statusBarText: fmt.Sprintf("current state: %s", m.stateMachine.CurrentStateType()),
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	m.logger.Debug("root tui model init called")
	return m.stateMachine.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.logger.Info("event received",
		slog.Any("event", fmt.Sprintf("%#v", msg)),
	)

	// update bars
	m.statusBarText = fmt.Sprintf("current state: %s", m.stateMachine.CurrentStateType())

	processedMsg := msg

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		numBars := 2
		if m.helpVisible {
			numBars = 3
		}

		smWidth := m.width - 2
		smHeight := m.height - numBars - 2

		processedMsg = tea.WindowSizeMsg{
			Width:  smWidth,
			Height: smHeight,
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.ToggleHelp):
			m.helpVisible = !m.helpVisible
			return m, nil
		}
	}

	// update state machine
	stateMachineModel, stateMachineCmds := m.stateMachine.Update(processedMsg)
	if sm, ok := stateMachineModel.(statemachine.Model); ok {
		m.stateMachine = sm
	}
	cmds = append(cmds, stateMachineCmds)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.renderModel()
}

func (m Model) Help() string {
	return "global[ " +
		kfmt.FmtKeymap(m.keys.ShortHelp()) + "] " +
		m.stateMachine.Help()
}
