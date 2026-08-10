package createplaylist

import (

	huh "charm.land/huh/v2"

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
	).WithTheme(huh.ThemeFunc(style.DefaultFormStyles))

	return f
}

func (f *form) Name() string {
	return f.name
}
