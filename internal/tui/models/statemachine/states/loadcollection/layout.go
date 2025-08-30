package loadcollection

import (
	lipgloss "github.com/charmbracelet/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/tui/style"
)

var (
	progressStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)

func (s LoadCollectionState) renderModel() string {
	// Get the rendered form view as a string.
	formView := s.selectFolderForm.View()

	// Use a style to add a border around the form content.
	// We'll place this styled block in the center.
	borderedForm := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(style.LightGreen).
		Render(formView)

	// lipgloss.Place is perfect for centering a block of content
	// within a larger space.
	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		borderedForm,
	)

}

// func newLoadedCollectionFormLayout(base *layout.Div, f *form) *layout.Div {
// 	base.ClearChildren()

// 	// 1) Grab current viewport dimensions (set by t.WindowSizeMsg → base.Resize)
// 	w := base.Width()
// 	h := base.Height()

// 	// 2) Compute box size as 1/3 of viewport
// 	boxW := w / 3
// 	boxH := h / 3

// 	if boxW > w {
// 		boxW = w
// 	}

// 	if boxH > h {
// 		boxH = h
// 	}

// 	// 3) Compute margins so the box is centered
// 	//    leftover width/height divided evenly on both sides
// 	marginV := (h - boxH) / 2
// 	marginH := (w - boxW) / 2

// 	// 4) Create the centered container with border, fixed size, and margins
// 	centerBox, _ := layout.New(layout.Column, lipgloss.NewStyle(),
// 		layout.WithName("centerbox"),
// 		layout.WithBorder(true),
// 		layout.WithFixedWidth(boxW),
// 		layout.WithFixedHeight(boxH),
// 		layout.WithMargin(marginV, marginH, marginV, marginH),
// 	)

// 	// 5) Render the Huh form inside
// 	centerBox.AddChild(&layout.TextNode{
// 		Body: f.View(),
// 	})

// 	// 6) Attach it as the lone child of the viewport
// 	base.AddChild(centerBox)

// 	return base
// }

// func newLoadedCollectionProgressLayout(
// 	base *layout.Div,
// 	pb progress.Model,
// 	sp spinner.Model,
// 	fetching bool,
// 	inserting bool,
// 	pct float64,
// ) (*layout.Div, error) {
// 	base.ClearChildren()

// 	var lines []string

// 	// 1) spinner vs. checkmark
// 	if fetching && !inserting {
// 		// spinner spinning, “fetching…”
// 		spinnerLine := lipgloss.NewStyle().
// 			AlignHorizontal(lipgloss.Center).
// 			Render(fmt.Sprintf("%s fetching collection from Discogs…", sp.View()))
// 		lines = append(lines, spinnerLine)
// 	} else {
// 		// checkmark when fetching→inserting transition is done
// 		check := lipgloss.NewStyle().
// 			Foreground(lipgloss.Color("10")).
// 			Bold(true).
// 			Render("✔ fetching collection complete")
// 		lines = append(lines, check)
// 	}

// 	// 2) always show the progress bar underneath
// 	barLine := lipgloss.NewStyle().
// 		AlignHorizontal(lipgloss.Center).
// 		Render(pb.ViewAs(pct))
// 	lines = append(lines, barLine)

// 	// join the two lines vertically, then wrap in a single TextNode
// 	joined := lipgloss.JoinVertical(lipgloss.Center, lines...)
// 	base.AddChild(&layout.TextNode{
// 		Body: style.ProgressStyle.Render(joined),
// 	})

// 	return base, nil
// }
