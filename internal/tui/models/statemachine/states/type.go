package states

import (
	tea "charm.land/bubbletea/v2"
)

type StateType int

type helper interface{
	Help() string
}


type State interface {
	tea.Model
	helper
	Type() StateType
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
