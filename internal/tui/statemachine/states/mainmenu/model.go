package mainmenu

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/statemachine"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

type MainMenuState struct {
	keys        keyMap
	loadedShelf *physical.Shelf
	shelves     list.Model
	nextState   statemachine.StateType
}

func New() MainMenuState {
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

	return MainMenuState{
		keys:        defaultKeybinds(),
		shelves:     shelves,
		loadedShelf: nil,
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
			i, ok := s.shelves.SelectedItem().(*physical.Shelf)
			if ok {
				s.loadedShelf = i
			}
			s.nextState = statemachine.MainMenu
		case key.Matches(msg, s.keys.NewShelf):
			s.nextState = statemachine.CreateShelf
		}
	case NewShelfMsg:
		s.nextState = statemachine.MainMenu

		if msg.Shelf != nil {
			insCmds := s.shelves.InsertItem(0, msg.Shelf)
			cmds = append(cmds, insCmds)
		}
	}

	var listCmd tea.Cmd
	s.shelves, listCmd = s.shelves.Update(msg)
	cmds = append(cmds, listCmd)

	return s, tea.Batch(cmds...)
}

func (s MainMenuState) View() string {
	if len(s.shelves.Items()) == 0 {
		return "no currently defined shelves"
	}

	list := s.shelves.View()

	if s.loadedShelf != nil {
		name := s.loadedShelf.Name
		nBins := len(s.loadedShelf.Bins)
		sz := s.loadedShelf.BinSize
		capacity := nBins * sz

		list = list + fmt.Sprintf(
			"\n\nshelf name: %s\nnum bins: %d\nbin size: %d\n\nshelf %s has a capacity of %d records!",
			name, nBins, sz, name, capacity,
		)
	}

	return list + "\n"
}

func (s MainMenuState) Next(msg tea.Msg) (*statemachine.StateType, error) {
	return &s.nextState, nil
}
