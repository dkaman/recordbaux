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
}

const (
	// states
	MainMenu StateType = iota
	LoadedShelf
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
		"LoadedShelf",
		"LoadedBin",
		"SelectShelf",
		"CreatePlaylist",
		"LoadedPlaylist",
		"Quit",
		"Undefined",
	}[s]
}
