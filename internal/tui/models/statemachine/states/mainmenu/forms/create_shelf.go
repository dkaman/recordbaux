package forms

import (
	"fmt"
	"strconv"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"

	"github.com/dkaman/recordbaux/internal/tui/style"

	tcmds "github.com/dkaman/recordbaux/internal/tui/cmds"
)

const (
	keyShelfName    = "shelf_name"
	keyShelfShape   = "shape"
	keyShelfBinDimX = "bin_dim_x"
	keyShelfBinDimY = "bin_dim_y"
	keyShelfNumBins = "num_bins"
	keyShelfBinSize = "bin_size"
	keyPlaylistName = "playlist_name"
)

func validateNum(s string) error {
	if s == "" {
		return fmt.Errorf("required")
	}

	if _, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("must be a number")
	}

	return nil
}

// createShelfForm collects the name, shape (rect/irrecgular) and dimensions of
// a shelf before creating a new blank one. irregular shape will just be a
// linear collection of bins rather than trying to arrange in a row/column form.

type CreateShelfForm struct {
	Form   *huh.Form
	active bool
}

func NewCreateShelfForm() CreateShelfForm {
	var shape string

	f := CreateShelfForm{
		active: true,
	}

	f.Form = huh.NewForm(
		// Page 1: name + shape
		huh.NewGroup(
			huh.NewInput().
				Key(keyShelfName).
				Title("shelf name").
				Placeholder("newshelf").
				Validate(huh.ValidateNotEmpty()),
			huh.NewSelect[string]().
				Key(keyShelfShape).
				Title("shelf shape").
				Options(
					huh.NewOption("square/rect", "rect"),
					huh.NewOption("irregular", "irregular"),
				).
				Value(&shape),
		),

		// Page 2: square dims, only when shape == "square/rect"
		huh.NewGroup(
			huh.NewInput().
				Key(keyShelfBinDimX).
				Title("Bin Dimension X").
				Placeholder("3").
				Validate(validateNum),
			huh.NewInput().
				Key(keyShelfBinDimY).
				Title("Bin Dimension Y").
				Placeholder("4").
				Validate(validateNum),
			huh.NewInput().
				Key(keyShelfBinSize).
				Title("Bin Size").
				Placeholder("50").
				Validate(validateNum),
		).WithHideFunc(func() bool {
			return shape != "rect"
		}),

		// Page 3: irregular bins, only when shape == "irregular"
		huh.NewGroup(
			huh.NewInput().
				Key(keyShelfNumBins).
				Title("Number of Bins").
				Placeholder("12").
				Validate(validateNum),
			huh.NewInput().
				Key(keyShelfBinSize).
				Title("Bin Size").
				Placeholder("50").
				Validate(validateNum),
		).WithHideFunc(func() bool {
			return shape != "irregular"
		}),
	).WithTheme(huh.ThemeFunc(style.DefaultFormStyles))

	return f
}

// tea.Model implementation

func (f CreateShelfForm) Init() tea.Cmd {
	return f.Form.Init()
}

func (f CreateShelfForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	mod, cmd := f.Form.Update(msg)
	if frm, ok := mod.(*huh.Form); ok {
		f.Form = frm
	}

	cmds = append(cmds, cmd)

	// 2. Trigger save on form completion
	if f.Form.State == huh.StateCompleted && f.active {
		name := f.Form.GetString(keyShelfName)
		shape := f.Form.GetString(keyShelfShape)
		x, _ := strconv.Atoi(f.Form.GetString(keyShelfBinDimX))
		y, _ := strconv.Atoi(f.Form.GetString(keyShelfBinDimY))
		binSize, _ := strconv.Atoi(f.Form.GetString(keyShelfBinSize))
		numBins, _ := strconv.Atoi(f.Form.GetString(keyShelfNumBins))

		cmds = append(cmds, tcmds.NewShelfCmd(name, shape, x, y, binSize, numBins))

		f.active = false
	}

	return f, tea.Batch(cmds...)
}

func (f CreateShelfForm) View() tea.View {
	return tea.NewView(f.Form.View())
}

// overlay.TeaActiveCloser

func (f CreateShelfForm) Active() bool {
	return f.active
}

func (f CreateShelfForm) Close() tea.Cmd {
	return func() tea.Msg {
		return tcmds.TransitionToMainMenuMsg{}
	}
}
