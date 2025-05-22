package selectshelf

import (
	"github.com/charmbracelet/huh"

	"github.com/dkaman/recordbaux/internal/tui/models/shelf"
)

type form struct {
	*huh.Form
	shelf string
}

func newShelfSelectForm(shelves []shelf.Model) *form {
	f := &form{}


	var shelfOptions []huh.Option[string]

	for _, sh := range shelves {
		name := sh.PhysicalShelf().Name
		o := huh.NewOption(name, name)
		shelfOptions = append(shelfOptions, o)
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("discogs folder select").
				Options(shelfOptions...).
				Value(&f.shelf),
		),
	)

	return f
}

func (f *form) Shelf() string {
	return f.shelf
}
