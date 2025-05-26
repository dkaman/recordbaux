package loadedbin

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

type refreshLoadedBinMsg struct{}

type LoadedBinState struct {
	app       *app.App
	help      help.Model
	nextState statemachine.StateType
	bin       bin.Model
	keys      keyMap

	records table.Model
}

// New constructs a LoadedBinState ready to receive a LoadShelfMsg
func New(a *app.App) LoadedBinState {
	return LoadedBinState{
		app:       a,
		help:      help.New(),
		nextState: statemachine.Undefined,
		keys:      defaultKeybinds(),
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
			s.nextState = statemachine.LoadedShelf
			return s, tea.Batch(cmds...)
		}
	}

	tableModel, tableUpdateCmds := s.records.Update(msg)
	s.records = tableModel
	cmds = append(cmds,
		tableUpdateCmds,
	)

	return s, tea.Batch(cmds...)
}

func (s LoadedBinState) View() string {
	view := s.records.View()
	return view
}

func (s LoadedBinState) Help() string {
	return s.help.View(s.keys)
}

func (s LoadedBinState) Next() (statemachine.StateType, bool) {
	if s.nextState != statemachine.Undefined {
		return s.nextState, true
	}

	return statemachine.Undefined, false
}

func (s LoadedBinState) Transition() {
	s.nextState = statemachine.Undefined
}
