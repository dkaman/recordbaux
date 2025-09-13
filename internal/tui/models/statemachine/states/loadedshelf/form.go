package loadedshelf

import (
	"context"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"

	"github.com/dkaman/discogs-golang"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

type loadCollectionForm struct {
	Form   *huh.Form
	folder string
}

func newFolderSelectForm(c *discogs.Client, u string) *loadCollectionForm {
	f := &loadCollectionForm{}

	folders, _ := c.Collection.ListFolders(context.TODO(), u)

	var folderOptions []huh.Option[string]

	for _, fol := range folders {
		name := fol.Name
		o := huh.NewOption(name, name)
		folderOptions = append(folderOptions, o)
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("discogs folder select").
				Options(folderOptions...).
				Value(&f.folder),
		),
	).WithTheme(style.DefaultFormStyles())

	return f
}

func (f *loadCollectionForm) Init() tea.Cmd {
	return f.Form.Init()
}

func (f *loadCollectionForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	mod, cmd := f.Form.Update(msg)
	if frm, ok := mod.(*huh.Form); ok {
		f.Form = frm
	}
	return f, cmd
}

func (f *loadCollectionForm) View() string {
	return f.Form.View()
}

func (f *loadCollectionForm) Folder() string {
	return f.folder
}
