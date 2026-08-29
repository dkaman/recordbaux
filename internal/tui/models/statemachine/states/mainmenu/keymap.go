package mainmenu

import (
	"charm.land/bubbles/v2/key"
)

type keyMap struct {
	Select      key.Binding
	New         key.Binding
	SwitchFocus key.Binding
	Edit        key.Binding
	Delete      key.Binding
}

func defaultKeybinds() keyMap {
	return keyMap{
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		New: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "new shelf or playlist"),
		),
		SwitchFocus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch focus"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Select,
		k.New,
		k.SwitchFocus,
		k.Edit,
		k.Delete,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.New, k.SwitchFocus, k.Edit, k.Delete},
	}
}
