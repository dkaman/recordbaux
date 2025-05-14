package createshelf

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
)

type shape int

const (
	Rect shape = iota
	Irregular
	NotDefined
)

type form struct {
	*huh.Form
	name    string
	shape   string
	dimX    string
	dimY    string
	binSize string
	numBins string
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

func newShelfCreateForm() *form {
	f := &form{}

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
				Value(&f.binSize).
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
				Value(&f.binSize).
				Validate(validateNum),
		).WithHideFunc(func() bool {
			return f.shape != "irregular"
		}),
	)

	return f
}

func (f *form) Name() string {
	return f.name
}


func (f *form) Shape() shape {
	if f.shape == "rect" {
		return Rect
	} else if f.shape == "irregular" {
		return Irregular

	}

	return NotDefined
}

func (f *form) DimX() int {
	x, _ := strconv.Atoi(f.dimX)
	return x
}

func (f *form) DimY() int {
	y, _ := strconv.Atoi(f.dimY)
	return y
}

func (f *form) BinSize() int {
	bs, _ := strconv.Atoi(f.binSize)
	return bs
}

func (f *form) NumBins() int {
	nb, _ := strconv.Atoi(f.numBins)
	return nb
}
