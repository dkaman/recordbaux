package mainmenu

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/list"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/shelf"
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
	// meta stuff
	svcs     *services.AllServices
	keys     keyMap
	logger   *slog.Logger
	handlers *handlers.Registry

	// main menu stuff
	shelves   list.Model
	playlists list.Model
	creating  bool

	// create shelf stuff
	createShelfForm *createShelfForm

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
		svcs:     svcs,
		keys:     defaultKeybinds(),
		logger:   log,
		handlers: getHandlers(),

		shelves:   shelfList,
		playlists: playlistList,

		focus: shelvesView,
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

	// form updates go first so they can accep enter keys etc.
	if s.creating {
		var formCmd tea.Cmd
		s.createShelfForm, formCmd = util.UpdateModel(s.createShelfForm, msg)
		cmds = append(cmds, formCmd)
		if s.createShelfForm.Form.State == huh.StateCompleted {
			s.creating = false
			cmds = append(cmds, formCmd, s.handleShelfCreation())
		}
		return s, tea.Batch(cmds...)
	}

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

	cmds = append(cmds, updateCmd)

	return s, tea.Batch(cmds...)
}

func (s MainMenuState) View() string {
	return s.renderModel()
}

func (s MainMenuState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}

func (s MainMenuState) handleShelfCreation() tea.Cmd {
	f := s.createShelfForm

	x := f.DimX()
	y := f.DimY()
	size := f.BinSize()

	var newShelf *shelf.Entity

	if f.Shape() == Rect {
		newShelf, _ = shelf.New(f.Name(), size,
			shelf.WithShapeRect(x, y, size, bin.SortAlphaByArtist),
		)
	} else {
		newShelf, _ = shelf.New(f.Name(), size)
	}

	s.logger.Debug("new shelf", slog.Any("shelf", newShelf))

	return tea.Sequence(
		s.svcs.SaveShelfCmd(newShelf),
		s.svcs.GetAllShelvesCmd(),
	)
}
