package handlers

import (
	"fmt"
	"reflect"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type Handler func(tea.Model, tea.Msg) (tea.Model, tea.Cmd, tea.Msg)

type Registry struct {
	handlers map[reflect.Type]Handler
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[reflect.Type]Handler),
	}
}

func Register[T tea.Model, U tea.Msg](r *Registry, handler func(T, U) (tea.Model, tea.Cmd, tea.Msg)) {
	// Get the reflect.Type of the message. We use a pointer and then
	// .Elem() to correctly get the type for both value and pointer
	// receivers.
	msgType := reflect.TypeOf((*U)(nil)).Elem()

	// Create a closure that wraps the strongly-typed handler. This wrapper
	// is what gets stored in the map.
	wrapper := func(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd, tea.Msg) {
		// Inside the wrapper, perform type assertions. We know the
		// concrete types T and U because they were captured by the
		// closure.
		concreteModel, ok := m.(T)
		if !ok {
			// This is a programming error. The model passed to
			// Update was not the type this handler expected.
			panic(fmt.Sprintf("handler registration mismatch: model is %T, but expected %T", m, *new(T)))
		}

		concreteMsg := msg.(U) // We know this will succeed because we look up by msg type.

		// Call the original, strongly-typed handler function.
		return handler(concreteModel, concreteMsg)
	}

	r.handlers[msgType] = wrapper
}

func (r *Registry) GetHandler(msg tea.Msg) (Handler, bool) {
	h, ok := r.handlers[reflect.TypeOf(msg)]
	return h, ok
}
