package mainmenu

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	SelectShelf    key.Binding
	NewShelf       key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		SelectShelf: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch shelf"),
		),
		NewShelf: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "new shelf"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.SelectShelf,
		k.NewShelf,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SelectShelf, k.NewShelf},
	}
}
