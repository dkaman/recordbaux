package loadedbin

import (
	"github.com/charmbracelet/bubbles/v2/key"
)

type keyMap struct {
	Back key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		Back: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "back"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back},
	}
}
