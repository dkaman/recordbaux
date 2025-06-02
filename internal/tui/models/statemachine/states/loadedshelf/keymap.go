package loadedshelf

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Next key.Binding
	Prev key.Binding
	Back key.Binding
	Load key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		Next: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next bin"),
		),
		Prev: key.NewBinding(
			key.WithKeys("N"),
			key.WithHelp("N", "prev bin"),
		),
		Back: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "back"),
		),
		Load: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "load bin"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Next,
		k.Prev,
		k.Back,
		k.Load,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Back, k.Load},
	}
}
