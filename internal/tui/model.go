package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/config"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"

	discogs "github.com/dkaman/discogs-golang"
	css "github.com/dkaman/recordbaux/internal/tui/models/states/createshelf"
	lss "github.com/dkaman/recordbaux/internal/tui/models/states/loadedshelf"
	mms "github.com/dkaman/recordbaux/internal/tui/models/states/mainmenu"
	lcs "github.com/dkaman/recordbaux/internal/tui/models/states/loadcollection"
)

type Model struct {
	// global application config
	cfg *config.Config

	// tea models
	stateMachine statemachine.Model

	// styling/layout
	layout *layouts.TallLayout
}

func New(c *config.Config) Model {
	tallLayout := layouts.NewTallLayout()

	tallLayout.WithSection(layouts.StatusLine, "state: main menu")

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
			statemachine.MainMenu:    mms.New(tallLayout),
			statemachine.CreateShelf: css.New(tallLayout),
			statemachine.LoadedShelf: lss.New(tallLayout),
			statemachine.LoadCollection: lcs.New(tallLayout, discogsClient, discogsUsername),
		},

		// state machine's ref to the layout too
		tallLayout,
	)

	return Model{
		cfg:          c,
		stateMachine: sm,
		layout:       tallLayout,
	}
}

func (m Model) Init() tea.Cmd {
	return m.stateMachine.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	stateMachineModel, stateMachineCmds := m.stateMachine.Update(msg)
	if sm, ok := stateMachineModel.(statemachine.Model); ok {
		m.stateMachine = sm
	}

	cmds = append(cmds, stateMachineCmds)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	_ = m.stateMachine.View()
	return m.layout.Render()
}
