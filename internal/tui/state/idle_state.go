package state

import (
	"fmt"
	"log"
	tea "github.com/charmbracelet/bubbletea"
)

type IdleState struct {

}

func NewIdleState() IdleState {
	return IdleState{

	}
}

func (s IdleState) Init() tea.Cmd {
	log.Println("in idle state init")
	return nil
}

func (s IdleState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("in idle state update")
	return nil, nil
}

func (s IdleState) View() string {
	log.Println("in idle state view")
	return fmt.Sprintf("idle\n")
}

func (s IdleState) Next(msg tea.Msg) StateType {
	log.Println("in idle state next")
	return NewShelf
}
