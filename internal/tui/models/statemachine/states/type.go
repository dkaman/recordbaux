package states

import (
	tea "github.com/charmbracelet/bubbletea/v2"
)

type StateType int

type helper interface{
	Help() string
}

type State interface {
	tea.Model
	tea.ViewModel
	helper
	Next() (StateType, bool)
	Transition() State
}

const (
	// states
	MainMenu StateType = iota
	CreateShelf
	LoadedShelf
	LoadCollection
	FetchFromDiscogs
	LoadedBin
	SelectShelf
	CreatePlaylist
	LoadedPlaylist
	Quit
	Undefined
)

func (s StateType) String() string {
	return [...]string{
		"MainMenu",
		"CreateShelf",
		"LoadedShelf",
		"LoadCollection",
		"FetchFromDiscogs",
		"LoadedBin",
		"SelectShelf",
		"CreatePlaylist",
		"LoadedPlaylist",
		"Quit",
		"Undefined",
	}[s]
}
