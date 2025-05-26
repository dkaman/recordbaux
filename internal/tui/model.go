package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

	discogs "github.com/dkaman/discogs-golang"
	css "github.com/dkaman/recordbaux/internal/tui/models/states/createshelf"
	lcs "github.com/dkaman/recordbaux/internal/tui/models/states/loadcollection"
	lbs "github.com/dkaman/recordbaux/internal/tui/models/states/loadedbin"
	lss "github.com/dkaman/recordbaux/internal/tui/models/states/loadedshelf"
	mms "github.com/dkaman/recordbaux/internal/tui/models/states/mainmenu"
	sss "github.com/dkaman/recordbaux/internal/tui/models/states/selectshelf"
)

type Model struct {
	// global application config/state
	cfg  *config.Config
	app  *app.App
	keys keyMap
	help help.Model

	// tea models
	stateMachine statemachine.Model

	// styling/layout
	helpVisible bool
	layout      *layout.Node
	topBar      string
	viewPort    string
	statusBar   string
	helpBar     string
	overlay     string
}

func New(c *config.Config, l *layout.Node) Model {
	var err error

	h := help.New()
	h.Styles = style.DefaultHelpStyles()

	m := Model{
		cfg:  c,
		app:  app.NewApp(),
		keys: defaultKeybinds(),
		help: h,
	}

	initialState := statemachine.MainMenu

	discogsUsername, _ := c.String("shelf.discogs.username")
	discogsAPIKey, _ := c.String("shelf.discogs.key")
	discogsClient, err := discogs.New(
		discogs.WithToken(discogsAPIKey),
	)
	if err != nil {
		log.Printf("error in discogs client creation %w", err)
	}

	sm, _ := statemachine.New(
		initialState,
		map[statemachine.StateType]statemachine.State{
			statemachine.MainMenu:       mms.New(m.app),
			statemachine.CreateShelf:    css.New(m.app),
			statemachine.LoadedShelf:    lss.New(m.app),
			statemachine.LoadCollection: lcs.New(m.app, discogsClient, discogsUsername),
			statemachine.LoadedBin:      lbs.New(m.app),
			statemachine.SelectShelf:    sss.New(m.app),
		},
	)

	m.stateMachine = sm
	m.layout, _ = newTUILayout(l)

	return m
}

func (m Model) Init() tea.Cmd {
	return m.stateMachine.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.ToggleHelp):
			m.helpVisible = !m.helpVisible
		}
	}

	m.helpBar = fmt.Sprintf("global: %s state: %s", m.help.View(m.keys), m.stateMachine.Help())

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
