package mainmenu

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/list"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"
)

type focusedView int

const (
	shelvesView focusedView = iota
	playlistsView
)

type MainMenuState struct {
	svcs     *services.AllServices
	keys     keyMap
	logger   *slog.Logger
	handlers *handlers.Registry

	shelves   list.Model
	playlists list.Model

	focus         focusedView
	width, height int
}

func New(svcs *services.AllServices, log *slog.Logger) MainMenuState {
	log = log.WithGroup("mainmenu")

	// Shelves List
	shelfDelegate := shelfDelegate{}
	shelfList := list.New([]list.Item{}, shelfDelegate, 0, 0)
	shelfList.Title = "shelves"
	shelfList.Styles = style.DefaultListStyles()

	// Playlists List
	playlistDelegate := playlistDelegate{}
	playlistList := list.New([]list.Item{}, playlistDelegate, 0, 0)
	playlistList.Title = "playlists"
	playlistList.Styles = style.DefaultListStyles()

	return MainMenuState{
		svcs:      svcs,
		keys:      defaultKeybinds(),
		logger:    log,
		handlers:  getHandlers(),
		shelves:   shelfList,
		playlists: playlistList,
		focus:     shelvesView,
	}
}

func (s MainMenuState) Init() tea.Cmd {
	s.logger.Debug("mainmenu state init",
		slog.Int("numShelves", len(s.shelves.Items())),
	)

	return tea.Sequence(
		s.svcs.GetAllShelvesCmd(),
		s.svcs.GetAllPlaylistsCmd(),
	)
}

func (s MainMenuState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if handler, ok := s.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(s, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		s = model.(MainMenuState)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	var updateCmd tea.Cmd
	if s.focus == shelvesView {
		s.shelves, updateCmd = s.shelves.Update(msg)
	} else {
		s.playlists, updateCmd = s.playlists.Update(msg)
	}

	return s, tea.Batch(updateCmd)
}

func (s MainMenuState) View() string {
	return s.renderModel()
}

func (s MainMenuState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}
