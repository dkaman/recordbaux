package statemachine

import (
	"errors"
	"log/slog"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/tui/models/overlay"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	lbs "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadedbin"
	lps "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadedplaylist"
	lss "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/loadedshelf"
	mms "github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/mainmenu"
)

var (
	StateNotFoundErr = errors.New("state not found in state map")
)


type Model struct {
	width, height int

	logger       *slog.Logger
	currentState states.State
}

func New(log *slog.Logger) (Model, error) {
	m := Model{}
	m.logger = log.WithGroup("statemachine")
	m.currentState = mms.New(m.logger)

	return m, nil
}

func (m Model) Init() tea.Cmd {
	m.logger.Debug("statemachine init",
		slog.String("currentState", m.currentState.Type().String()),
	)

	return m.currentState.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tcmds.TransitionToMainMenuMsg:
		newState := mms.New(m.logger)
		m.currentState = newState
		return m, newState.Init()

	case tcmds.TransitionToLoadedShelfMsg:
		newState, err := lss.New(m.logger, msg.ShelfID)
		if err != nil {
			m.logger.Error("error during transition to loadedshelfstate",
				slog.Any("err", err),
			)
			return m, nil
		}
		m.currentState = newState
		return m, newState.Init()

	case tcmds.TransitionToLoadedPlaylistMsg:
		newState, err := lps.New(m.logger, msg.PlaylistID)
		if err != nil {
			m.logger.Error("error during transition to loadedplayliststate",
				slog.Any("err", err),
			)
			return m, nil
		}
		m.currentState = newState
		return m, newState.Init()

	case tcmds.TransitionToLoadedBinMsg:
		newState, err := lbs.New(m.logger, msg.BinID)
		if err != nil {
			m.logger.Error("error during transition to loadedbinstate",
				slog.Any("err", err),
			)
			return m, nil
		}
		m.currentState = newState
		return m, newState.Init()
	}

	var stateCmds tea.Cmd
	m.currentState, stateCmds = util.UpdateModel(m.currentState, msg)

	return m, tea.Batch(stateCmds)
}

func (m Model) View() tea.View {
	currentState := m.currentState.View()

	viewportStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height)

	content := viewportStyle.Render(currentState.Content)

	return tea.NewView(content)
}

// overlay.teaFocusBlurer implementation

func (m Model) Focus() overlay.TeaFocusBlurer {
	return m
}

func (m Model) Blur() overlay.TeaFocusBlurer {
	return m
}

func (m Model) Help() string {
	return "statemachine: " + m.currentState.Help()
}

func (m Model) CurrentState() states.State {
	return m.currentState
}
