package forms

import (
	"io"
	"fmt"
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/db/track"
	"github.com/dkaman/recordbaux/internal/tui/style"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	ttrack "github.com/dkaman/recordbaux/internal/tui/models/track"
)

type createPlaylistFormKeyMap struct {
	Select key.Binding
	Create key.Binding
}

func defaultCreatePlaylistFormKeybinds() createPlaylistFormKeyMap {
	return createPlaylistFormKeyMap{
		Select: key.NewBinding(
			key.WithKeys("space"), // spacebar
			key.WithHelp("space", "toggle select"),
		),
		Create: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "create playlist"),
		),
	}
}

func (k createPlaylistFormKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Select,
		k.Create,
	}
}

func (k createPlaylistFormKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Create},
	}
}

// createPlaylistForm implements the workflow to collect tracks and a name to
// create a new playlist with. it allows the user to scroll a list of tracks,
// hit space to mark it, and then finally hit enter to create.

type CreatePlaylistForm struct {
	width, height int
	active        bool
	keys          createPlaylistFormKeyMap
	logger        *slog.Logger

	shelves        []uint
	namingPlaylist bool
	nameForm       *huh.Form

	trackList list.Model
}

func NewCreatePlaylistForm(shelfIDs []uint, l *slog.Logger) CreatePlaylistForm {
	f := CreatePlaylistForm{
		shelves: shelfIDs,
		logger:  l,
		active:  true,
	}

	// create naming form
	f.nameForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key(keyPlaylistName).
				Title("playlist name").
				Validate(huh.ValidateNotEmpty()),
		),
	).WithTheme(huh.ThemeFunc(style.DefaultFormStyles))

	f.keys = defaultCreatePlaylistFormKeybinds()

	return f
}

// tea.Model implementation

func (f CreatePlaylistForm) Init() tea.Cmd {
	return tcmds.GetAllTracksFromShelvesCmd(f.shelves)
}

func (f CreatePlaylistForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	passthru := msg

	if f.namingPlaylist {
		fModel, formUpdatesCmds := f.nameForm.Update(msg)

		if frm, ok := fModel.(*huh.Form); ok {
			f.nameForm = frm
		}

		cmds = append(cmds, formUpdatesCmds)

		if f.nameForm.State == huh.StateCompleted {
			var tracks []*track.Entity

			name := f.nameForm.GetString(keyPlaylistName)

			for _, item := range f.trackList.Items() {
				if trackModel, ok := item.(ttrack.Model); ok && trackModel.Selected {
					tracks = append(tracks, trackModel.PhysicalTrack())
				}
			}

			cmds = append(cmds, tcmds.NewPlaylistCmd(name, tracks))

			f.namingPlaylist = false
			f.active = false
		}

		return f, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width, f.height = msg.Width, msg.Height

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, f.keys.Select):
			if i, ok := f.trackList.SelectedItem().(ttrack.Model); ok {
				i.Selected = !i.Selected
				cmd := f.trackList.SetItem(f.trackList.Index(), i)
				return f, cmd
			}

			return f, nil

		case key.Matches(msg, f.keys.Create):
			var selectedCount int

			for _, item := range f.trackList.Items() {
				if trackModel, ok := item.(ttrack.Model); ok && trackModel.Selected {
					selectedCount++
				}
			}

			if selectedCount > 0 {
				f.namingPlaylist = true
				cmds = append(cmds, f.nameForm.Init())
			}

			return f, tea.Batch(cmds...)
		}

	case tcmds.ShelvesAllTracksMsg:
		if msg.Err != nil {
			f.logger.Error("error getting all tracks", slog.Any("err", msg.Err))
			return f, nil
		}

		items := make([]list.Item, len(msg.Tracks))

		for i, t := range msg.Tracks {
			items[i] = ttrack.New(t)
		}

		delegate := trackDelegate{}
		trackList := list.New(items, delegate, 0, 0)
		trackList.Styles = style.DefaultListStyles()
		trackList.Title = "select tracks for new playlist"

		f.trackList = trackList

		return f, nil
	}

	var listCmd tea.Cmd
	f.trackList, listCmd = f.trackList.Update(passthru)
	cmds = append(cmds, listCmd)

	return f, tea.Batch(cmds...)
}

func (f CreatePlaylistForm) View() tea.View {
	if f.namingPlaylist {
		return tea.NewView(f.nameForm.View())
	}

	content := f.trackList.View()

	return tea.NewView(content)
}

// modal.TeaActiveCloser implementation

func (f CreatePlaylistForm) Active() bool {
	return f.active
}

func (f CreatePlaylistForm) Close() tea.Cmd {
	return func() tea.Msg {
		return tcmds.TransitionToMainMenuMsg{}
	}
}

type trackDelegate struct{}

func (d trackDelegate) Height() int { return 2 }

func (d trackDelegate) Spacing() int { return 1 }

func (d trackDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d trackDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	boldStyle := lipgloss.NewStyle().
		Bold(true)

	dimStyle := lipgloss.NewStyle().
		Foreground(style.LightGrey)

	trackM, ok := listItem.(ttrack.Model)
	if !ok {
		return
	}

	selectedMarker := "[ ] "
	if trackM.Selected {
		selectedMarker = "[" + boldStyle.Render("*") + "] "
	}

	pTrack := trackM.PhysicalTrack()

	titleLine := fmt.Sprintf("%s - %s", boldStyle.Render(pTrack.Title), pTrack.Duration)
	descLine := fmt.Sprintf("key: '%s' bpm: %d", pTrack.Key, pTrack.BPM)
	prefix := "  "
	if m.Index() == index {
		prefix = "> "
	}

	lines := lipgloss.JoinVertical(lipgloss.Left, titleLine, dimStyle.Render(descLine))
	display := lipgloss.JoinHorizontal(lipgloss.Top, prefix, selectedMarker, lines)

	sty := style.TextStyle
	if m.Index() == index {
		sty = style.ActiveTextStyle
	}

	fmt.Fprint(w, sty.Render(display))
}
