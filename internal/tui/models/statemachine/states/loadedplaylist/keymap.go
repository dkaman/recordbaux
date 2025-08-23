package loadedplaylist

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Back     key.Binding
	Checkout key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		Back: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q/esc", "back"),
		),
		Checkout: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "checkout"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Checkout}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Back, k.Checkout}}
}
