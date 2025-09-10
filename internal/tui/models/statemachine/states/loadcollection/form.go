package loadcollection

import (
	"context"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"

	"github.com/dkaman/discogs-golang"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

type shape int

const (
	Rect shape = iota
	Irregular
	NotDefined
)

type form struct {
	Form   *huh.Form
	folder string
}

func validateNum(s string) error {
	if s == "" {
		return fmt.Errorf("required")
	}
	if _, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("must be a number")
	}
	return nil
}

func newFolderSelectForm(c *discogs.Client, u string) *form {
	f := &form{}

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

func (f *form) Init() tea.Cmd {
	return f.Form.Init()
}

func (f *form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	mod, cmd := f.Form.Update(msg)
	if frm, ok := mod.(*huh.Form); ok {
		f.Form = frm
	}
	return f, cmd
}

func (f *form) View() string {
	return f.Form.View()
}

func (f *form) Folder() string {
	return f.folder
}
