package mainmenu

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/playlist"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	keyFmt "github.com/dkaman/recordbaux/internal/tui/key"
	tplaylist "github.com/dkaman/recordbaux/internal/tui/models/playlist"
)

type focusedView int

const (
	shelvesView focusedView = iota
	playlistsView
)

type MainMenuState struct {
	shelfService    *services.ShelfService
	trackService    *services.TrackService
	playlistService *services.PlaylistService
	nextState       states.StateType

	keys keyMap

	logger        *slog.Logger
	shelves       list.Model
	playlists     list.Model
	focus         focusedView
	width, height int
}

func New(s *services.ShelfService, t *services.TrackService, p *services.PlaylistService, log *slog.Logger) MainMenuState {
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
		shelfService:    s,
		trackService:    t,
		playlistService: p,
		keys:            defaultKeybinds(),
		shelves:         shelfList,
		playlists:       playlistList,
		logger:          log,
		focus:           shelvesView,
		nextState:       states.Undefined,
	}
}

func (s MainMenuState) Init() tea.Cmd {
	s.logger.Debug("mainmenu state init",
		slog.Int("numShelves", len(s.shelves.Items())),
	)

	return tea.Sequence(
		tcmds.GetAllShelvesCmd(s.shelfService.Shelves, s.logger),
		tcmds.GetAllPlaylistsCmd(s.playlistService.Playlists, s.logger),
	)
}

func (s MainMenuState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height
		return s, nil

	case tcmds.ShelvesLoadedMsg:
		s.logger.Debug("refreshing shelves from service")

		shlvs := msg.Shelves
		items := make([]list.Item, len(shlvs))

		for i, sh := range shlvs {
			items[i] = shelf.New(sh, s.logger)
		}

		s.shelves.SetItems(items)

		return s, nil

	case tcmds.PlaylistsLoadedMsg:
		s.logger.Debug("refreshing playlists from service")

		playlists := msg.Playlists
		playlistItems := make([]list.Item, len(playlists))

		for i, p := range playlists {
			playlistItems[i] = tplaylist.New(p)
		}

		s.playlists.SetItems(playlistItems)

		return s, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.SwitchFocus):
			if s.focus == shelvesView {
				s.focus = playlistsView
			} else {
				s.focus = shelvesView
			}

			return s, nil

		case key.Matches(msg, s.keys.NewShelf):
			if s.focus == shelvesView {
				s.logger.Debug("create shelf selected")
				return s, tcmds.WithNextState(states.CreateShelf, nil, nil)
			} else {
				s.logger.Debug("create playlist selected")
				return s, tcmds.WithNextState(
					states.CreatePlaylist,
					nil,
					[]tea.Cmd{tcmds.GetAllTracksCmd(s.trackService.Tracks, s.logger)},
				)
			}

		case key.Matches(msg, s.keys.Select):
			if s.focus == shelvesView {
				if sel, ok := s.shelves.SelectedItem().(shelf.Model); ok {
					return s, tcmds.WithNextState(
						states.LoadedShelf,
						nil,
						[]tea.Cmd{shelf.WithPhysicalShelf(sel.PhysicalShelf())},
					)
				}
			} else {
				if sel, ok := s.playlists.SelectedItem().(tplaylist.Model); ok {
					return s, tcmds.WithNextState(
						states.LoadedPlaylist,
						nil,
						[]tea.Cmd{playlist.WithPhysicalPlaylist(sel.PhysicalPlaylist())},
					)
				}
			}
		}
	}

	var updateCmds tea.Cmd
	if s.focus == shelvesView {
		s.shelves, updateCmds = s.shelves.Update(msg)
	} else {
		s.playlists, updateCmds = s.playlists.Update(msg)
	}

	cmds = append(cmds, updateCmds)

	return s, tea.Batch(cmds...)
}

func (s MainMenuState) View() string {
	return s.renderModel()
}

func (s MainMenuState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}

	return states.Undefined, false
}

func (s MainMenuState) Transition() states.State {
	s.nextState = states.Undefined
	return s
}

func (s MainMenuState) Help() string {
	return keyFmt.FmtKeymap(s.keys.ShortHelp())
}
