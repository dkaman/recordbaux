package createplaylist

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/v2/list"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"

	"github.com/dkaman/recordbaux/internal/db/playlist"
	"github.com/dkaman/recordbaux/internal/db/track"
	"github.com/dkaman/recordbaux/internal/services"
	"github.com/dkaman/recordbaux/internal/tui/handlers"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
	ttrack "github.com/dkaman/recordbaux/internal/tui/models/track"
)

type CreatePlaylistState struct {
	svcs     *services.AllServices
	logger   *slog.Logger
	keys     keyMap
	handlers *handlers.Registry

	list           list.Model
	namingPlaylist bool
	nameForm       *form
	playlistName   string

	width, height int
}

func New(svcs *services.AllServices, log *slog.Logger) CreatePlaylistState {
	logger := log.WithGroup("createplayliststate")

	delegate := trackDelegate{}
	trackList := list.New([]list.Item{}, delegate, 0, 0)
	trackList.Styles = style.DefaultListStyles()
	trackList.Title = "select tracks for new playlist"

	return CreatePlaylistState{
		svcs:   svcs,
		logger: logger,
		handlers: getHandlers(),
		list:   trackList,
		keys:   defaultKeybinds(),
		namingPlaylist: false,
	}
}

func (s CreatePlaylistState) Init() tea.Cmd {
	return nil
}

func (s CreatePlaylistState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if s.namingPlaylist {
		s.logger.Debug("begin naming playlist")
		fModel, formUpdatesCmds := s.nameForm.Update(msg)
		if f, ok := fModel.(*form); ok {
			s.nameForm = f
		}
		cmds = append(cmds, formUpdatesCmds)

		if s.nameForm.State == huh.StateCompleted {
			s.logger.Debug("end naming playlist")
			name := s.nameForm.Name()
			newPlaylist := &playlist.Entity{
				Name:   name,
				Tracks: make([]*track.Entity, 0),
			}

			for _, item := range s.list.Items() {
				if trackModel, ok := item.(ttrack.Model); ok && trackModel.Selected {
					newPlaylist.Tracks = append(newPlaylist.Tracks, trackModel.PhysicalTrack())
				}
			}
			cmds = append(cmds, s.svcs.SavePlaylistCmd(newPlaylist))

			items := s.list.Items()
			for i, item := range items {
				if trackModel, ok := item.(ttrack.Model); ok && trackModel.Selected {
					trackModel.Selected = false
					items[i] = trackModel
				}
			}

			s.list.SetItems(items)
			s.namingPlaylist = false
			s.nameForm = newNameForm()

			return s, tcmds.Transition(
				states.MainMenu, cmds, nil,
			)
		}

		return s, tea.Batch(cmds...)
	}

	if handler, ok := s.handlers.GetHandler(msg); ok {
		model, cmd, passthruMsg := handler(s, msg)
		if passthruMsg == nil {
			return model, cmd
		}
		s = model.(CreatePlaylistState)
		msg = passthruMsg
		cmds = append(cmds, cmd)
	}

	var listCmd tea.Cmd
	s.list, listCmd = s.list.Update(msg)
	cmds = append(cmds, listCmd)
	return s, tea.Batch(cmds...)
}

func (s CreatePlaylistState) View() string {
	return s.renderModel()
}

func (s CreatePlaylistState) Title() string {
	return "create playlist"
}

func (s CreatePlaylistState) State() states.StateType {
	return states.CreatePlaylist
}

func (s CreatePlaylistState) Help() string {
	return "create a playlist"
}
