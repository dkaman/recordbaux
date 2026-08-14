package mainmenu

import (
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/models/flist"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	tshelf "github.com/dkaman/recordbaux/internal/tui/models/shelf"
	tplaylist "github.com/dkaman/recordbaux/internal/tui/models/playlist"
)

type focusedView int

const (
	shelvesView focusedView = iota
	playlistsView
)

type MainMenuState struct {
	// meta stuff
	svcs   *services.AllServices
	keys   keyMap
	logger *slog.Logger

	// main menu stuff
	shelves   flist.Model
	playlists flist.Model
	creating  bool

	// create shelf stuff
	createShelfForm *createShelfForm

	focus         focusedView
	width, height int
}

func New(svcs *services.AllServices, log *slog.Logger) MainMenuState {
	log = log.WithGroup("mainmenu")

	// Shelves List
	shelfDelegate := newShelfDelegate(shelfDelegateStyles{
		ItemStyle:            style.TextStyle,
		ItemStyleBlurred:     style.TextStyleDimmed,
		SelectedStyle:        style.ActiveTextStyle,
		SelectedStyleBlurred: style.ActiveLabelStyleDimmed,
	})

	shelfList := flist.New([]list.Item{}, shelfDelegate)
	shelfList.Title = "shelves"
	shelfList.Styles = style.DefaultListStyles()

	// Playlists List
	playlistDelegate := newPlaylistDelegate(playlistDelegateStyles{
		ItemStyle:            style.TextStyle,
		ItemStyleBlurred:     style.TextStyleDimmed,
		SelectedStyle:        style.ActiveTextStyle,
		SelectedStyleBlurred: style.ActiveLabelStyleDimmed,
	})

	playlistList := flist.New([]list.Item{}, playlistDelegate)
	playlistList.Title = "playlists"
	playlistList.Styles = style.DefaultListStyles()

	return MainMenuState{
		svcs:   svcs,
		keys:   defaultKeybinds(),
		logger: log,

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
	var passthru tea.Msg

	passthru = msg

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height

		halfWidth := s.width / 2
		otherWidth := s.width - halfWidth

		var shelfUpdateCmds tea.Cmd

		shelfSizeMsg := tea.WindowSizeMsg{
			Width:  halfWidth,
			Height: s.height,
		}

		s.shelves, shelfUpdateCmds = util.UpdateModel(s.shelves, shelfSizeMsg)
		cmds = append(cmds, shelfUpdateCmds)

		var playlistUpdateCmds tea.Cmd

		playlistSizeMsg := tea.WindowSizeMsg{
			Width:  otherWidth,
			Height: s.height,
		}

		s.playlists, playlistUpdateCmds = util.UpdateModel(s.playlists, playlistSizeMsg)
		cmds = append(cmds, playlistUpdateCmds)

		s.logger.Debug("window sizes for child models",
			slog.Any("shelf width", s.shelves.Width()),
			slog.Any("shelf height", s.shelves.Height()),
			slog.Any("playlist width", s.playlists.Width()),
			slog.Any("playlist height", s.playlists.Height()),
		)

		return s, tea.Batch(cmds...)

	case tea.KeyPressMsg:
		if s.creating {
			return s, nil
		}

		switch {
		case key.Matches(msg, s.keys.Quit):
			return s, nil

		case key.Matches(msg, s.keys.SwitchFocus):
			if s.focus == shelvesView {
				s.focus = playlistsView
				s.playlists = s.playlists.Focus()
				s.shelves = s.shelves.Blur()
			} else {
				s.focus = shelvesView
				s.shelves = s.shelves.Focus()
				s.playlists = s.playlists.Blur()
			}

			return s, nil

		case key.Matches(msg, s.keys.NewShelf):
			if s.focus == shelvesView {
				s.logger.Debug("create shelf selected")
				s.creating = true
				s.shelves = s.shelves.Blur()
				s.playlists = s.playlists.Blur()
				s.createShelfForm = newCreateShelfForm()
				return s, s.createShelfForm.Init()
			} else {
				s.logger.Debug("create playlist selected")
				return s, tcmds.Transition(
					states.CreatePlaylist,
					nil,
					[]tea.Cmd{s.svcs.GetAllTracksCmd()},
				)
			}

		case key.Matches(msg, s.keys.Select):
			if s.focus == shelvesView {
				if sel, ok := s.shelves.SelectedItem().(tshelf.Model); ok {
					return s, tcmds.Transition(
						states.LoadedShelf,
						nil,
						[]tea.Cmd{tshelf.WithPhysicalShelf(sel.PhysicalShelf())},
					)
				}
			} else {
				if sel, ok := s.playlists.SelectedItem().(tplaylist.Model); ok {
					return s, tcmds.Transition(
						states.LoadedPlaylist,
						nil,
						[]tea.Cmd{tplaylist.WithPhysicalPlaylist(sel.PhysicalPlaylist())},
					)
				}
			}
		}

		// pass key message through if we aren't handling it
		passthru = msg

	case services.ShelvesLoadedMsg:
		shlvs := msg.Shelves
		items := make([]list.Item, len(shlvs))

		for i, sh := range shlvs {
			items[i] = tshelf.New(sh, s.logger)
		}

		s.shelves.SetItems(items)

		return s, nil

	case services.PlaylistsLoadedMsg:
		playlists := msg.Playlists
		playlistItems := make([]list.Item, len(playlists))

		for i, p := range playlists {
			play := tplaylist.New()
			play.SetEntity(p)
			playlistItems[i] = play
		}

		s.playlists.SetItems(playlistItems)

		return s, nil
	}

	// form updates go first so they can accep enter keys etc.
	if s.creating {
		var formCmd tea.Cmd
		s.createShelfForm, formCmd = util.UpdateModel(s.createShelfForm, passthru)
		cmds = append(cmds, formCmd)

		if s.createShelfForm.Form.State == huh.StateCompleted {
			s.creating = false
			if s.focus == shelvesView {
				s.shelves = s.shelves.Focus()
			} else {
				s.playlists = s.playlists.Focus()
			}

			cmds = append(cmds, formCmd, s.handleShelfCreation())
		}

		return s, tea.Batch(cmds...)
	}

	var updateCmd tea.Cmd

	if s.focus == shelvesView {
		s.shelves, updateCmd = util.UpdateModel(s.shelves, passthru)
	} else {
		s.playlists, updateCmd = util.UpdateModel(s.playlists, passthru)
	}

	cmds = append(cmds, updateCmd)

	return s, tea.Batch(cmds...)
}

func (s MainMenuState) View() tea.View {
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
