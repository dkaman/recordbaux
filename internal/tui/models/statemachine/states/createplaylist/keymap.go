package createplaylist

import (
	"github.com/charmbracelet/bubbles/v2/key"
)

type keyMap struct {
	Back   key.Binding
	Select key.Binding
	Create key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		Back: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q/esc", "back"),
		),
		Select: key.NewBinding(
			key.WithKeys("space"), // spacebar
			key.WithHelp("space", "toggle select"),
		),
		Create: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "create playlist"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Select,
		k.Create,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Select, k.Create},
	}
}
