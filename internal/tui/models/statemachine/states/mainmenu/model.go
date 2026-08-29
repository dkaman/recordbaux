package mainmenu

import (
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/models/flist"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states/mainmenu/forms"
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
	keys          keyMap
	logger        *slog.Logger

	shelves   flist.Model
	playlists flist.Model
	creating  bool

	focus focusedView
}

func New(log *slog.Logger) MainMenuState {
	m := MainMenuState{
		keys:   defaultKeybinds(),
		logger: log,
	}

	m.logger = log.WithGroup("mainmenu")

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
	m.shelves = shelfList

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
	m.playlists = playlistList

	m.focus = shelvesView

	return m
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

		case key.Matches(msg, s.keys.New):
			if s.focus == shelvesView {
				s.logger.Debug("create shelf selected")
				return s, tcmds.ShowModalCmd(forms.NewCreateShelfForm())
			} else {
				s.logger.Debug("create playlist selected")

				var shelfIDs []uint
				for _, item := range s.shelves.Items() {
					if shM, ok := item.(tshelf.Model); ok {
						shelfIDs = append(shelfIDs, shM.PhysicalShelf().ID)
					}
				}

				return s, tcmds.ShowModalCmd(forms.NewCreatePlaylistForm(shelfIDs, s.logger))
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
