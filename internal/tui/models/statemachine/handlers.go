package statemachine

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

func getHandlers() *handlers.Registry {
	r := handlers.NewRegistry()

	handlers.Register(r, handleTeaWindowSizeMsg)
	handlers.Register(r, handleStateTransitionMsg)

	return r
}

func handleTeaWindowSizeMsg(m Model, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd, tea.Msg) {
	m.width, m.height = msg.Width, msg.Height
	return m, nil, msg
}

func handleStateTransitionMsg(m Model, msg tcmds.StateTransitionMsg) (tea.Model, tea.Cmd, tea.Msg) {
	next := msg.Transition.Next

	nextState := m.allStates[next]

	sizeMsg := tea.WindowSizeMsg{
		Width: m.width,
		Height: m.height,
	}

	resizedNextState, sizeUpdateCmd := util.UpdateModel(nextState, sizeMsg)

	m.logger.Info("state transition",
		slog.String("from", m.currentStateType.String()),
		slog.String("to", next.String()),
	)

	m.allStates[m.currentStateType] = m.currentState
	m.currentState = resizedNextState
	m.currentStateType = next
	m.allStates[next] = resizedNextState

	// new state will be initialized and then post-transition commands will
	// run
	return m, tea.Sequence(
		m.currentState.Init(),
		sizeUpdateCmd,
		tea.Batch(msg.Transition.PostCmds...),
	), nil
}
