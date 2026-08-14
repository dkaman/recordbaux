package tui

import (
	"errors"
	"fmt"
	"log/slog"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"
)

var (
	LoggerIsNilErr = errors.New("supplied slog logger is nil")
)

// root tui model
type Model struct {
	// global application config/state
	cfg    *config.Config
	keys   keyMap
	logger *slog.Logger
	// handlers *handlers.Registry

	ready         bool
	stateMachine  statemachine.Model
	topBarText    string
	statusBarText string
	helpVisible   bool

	width, height int
}

// local type def to define keybindings for this model
type keyMap struct {
	ToggleHelp key.Binding
	Quit       key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		ToggleHelp: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("C-c", "quit"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.ToggleHelp,
		k.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ToggleHelp, k.Quit},
	}
}

func New(c *config.Config, log *slog.Logger, svcs *services.AllServices) (Model, error) {
	var m Model

	if log == nil {
		return m, LoggerIsNilErr
	}

	sm, err := statemachine.New(svcs, c, log)
	if err != nil {
		return m, fmt.Errorf("error creating state machine: %w", err)
	}

	m = Model{
		cfg:  c,
		keys: defaultKeybinds(),
		// handlers:      getHandlers(),
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
	return m.stateMachine.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var passthruMsg tea.Msg

	m.logger.Info("event received",
		slog.Any("event", fmt.Sprintf("%#v", msg)),
	)

	passthruMsg = msg

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// this code updates the root element's size, then constructs a
		// size msg that represents the whole size of the viewport to
		// pass thru to child model updates
		m.width, m.height = msg.Width, msg.Height

		if !m.ready {
			m.ready = true
		}

		numBars := 2
		if m.helpVisible {
			numBars = 3
		}

		// modify window size msg and pass thru
		passthruMsg = tea.WindowSizeMsg{
			Width:  m.width - 2,
			Height: m.height - numBars - 2,
		}

	case tea.KeyPressMsg:
		// this branch also passes through the key message if there is
		// no match so the child models can handle the keys
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.ToggleHelp):
			m.helpVisible = !m.helpVisible
			return m, func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
			}
		}
	}

	var stateMachineCmd tea.Cmd
	m.stateMachine, stateMachineCmd = util.UpdateModel(m.stateMachine, passthruMsg)
	cmds = append(cmds, stateMachineCmd)

	m.statusBarText = fmt.Sprintf("current state: %s", m.stateMachine.CurrentStateType())

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("\n initializing...")
	}

	numBars := 2
	if m.helpVisible {
		numBars = 3
	}

	barStyle := style.BarStyle.
		Width(m.width).
		Height(1)

	helpStyle := style.HelpBarStyle.
		Width(m.width).
		Height(1)

	viewportStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width).
		Height(max(0, m.height-numBars))

	topBar := barStyle.Render(m.topBarText)
	statusBar := barStyle.Render(m.statusBarText)
	viewPort := viewportStyle.Render(m.stateMachine.View().Content)

	var content string

	if m.helpVisible {
		helpBar := helpStyle.Render(m.Help())
		content = lipgloss.JoinVertical(lipgloss.Left,
			topBar,
			viewPort,
			helpBar,
			statusBar,
		)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Left,
			topBar,
			viewPort,
			statusBar,
		)
	}

	return tea.NewView(content)
}

func (m Model) Help() string {
	return "global[ " +
		util.FmtKeymap(m.keys.ShortHelp()) + "] " +
		m.stateMachine.Help()
}
