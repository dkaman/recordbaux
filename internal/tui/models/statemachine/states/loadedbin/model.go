package loadedbin

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine/states"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/div"

	keyFmt "github.com/dkaman/recordbaux/internal/tui/key"
)

type refreshLoadedBinMsg struct{}

type LoadedBinState struct {
	app       *app.App
	nextState states.StateType
	bin       bin.Model
	keys      keyMap
	layout    *div.Div

	records table.Model
}

// New constructs a LoadedBinState ready to receive a LoadShelfMsg
func New(a *app.App, l *div.Div) LoadedBinState {
	h := help.New()
	h.Styles = style.DefaultHelpStyles()

	return LoadedBinState{
		app:       a,
		nextState: states.Undefined,
		keys:      defaultKeybinds(),
		layout:    l,
	}
}

func (s LoadedBinState) Init() tea.Cmd {
	return func() tea.Msg {
		return refreshLoadedBinMsg{}
	}
}

func (s LoadedBinState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case refreshLoadedBinMsg:
		s.bin = s.app.CurrentBin

		columns := []table.Column{
			{Title: "catalog no.", Width: 15},
			{Title: "release name", Width: 50},
			{Title: "artist", Width: 30},
		}

		var rows []table.Row

		for _, r := range s.bin.PhysicalBin().Records {
			catno := r.Release.BasicInfo.Labels[0].CatNo
			name := r.Release.BasicInfo.Title
			artist := r.Release.BasicInfo.Artists[0].Name
			row := table.Row{catno, name, artist}
			rows = append(rows, row)
		}

		s.records = table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(10),
			table.WithStyles(style.DefaultTableStyles()),
		)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			s.nextState = states.LoadedShelf
			return s, tea.Batch(cmds...)
		}
	}

	tableModel, tableUpdateCmds := s.records.Update(msg)
	s.records = tableModel
	cmds = append(cmds,
		tableUpdateCmds,
	)

	s.layout, _  = newLoadedBinLayout(s.layout, s.records)

	return s, tea.Batch(cmds...)
}

func (s LoadedBinState) View() string {
	return s.layout.Render()
}

func (s LoadedBinState) Help() string {
	return keyFmt.FmtKeymap(s.keys.ShortHelp())
}

func (s LoadedBinState) Next() (states.StateType, bool) {
	if s.nextState != states.Undefined {
		return s.nextState, true
	}

	return states.Undefined, false
}

func (s LoadedBinState) Transition() {
	s.nextState = states.Undefined
}
