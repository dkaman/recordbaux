package loadcollection

import (
	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
)

func getHandlers() *handlers.Registry {
	r := handlers.NewRegistry()

	handlers.Register(r, handleTeaWindowSizeMsg)
	handlers.Register(r, handleLoadShelfMsg)

	return r
}

func handleTeaWindowSizeMsg(s LoadCollectionState, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.width, s.height = msg.Width, msg.Height
	return s, nil, nil
}

func handleLoadShelfMsg(s LoadCollectionState, msg shelf.LoadShelfMsg) (tea.Model, tea.Cmd, tea.Msg) {
	s.shelf = msg.Phy
	return s, nil, nil
}
