package selectshelf

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Select key.Binding
	Back   key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select shelf"),
		),
		Back: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "go back"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Select,
		k.Back,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Back},
	}
}
