package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"golang.org/x/term"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"

	discogs "github.com/dkaman/discogs-golang"
	teaCmds "github.com/dkaman/recordbaux/internal/tui/cmds"
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
	layout      *layouts.TwoBarViewportLayout
	topBar      string
	viewPort    string
	statusBar   string
	helpBar     string
	overlay     string
}

func New(c *config.Config) Model {
	initialState := statemachine.MainMenu

	m := Model{
		cfg:  c,
		app:  app.NewApp(),
		keys: defaultKeybinds(),
		help: help.New(),

		topBar:      "recordbaux - organize your vinyl record collection",
		viewPort:    "welcome to recordbaux",
		statusBar:   fmt.Sprintf("current state: %s", initialState),
		helpVisible: false,
		helpBar:     "",
	}

	m.layout = layouts.NewTwoBarViewportLayout()

	discogsAPIKey, _ := c.String("shelf.discogs.key")

	discogsClient, err := discogs.New(
		discogs.WithToken(discogsAPIKey),
	)
	if err != nil {
		log.Printf("error in discogs client creation %w", err)
	}

	discogsUsername, _ := c.String("shelf.discogs.username")

	sm, _ := statemachine.New(
		// our initial state is main menu
		statemachine.MainMenu,

		// pass the layout to all the states so they can add if they want
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

	return m
}

func (m Model) Init() tea.Cmd {
	return m.stateMachine.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.helpBar = fmt.Sprintf("global %sstate: %s", m.help.View(m.keys), m.stateMachine.Help())

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.ToggleHelp):
			m.helpVisible = !m.helpVisible
		}

	case teaCmds.LayoutUpdateMsg:
		sec := msg.Section
		cont := msg.Content

		switch sec {
		case layouts.TopBar:
			m.topBar = cont
		case layouts.Viewport:
			m.viewPort = cont
		case layouts.StatusBar:
			m.statusBar = cont
		case layouts.Overlay:
			m.overlay = cont
		}

		// if we're just updating the layout, dont pass the message on
		return m, nil
	}

	stateMachineModel, stateMachineCmds := m.stateMachine.Update(msg)
	if sm, ok := stateMachineModel.(statemachine.Model); ok {
		m.stateMachine = sm
	}

	cmds = append(cmds, stateMachineCmds)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	totalW, totalH, _ := term.GetSize(int(os.Stdout.Fd()))

	m.layout.SetSize(totalW, totalH)

	h := ""
	if m.helpVisible {
		h = m.helpBar
	}

	return m.layout.Render(m.topBar, m.viewPort, m.statusBar, m.overlay, h)
}

func (m Model) Help() string {
	return ""
}
