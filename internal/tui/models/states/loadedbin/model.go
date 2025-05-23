package loadedbin

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/app"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/models/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"

	teaCmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

type refreshLoadedBinMsg struct{}

type LoadedBinState struct {
	app       *app.App
	nextState statemachine.StateType
	bin       bin.Model
	keys      keyMap

	records table.Model
}

// New constructs a LoadedBinState ready to receive a LoadShelfMsg
func New(a *app.App) LoadedBinState {
	columns := []table.Column{
		{Title: "catalog no.", Width: 10},
		{Title: "release name", Width: 30},
		{Title: "artist", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	return LoadedBinState{
		app:       a,
		nextState: statemachine.Undefined,
		keys:      defaultKeybinds(),
		records:   t,
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
		)

		cmds = append(cmds,
			teaCmds.WithLayoutUpdate(layouts.Overlay, s.records.View()),
		)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			cmds = append(cmds,
				teaCmds.WithLayoutUpdate(layouts.Overlay, ""),
			)
			s.nextState = statemachine.LoadedShelf
			return s, tea.Batch(cmds...)
		}
	}

	tableModel, tableUpdateCmds := s.records.Update(msg)
	s.records = tableModel
	cmds = append(cmds,
		tableUpdateCmds,
		teaCmds.WithLayoutUpdate(layouts.Overlay, s.records.View()),
	)

	return s, tea.Batch(cmds...)
}

func (s LoadedBinState) View() string {
	view := s.records.View()
	return view
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
