package createplaylist

import (

	"github.com/charmbracelet/huh"

	"github.com/dkaman/recordbaux/internal/tui/style"
)

type form struct {
	*huh.Form
	name string
}

func newNameForm() *form {
	f := &form{}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Playlist Name").
				Value(&f.name).
				Validate(huh.ValidateNotEmpty()),
		),
	).WithTheme(style.DefaultFormStyles())

	return f
}

func (f *form) Name() string {
	return f.name
}
