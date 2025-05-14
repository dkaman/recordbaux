package main_menu_state

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	NextShelf   key.Binding
	PrevShelf   key.Binding
	SelectShelf key.Binding
	NewShelf    key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		NextShelf:   key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/down", "next shelf"),
		),
		PrevShelf:   key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/up", "previous shelf"),
		),
		SelectShelf: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "load shelf"),
		),
		NewShelf:    key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "new shelf"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.NextShelf,
		k.PrevShelf,
		k.SelectShelf,
		k.NewShelf,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextShelf, k.PrevShelf, k.SelectShelf, k.NewShelf},
	}
}
