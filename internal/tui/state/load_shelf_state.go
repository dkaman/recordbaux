package state

import (
	"fmt"
	"log"
	tea "github.com/charmbracelet/bubbletea"
)

type LoadShelfState struct {

}

func NewLoadShelfState() LoadShelfState {
	return LoadShelfState{}
}

func (s LoadShelfState) Init() tea.Cmd {
	log.Println("in load shelf state init")
	return nil
}

func (s LoadShelfState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("in load shelf state update")
	return nil, nil
}

func (s LoadShelfState) View() string {
	log.Println("in load shelf state view")
	return fmt.Sprint("load shelf\n")
}

func (s LoadShelfState) Next(msg tea.Msg) StateType {
	log.Println("in load shelf state next")
	return Idle
}
