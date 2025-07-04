package tui

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	kfmt "github.com/dkaman/recordbaux/internal/tui/key"
)

var (
	LoggerIsNilErr = errors.New("supplied slog logger is nil")
)

type Model struct {
	// global application config/state
	cfg          *config.Config
	shelfService *services.ShelfService
	keys         keyMap
	logger       *slog.Logger

	// tea models
	stateMachine statemachine.Model

	// styling/layout
	helpVisible bool
	layout      *layout.Div
}

func New(c *config.Config, log *slog.Logger, d db.Repository[*shelf.Entity]) (Model, error) {
	var m Model

	if log == nil {
		return m, LoggerIsNilErr
	}

	h := help.New()
	h.Styles = style.DefaultHelpStyles()

	s := services.NewShelfService(d)

	l, err := newTUILayout()
	if err != nil {
		return m, fmt.Errorf("error creating layout: %w", err)
	}

	vp := l.Find("viewport")

	sm, err := statemachine.New(s, c, vp, log)
	if err != nil {
		return m, fmt.Errorf("error creating state machine: %w", err)
	}

	m = Model{
		cfg:          c,
		shelfService: s,
		keys:         defaultKeybinds(),
		helpVisible:  false,
		layout:       l,
		stateMachine: sm,
		logger:       log.WithGroup("root"),
	}

	err = addTopBarText(l, "recordbaux - organize your record collection")
	if err != nil {
		return m, fmt.Errorf("error adding top bar text: %w", err)
	}

	err = addStatusBarText(m.layout, fmt.Sprintf("current state: %s", m.stateMachine.CurrentStateType()))
	if err != nil {
		return m, fmt.Errorf("error adding status bar text: %w", err)
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

	// we're going to handle all db messages in root
	case tcmds.ShelfSavedMsg:
		if msg.Err != nil {
			m.logger.Error("error saving shelf",
				slog.String("error", msg.Err.Error()),
			)
		}
		return m, tea.Batch(cmds...)

	case tcmds.ShelfLoadedMsg:
		if err := msg.Err; err != nil {
			m.logger.Error("error loading shelf",
				slog.String("error", err.Error()),
			)
			return m, tea.Batch(cmds...)
		}

		m.shelfService.CurrentShelf = msg.Shelf

		return m, tea.Batch(cmds...)

	case tcmds.ShelvesLoadedMsg:
		if err := msg.Err; err != nil {
			m.logger.Error("error loading all shelves",
				slog.String("error", err.Error()),
			)
			return m, tea.Batch(cmds...)
		}

		m.shelfService.AllShelves = msg.Shelves

		return m, tea.Batch(cmds...)

	case tcmds.ShelfDeletedMsg:
		if err := msg.Err; err != nil {
			m.logger.Error("error deleting shelf",
				slog.String("error", err.Error()),
			)
		}

		return m, tea.Batch(cmds...)
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
		kfmt.FmtKeymap(m.keys.ShortHelp()) + " " +
		m.stateMachine.Help()
}
