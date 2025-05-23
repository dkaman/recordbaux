package statemachine

type StateType int

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
