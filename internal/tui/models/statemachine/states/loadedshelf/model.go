package loadedshelf

import (
	"log/slog"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/util"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type LoadedShelfState struct {
	keys   keyMap
	logger *slog.Logger

	// loaded shelf
	shelfID     uint
	shelf       shelf.Model
	selectedBin int

	loading            bool
	fetching           bool
	loadCollectionForm *loadCollectionForm

	// Fields for fetching progress UI and state
	spin          spinner.Model
	prog          progress.Model
	releases      []*record.Entity
	currentIndex  int
	totalReleases int
	pct           float64
	currentTitle  string

	width, height int
}

type loadNextMsg struct{}

func (s LoadedShelfState) loadNextRecord() tea.Cmd {
	return func() tea.Msg {
		return loadNextMsg{}
	}
}

// New constructs a LoadedShelfState
func New(log *slog.Logger, shelfID uint) (LoadedShelfState, error) {
	s := LoadedShelfState{
		shelfID:  shelfID,
		keys:     defaultKeybinds(),
		fetching: false,
		loading:  false,
	}

	s.logger = log.WithGroup(states.LoadedShelf.String())

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(style.LightMagenta)
	s.spin = sp

	prg := progress.New(progress.WithDefaultBlend())
	s.prog = prg

	return s, nil
}

func (s LoadedShelfState) Init() tea.Cmd {
	s.logger.Debug("loadedshelf state init",
		slog.Int("shelfID", int(s.shelfID)),
	)

	return tea.Sequence(
		tcmds.GetShelfCmd(s.shelfID),
		tcmds.RefreshWindowSizeCmd(),
	)
}

func (s LoadedShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	passthru := msg

	if s.loading {
		var formCmd tea.Cmd

		if ws, ok := passthru.(tea.WindowSizeMsg); ok {
			targetWidth := max(ws.Width / 3, 40)
			ws.Width = targetWidth - style.ModalStyle.GetHorizontalFrameSize()
			targetHeight := max(ws.Height / 3, 15)
			ws.Height = targetHeight - style.ModalStyle.GetVerticalFrameSize()
			passthru = ws
		}

		s.loadCollectionForm, formCmd = util.UpdateModel(s.loadCollectionForm, passthru)
		cmds = append(cmds, formCmd)

		if s.loadCollectionForm.Form.State == huh.StateCompleted {
			folder := s.loadCollectionForm.Folder()
			s.logger.Debug("folder selected, starting fetch", slog.String("folder", folder))
			s.loading = false
			s.fetching = true
			cmds = append(cmds, s.spin.Tick, tcmds.RetrieveDiscogsCollectionCmd(folder))
		} else if s.loadCollectionForm.Form.State == huh.StateAborted {
			s.loading = false
			s.shelf.Focus()
		}
		return s, tea.Batch(cmds...)

	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height

	case tea.KeyPressMsg:
		sh := s.shelf.PhysicalShelf()
		if sh == nil {
			s.logger.Debug("physical shelf in loadedshelf is nil")
			return s, nil
		}

		switch {
		case key.Matches(msg, s.keys.Next):
			s.shelf = s.shelf.SelectNextBin()
			return s, nil

		case key.Matches(msg, s.keys.Prev):
			s.shelf = s.shelf.SelectPrevBin()
			return s, nil

		case key.Matches(msg, s.keys.Back):
			return s, func() tea.Msg {
				return tcmds.TransitionToMainMenuMsg{}
			}

		case key.Matches(msg, s.keys.Load):
			s.shelf.Blur()
			return s, tcmds.ListDiscogsFoldersCmd()

		case msg.String() == "enter":
			b := s.shelf.GetSelectedBin().PhysicalBin()
			return s, func() tea.Msg {
				return tcmds.TransitionToLoadedBinMsg{
					BinID: b.ID,
				}
			}
		}

	case spinner.TickMsg:
		if s.fetching {
			var cmd tea.Cmd
			s.spin, cmd = s.spin.Update(msg)
			return s, cmd
		}

	case tcmds.GetDiscogsFoldersMsg:
		s.loadCollectionForm = newFolderSelectForm(msg.Folders)
		s.loading = true
		return s, s.loadCollectionForm.Init()

	case tcmds.NewDiscogsCollectionMsg:
		s.releases = msg.Releases
		s.totalReleases = len(msg.Releases)
		s.currentIndex = 0
		s.pct = 0.0
		s.logger.Debug("new discogs collection to process", slog.Int("count", s.totalReleases))
		return s, s.loadNextRecord()

	case loadNextMsg:
		if s.currentIndex >= s.totalReleases {
			s.logger.Debug("finished processing all releases")
			s.fetching = false
			s.shelf.Focus()
			// Reload the shelf from the DB to get all relationships correctly
			return s, tcmds.GetShelfCmd(s.shelf.ID())
		}

		s.logger.Debug("enriching next release", slog.Int("index", s.currentIndex))
		rec := s.releases[s.currentIndex]
		return s, tcmds.EnrichReleaseInstanceCmd(rec)

	case tcmds.ShelvesLoadedMsg:
		if err := msg.Err; err != nil {
			s.logger.Debug("error loading physical shelf",
				slog.Any("err", err),
			)
			return s, nil
		}

		if len(msg.Shelves) != 1 {
			s.logger.Warn("for some reason the database returned multiple shelves, this should be impossible, selecting first result to continume")
		}

		sh := msg.Shelves[0]

		s.shelf = shelf.New(sh, s.logger).
			SelectBin(0)

		return s, nil

	case tcmds.ShelfSavedMsg:
		if msg.Err != nil {
			s.logger.Error("failed to save shelf", slog.String("error", msg.Err.Error()))
		}
		s.currentIndex++
		s.pct = float64(s.currentIndex) / float64(s.totalReleases)
		return s, tea.Batch(s.prog.SetPercent(s.pct), s.loadNextRecord())

	case tcmds.NewDiscogsEnrichRecordMsg:
		if msg.Err != nil {
			s.logger.Error("failed to enrich record, skipping", slog.String("error", msg.Err.Error()))
			s.currentIndex++
			return s, s.loadNextRecord()
		}

		if s.currentIndex < len(s.releases) {
			s.currentTitle = msg.Record.Title
		}

		// This updates the in-memory representation
		_, err := s.shelf.PhysicalShelf().Insert(msg.Record)
		if err != nil {
			s.logger.Error("failed to insert record into shelf model", slog.String("error", err.Error()))
		}

		// Persist the change to the database
		return s, tcmds.SaveShelfCmd(s.shelf.PhysicalShelf())

	}

	var shelfCmd tea.Cmd
	s.shelf, shelfCmd = util.UpdateModel(s.shelf, passthru)
	cmds = append(cmds, shelfCmd)

	return s, tea.Batch(cmds...)
}

func (s LoadedShelfState) View() tea.View {
	return s.renderModel()
}

func (s LoadedShelfState) Help() string {
	return util.FmtKeymap(s.keys.ShortHelp())
}

func (s LoadedShelfState) Type() states.StateType {
	return states.LoadedShelf
}
