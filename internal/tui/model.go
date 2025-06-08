package tui

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/div"

	keyFmt "github.com/dkaman/recordbaux/internal/tui/key"
)

type Model struct {
	// global application config/state
	cfg  *config.Config
	app  *app.App
	keys keyMap

	logger *slog.Logger

	// tea models
	stateMachine statemachine.Model

	// styling/layout
	helpVisible bool
	layout      *div.Div
}

func New(c *config.Config, log *slog.Logger) Model {
	h := help.New()
	h.Styles = style.DefaultHelpStyles()

	l, _ := newTUILayout()
	vp := l.Find("viewport")

	a := app.NewApp()

	if log == nil {
		// TODO handle errors
		f, _ := os.OpenFile("./logs/recordbaux.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644);
		log = slog.New(slog.NewTextHandler(f, nil))
	}

	sm, _ := statemachine.New(a, c, vp, log)

	m := Model{
		cfg:          c,
		app:          a,
		keys:         defaultKeybinds(),
		helpVisible:  false,
		layout:       l,
		stateMachine: sm,
		logger: log,
	}

	_ = addTopBarText(l, "recordbaux - organize your record collection")
	_ = addStatusBarText(m.layout, fmt.Sprintf("current state: %s", m.stateMachine.CurrentStateType()))

	return m
}

func (m Model) Init() tea.Cmd {
	return m.stateMachine.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.logger.Info("event received",
		slog.Any("event", fmt.Sprintf("%#v", msg)),
	)

	// update bars
	statusBarText := fmt.Sprintf("current state: %s", m.stateMachine.CurrentStateType())
	_ = addStatusBarText(m.layout, statusBarText)

	_ = addHelpBarText(m.layout, m.Help())

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w, h := msg.Width, msg.Height
		m.layout.Resize(w, h)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.ToggleHelp):
			m.helpVisible = !m.helpVisible

			helpBar := m.layout.Find("helpbar")
			if m.helpVisible {
				helpBar.Show()
			} else {
				helpBar.Hide()
			}

			w, h := m.layout.Width(), m.layout.Height()
			m.layout.Resize(w, h)

			return m, nil
		}
	}

	// update state machine
	stateMachineModel, stateMachineCmds := m.stateMachine.Update(msg)
	if sm, ok := stateMachineModel.(statemachine.Model); ok {
		m.stateMachine = sm
	}
	cmds = append(cmds, stateMachineCmds)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.layout.Render()
}

func (m Model) Help() string {
	return "global: " +
		keyFmt.FmtKeymap(m.keys.ShortHelp()) + " " +
		m.stateMachine.Help()
}
