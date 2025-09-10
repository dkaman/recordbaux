package util

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea/v2"
)

func UpdateModel[T tea.Model](m T, msg tea.Msg) (T, tea.Cmd) {
	updatedModel, cmd := m.Update(msg)

	concreteModel, ok := updatedModel.(T)
	if !ok {
		// A model's Update method should ALWAYS return a model of its
		// own type. If it doesn't, it's a critical programming error.
		panic(fmt.Sprintf("model update returned wrong type: expected %T but got %T", m, updatedModel))
	}

	return concreteModel, cmd
}
