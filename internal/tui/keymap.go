package tui

import (
	"github.com/charmbracelet/bubbles/v2/key"
)

type keyMap struct {
	ToggleHelp key.Binding
	Quit       key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		ToggleHelp: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("C-c", "quit"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.ToggleHelp,
		k.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ToggleHelp, k.Quit},
	}
}
