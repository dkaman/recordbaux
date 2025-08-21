package mainmenu

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"

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

	keys   keyMap
	layout *layout.Div

	logger    *slog.Logger
	shelves   list.Model
	playlists list.Model
	focus     focusedView
}

type refreshMsg struct{}

func New(s *services.ShelfService, t *services.TrackService, p *services.PlaylistService, l *layout.Div, log *slog.Logger) MainMenuState {
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

	lay, _ := newMainMenuLayout(l, shelfList, playlistList, shelvesView)

	return MainMenuState{
		shelfService:    s,
		trackService:    t,
		playlistService: p,
		keys:            defaultKeybinds(),
		layout:          lay,
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
		s.refresh(),
	)
}

func (s MainMenuState) refresh() tea.Cmd {
	return func() tea.Msg {
		return refreshMsg{}
	}
}

func (s MainMenuState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case refreshMsg:
		s.logger.Debug("refreshing shelves from service")

		shlvs := s.shelfService.AllShelves
		items := make([]list.Item, len(shlvs))

		for i, sh := range shlvs {
			items[i] = shelf.New(sh, s.logger)
		}
		s.shelves.SetItems(items)

		playlists := s.playlistService.AllPlaylists
		playlistItems := make([]list.Item, len(playlists))
		for i, p := range playlists {
			playlistItems[i] = tplaylist.New(p)
		}
		s.playlists.SetItems(playlistItems)

		s.layout, _ = newMainMenuLayout(s.layout, s.shelves, s.playlists, s.focus)
		return s, nil

	case tea.KeyMsg:
		s.logger.Debug("key pressed",
			slog.String("key", msg.Type.String()),
		)

		switch {
		case key.Matches(msg, s.keys.SwitchFocus):
			if s.focus == shelvesView {
				s.focus = playlistsView
			} else {
				s.focus = shelvesView
			}
			s.layout, _ = newMainMenuLayout(s.layout, s.shelves, s.playlists, s.focus)
			return s, nil

		case key.Matches(msg, s.keys.NewShelf):
			if s.focus == shelvesView {
				s.logger.Debug("create shelf selected")
				s.nextState = states.CreateShelf
				return s, tea.Batch(cmds...)
			} else {
				s.logger.Debug("create playlist selected")
				cmds = append(cmds, tcmds.GetAllTracksCmd(s.trackService.Tracks, s.logger))
				s.nextState = states.CreatePlaylist
				return s, tea.Batch(cmds...)
			}

		case key.Matches(msg, s.keys.Select):
			if s.focus == shelvesView {
				if sel, ok := s.shelves.SelectedItem().(shelf.Model); ok {
					s.logger.Debug("selected shelf", slog.Any("id", sel.ID()))
					s.nextState = states.LoadedShelf
					s.shelfService.CurrentShelf = sel.PhysicalShelf()
					return s, tea.Batch(cmds...)
				}
				s.logger.Warn("somehow a shelf was selected that wasn't a shelf.Model")
				return s, tea.Batch(cmds...)
			} else {
				s.logger.Debug("selected playlist")
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

	s.layout, _ = newMainMenuLayout(s.layout, s.shelves, s.playlists, s.focus)
	return s, tea.Batch(cmds...)
}

func (s MainMenuState) View() string {
	return s.layout.Render()
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
