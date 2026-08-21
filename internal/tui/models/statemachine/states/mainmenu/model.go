package mainmenu

import (
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/flist"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	tplaylist "github.com/dkaman/recordbaux/internal/tui/models/playlist"
	tshelf "github.com/dkaman/recordbaux/internal/tui/models/shelf"
)

type focusedView int

const (
	shelvesView focusedView = iota
	playlistsView
)

type MainMenuState struct {
	width, height int
	keys   keyMap
	logger *slog.Logger

	shelves   flist.Model
	playlists flist.Model
	creating  bool

	createShelfForm *createShelfForm

	focus         focusedView
}

func New(log *slog.Logger) MainMenuState {
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
		tcmds.GetAllShelvesCmd(),
		tcmds.GetAllPlaylistsCmd(),
		tcmds.RefreshWindowSizeCmd(),
	)
}

func (s MainMenuState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	passthru := msg

	// form updates go first so they can accep enter keys etc. i need to see
	// if this is removable, i want all message handling first in all cases.
	if s.creating {
		var formCmd tea.Cmd

		// Intercept the resize message to enforce a 1/3 ratio
		if ws, ok := passthru.(tea.WindowSizeMsg); ok {
			targetWidth := max(ws.Width / 3 , 40)
			ws.Width = targetWidth - style.ModalStyle.GetHorizontalFrameSize()
			targetHeight := max(ws.Height / 3, 15)
			ws.Height = targetHeight - style.ModalStyle.GetVerticalFrameSize()
			passthru = ws
		}

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

		return s, tea.Batch(cmds...)

	case tea.KeyPressMsg:
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
				var shelfIDs []uint

				for _, item := range s.shelves.Items() {
					if shM, ok := item.(tshelf.Model); ok {
						shelfIDs = append(shelfIDs, shM.PhysicalShelf().ID)
					}
				}
				return s, func() tea.Msg{
					return tcmds.TransitionToCreatePlaylistMsg{
						ShelfIDs: shelfIDs,
					}
				}
			}

		case key.Matches(msg, s.keys.Select):
			if s.focus == shelvesView {
				if sel, ok := s.shelves.SelectedItem().(tshelf.Model); ok {
					return s, func() tea.Msg {
						return tcmds.TransitionToLoadedShelfMsg{
							ShelfID: sel.PhysicalShelf().ID,
						}
					}
				}
			} else {
				if sel, ok := s.playlists.SelectedItem().(tplaylist.Model); ok {
					return s, func() tea.Msg {
						return tcmds.TransitionToLoadedPlaylistMsg{
							PlaylistID: sel.PhysicalPlaylist().ID,
						}
					}
				}
			}
		}

		passthru = msg

	case tcmds.ShelvesLoadedMsg:
		if err := msg.Err; err != nil {
			s.logger.Error("error loading shelves",
				slog.Any("err", err),
			)
			return s, nil
		}

		shlvs := msg.Shelves
		items := make([]list.Item, len(shlvs))

		for i, sh := range shlvs {
			items[i] = tshelf.New(sh, s.logger)
		}

		s.shelves.SetItems(items)

		return s, nil

	case tcmds.ShelfSavedMsg:
		if err := msg.Err; err != nil {
			s.logger.Error("error saving shelf to database",
				slog.Any("err", err),
			)
		}

		return s, tcmds.GetAllShelvesCmd()

	case tcmds.PlaylistsLoadedMsg:
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

func (s MainMenuState) Type() states.StateType {
	return states.MainMenu
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

	// change this to save shelf and catch reply in update
	return tcmds.SaveShelfCmd(newShelf)
}
