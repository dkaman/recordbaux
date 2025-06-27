package states

import (
	tea "github.com/charmbracelet/bubbletea"
)

type StateType int

type helper interface{
	Help() string
}

type State interface {
	tea.Model
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
	LoadedBin
	SelectShelf
	Quit
	Undefined
)

func (s StateType) String() string {
	return [...]string{
		"MainMenu",
		"CreateShelf",
		"LoadedShelf",
		"LoadCollection",
		"LoadedBin",
		"SelectShelf",
		"Quit",
		"Undefined",
	}[s]
}
