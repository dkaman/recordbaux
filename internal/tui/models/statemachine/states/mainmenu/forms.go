package mainmenu

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"

	"github.com/dkaman/recordbaux/internal/tui/style"
)

type shape int

const (
	Rect shape = iota
	Irregular
	NotDefined
)

type createShelfForm struct {
	Form   *huh.Form
	name  string
	shape string

	// square/rect shape
	dimX        string
	dimY        string
	binSizeRect string

	// irregular shape
	numBins    string
	binSizeIrr string
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

func newCreateShelfForm() *createShelfForm {
	f := &createShelfForm{}

	f.Form = huh.NewForm(
		// Page 1: name + shape
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("shelf name").
				Placeholder("newshelf").
				Value(&f.name).
				Validate(huh.ValidateNotEmpty()),
			huh.NewSelect[string]().
				Key("shape").
				Title("shelf shape").
				Options(
					huh.NewOption("square/rect", "rect"),
					huh.NewOption("irregular", "irregular"),
				).
				Value(&f.shape),
		),

		// Page 2: square dims, only when shape == "square/rect"
		huh.NewGroup(
			huh.NewInput().
				Key("bin_dim_x").
				Title("Bin Dimension X").
				Placeholder("3").
				Value(&f.dimX).
				Validate(validateNum),
			huh.NewInput().
				Key("bin_dim_y").
				Title("Bin Dimension Y").
				Placeholder("4").
				Value(&f.dimY).
				Validate(validateNum),
			huh.NewInput().
				Key("bin_size").
				Title("Bin Size").
				Placeholder("50").
				Value(&f.binSizeRect).
				Validate(validateNum),
		).WithHideFunc(func() bool {
			return f.shape != "rect"
		}),

		// Page 3: irregular bins, only when shape == "irregular"
		huh.NewGroup(
			huh.NewInput().
				Key("num_bins").
				Title("Number of Bins").
				Placeholder("12").
				Value(&f.numBins).
				Validate(validateNum),
			huh.NewInput().
				Key("bin_size").
				Title("Bin Size").
				Placeholder("50").
				Value(&f.binSizeIrr).
				Validate(validateNum),
		).WithHideFunc(func() bool {
			return f.shape != "irregular"
		}),
	).WithTheme(style.DefaultFormStyles())

	return f
}

func (f *createShelfForm) Init() tea.Cmd {
	return f.Form.Init()
}

func (f *createShelfForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	mod, cmd := f.Form.Update(msg)
	if frm, ok := mod.(*huh.Form); ok {
		f.Form = frm
	}
	return f, cmd
}

func (f *createShelfForm) View() string {
	return f.Form.View()
}

func (f *createShelfForm) Name() string {
	return f.name
}

func (f *createShelfForm) Shape() shape {
	if f.shape == "rect" {
		return Rect
	} else if f.shape == "irregular" {
		return Irregular

	}

	return NotDefined
}

func (f *createShelfForm) DimX() int {
	x, _ := strconv.Atoi(f.dimX)
	return x
}

func (f *createShelfForm) DimY() int {
	y, _ := strconv.Atoi(f.dimY)
	return y
}

func (f *createShelfForm) BinSize() int {
	if f.Shape() == Rect {
		n, _ := strconv.Atoi(f.binSizeRect)
		return n
	}
	m, _ := strconv.Atoi(f.binSizeIrr)
	return m
}

func (f *createShelfForm) NumBins() int {
	nb, _ := strconv.Atoi(f.numBins)
	return nb
}
