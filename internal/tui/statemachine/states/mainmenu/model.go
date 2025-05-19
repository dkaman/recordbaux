package mainmenu

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
	"github.com/dkaman/recordbaux/internal/tui/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layouts"

	lss "github.com/dkaman/recordbaux/internal/tui/statemachine/states/loadedshelf"
)

type MainMenuState struct {
	keys      keyMap

	// tea models
	loadedShelf shelf.Model
	shelves     list.Model

	// styling/layout
	layout *layouts.TallLayout
}

func New(l *layouts.TallLayout) MainMenuState {
	shelves := list.New([]list.Item{}, shelfDelegate{focused: true}, 0, 10)
	shelves.Title = "shelves"
	shelves.SetStatusBarItemName("shelf", "shelves")

	shelves.Styles.TitleBar = style.TableTitleBarStyle
	shelves.Styles.Title = style.TableTitleStyle
	shelves.Styles.Spinner = style.TableSpinnerStyle
	shelves.Styles.FilterPrompt = style.TableFilterPromptStyle
	shelves.Styles.FilterCursor = style.TableFilterCursorStyle
	shelves.Styles.DefaultFilterCharacterMatch = style.TableDefaultFilterCharacterMatchStyle
	shelves.Styles.StatusBar = style.TableStatusBarStyle
	shelves.Styles.StatusEmpty = style.TableStatusEmptyStyle
	shelves.Styles.StatusBarActiveFilter = style.TableStatusBarActiveFilterStyle
	shelves.Styles.StatusBarFilterCount = style.TableStatusBarFilterCountStyle
	shelves.Styles.NoItems = style.TableNoItemsStyle
	shelves.Styles.ArabicPagination = style.TableArabicPaginationStyle
	shelves.Styles.PaginationStyle = style.TablePaginationStyleStyle
	shelves.Styles.HelpStyle = style.TableHelpStyleStyle
	shelves.Styles.ActivePaginationDot = style.TableActivePaginationDotStyle
	shelves.Styles.InactivePaginationDot = style.TableInactivePaginationDotStyle
	shelves.Styles.DividerDot = style.TableDividerDotStyle

	loadedShelf := shelf.New(nil, style.ActiveTextStyle)

	return MainMenuState{
		keys:        defaultKeybinds(),
		shelves:     shelves,
		loadedShelf: loadedShelf,
		layout:      l,
	}
}

func (s MainMenuState) Init() tea.Cmd {
	return nil
}

func (s MainMenuState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.SelectShelf):
			i, ok := s.shelves.SelectedItem().(shelf.Model)
			if ok {
				cmds = append(cmds,
					lss.WithShelf(i.PhysicalShelf()),
					statemachine.WithNextState(statemachine.LoadedShelf),
				)
			}
		case key.Matches(msg, s.keys.NewShelf):
			cmds = append(cmds, statemachine.WithNextState(statemachine.CreateShelf))
		}
	case NewShelfMsg:
		if msg.Shelf != nil {
			m := shelf.New(msg.Shelf, style.ActiveTextStyle)
			insCmds := s.shelves.InsertItem(0, m)
			cmds = append(cmds, insCmds)
		}
	}

	var listCmd tea.Cmd
	s.shelves, listCmd = s.shelves.Update(msg)
	cmds = append(cmds, listCmd)

	selectedShelfModel := s.shelves.SelectedItem()
	if sel, ok := selectedShelfModel.(shelf.Model); ok {
		s.loadedShelf = sel
	}

	shelfModel, shelfCmds := s.loadedShelf.Update(msg)
	if sh, ok := shelfModel.(shelf.Model); ok {
		s.loadedShelf = sh
	}
	cmds = append(cmds, shelfCmds)

	return s, tea.Batch(cmds...)
}

func (s MainMenuState) View() string {
	list := s.shelves.View()
	shelf := s.loadedShelf.View()

	s.layout.WithSection(layouts.SideBar, list)
	s.layout.WithSection(layouts.TopWindow, shelf)

	view := fmt.Sprintf("state: main menu\n\nshelf list:\n%s\n\nloaded shelf:\n%s\n", list, shelf)

	return view
}

func (s MainMenuState) SelectedShelf() shelf.Model {
	if sh, ok := s.shelves.SelectedItem().(shelf.Model); ok {
		return sh
	}

	return shelf.Model{}
}

func (s MainMenuState) Shelves() list.Model {
	return s.shelves
}

func (s MainMenuState) LoadedShelf() shelf.Model {
	return s.loadedShelf
}
