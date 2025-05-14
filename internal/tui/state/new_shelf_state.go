package state

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type NewShelfState struct {

}

func NewNewShelfState() NewShelfState {
	return NewShelfState{}
}

func (s NewShelfState) Init() tea.Cmd {
	log.Println("in new shelf state init")
	return nil
}

func (s NewShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("in new shelf state update")
	return nil, nil
}

func (s NewShelfState) View() string {
	log.Println("in new shelf state view")
	return fmt.Sprint("new shelf\n")
}

func (s NewShelfState) Next(msg tea.Msg) StateType {
	log.Println("in new shelf state next")
	return LoadShelf
}
