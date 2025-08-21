package mainmenu

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Select      key.Binding
	NewShelf    key.Binding
	SwitchFocus key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		NewShelf: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "new shelf"),
		),
		SwitchFocus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch focus"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Select,
		k.NewShelf,
		k.SwitchFocus,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.NewShelf, k.SwitchFocus},
	}
}
