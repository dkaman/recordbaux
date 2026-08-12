package shelf

import (
	"fmt"
	"log/slog"
	"math"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/style"
)

// Define constants for bin styling
const (
	// Minimum height for bin content (e.g., for "ID\n0/0" label)
	minBinContentHeight = 3 // ID (1) + count/size (1) + extra buffer (1) for internal text node

	// Default internal padding for each bin's content
	defaultBinDivPaddingHorizontal = 0
	defaultBinDivPaddingVertical   = 0

	// Default margin around each bin div
	defaultBinDivMarginHorizontal = 0
	defaultBinDivMarginVertical   = 0

	// Default border setting for each bin div
	defaultBinDivBorder = true
)

type topRightBottomLeft struct{ Top, Right, Bottom, Left int }

type option func(Model) Model

type Model struct {
	id     uint
	logger *slog.Logger

	selectedBin   int
	physicalShelf *shelf.Entity
	bins          []bin.Model

	width   int
	height  int
	focused bool
}

func New(p *shelf.Entity, log *slog.Logger, opts ...option) Model {
	logger := log.WithGroup("shelf")

	m := Model{
		id:            p.ID,
		selectedBin:   0,
		width:         0,
		height:        0,
		logger:        logger,
		focused:       true,
	}

	if p != nil {
		m.physicalShelf = p
		m.loadPhysicalShelf()
	}

	for _, o := range opts {
		m = o(m)
	}

	return m
}

func (m *Model) Focus() {
	m.focused = true
	m.loadPhysicalShelf()

}

func (m *Model) Blur() {
	m.focused = false
	m.loadPhysicalShelf()
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// update was called on a non-initialized model
	if m.id == 0 {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case LoadShelfMsg:
		m.physicalShelf = msg.Phy
		m.loadPhysicalShelf()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	if m.id == 0 {
		return tea.NewView("")
	}

	availableWidth := m.width
	availableHeight := m.height

	// if no physical shelf return message (likely won't happen)
	if m.physicalShelf == nil {
		return tea.NewView("no shelves loaded")
	}

	// don't render if no space is available
	if availableWidth <= 0 || availableHeight <= 0 {
		return tea.NewView("")
	}

	if len(m.physicalShelf.Bins) == 0 {
		return tea.NewView("shelf has no bins to display")
	}

	binDivMargin := topRightBottomLeft{
		Top:    defaultBinDivMarginVertical,
		Right:  defaultBinDivMarginHorizontal,
		Bottom: defaultBinDivMarginVertical,
		Left:   defaultBinDivMarginHorizontal,
	}

	binDivPadding := topRightBottomLeft{
		Top:    defaultBinDivPaddingVertical,
		Right:  defaultBinDivPaddingHorizontal,
		Bottom: defaultBinDivPaddingVertical,
		Left:   defaultBinDivPaddingHorizontal,
	}

	binDivBorder := defaultBinDivBorder

	// Calculate the combined horizontal space consumed by padding and border for a single bin div
	binInternalHorizontalPaddingAndBorder := (binDivPadding.Left + binDivPadding.Right)
	if binDivBorder {
		binInternalHorizontalPaddingAndBorder += 2 // 1 char for left border, 1 for right
	}

	// Calculate the combined vertical space consumed by padding and border for a single bin div
	binInternalVerticalPaddingAndBorder := (binDivPadding.Top + binDivPadding.Bottom)
	if binDivBorder {
		binInternalVerticalPaddingAndBorder += 2 // 1 char for top border, 1 for bottom
	}

	// Calculate the total horizontal margin consumed by *one* bin div (left + right)
	totalHorizontalMarginPerBin := (binDivMargin.Left + binDivMargin.Right)

	// Calculate the total vertical margin consumed by *one* bin div (top + bottom)
	totalVerticalMarginPerBin := (binDivMargin.Top + binDivMargin.Bottom)

	// Minimum height a bin div must have to show content + padding + border
	minBinDivTotalHeight := minBinContentHeight + binInternalVerticalPaddingAndBorder

	s, err := m.physicalShelf.GetShape()
	if err != nil {
		m.logger.Error("error getting shape from entity",
			slog.Any("error", err),
		)
		return tea.NewView("error: could not read shelf shape")
	}

	cols := s.X
	rows := s.Y

	if cols <= 0 || rows <= 0 {
		return tea.NewView("shelf shape has invalid dimensions or insufficient space")
	}

	// Now calculate the effective available space for the *bin content, padding, and border*
	// considering the determined `cols` and `rows` and all margins.
	effectiveAvailableWidthForBinsContentArea := availableWidth - (totalHorizontalMarginPerBin * cols)
	effectiveAvailableHeightForBinsContentArea := availableHeight - (totalVerticalMarginPerBin * rows)

	// Ensure these effective available dimensions are not negative
	if effectiveAvailableWidthForBinsContentArea < 0 {
		effectiveAvailableWidthForBinsContentArea = 0
	}

	if effectiveAvailableHeightForBinsContentArea < 0 {
		effectiveAvailableHeightForBinsContentArea = 0
	}

	// Calculate the "candidate" width and height for each bin div's *total area*
	// (including its own padding and border but not outer margins).
	candidateBinWidth := effectiveAvailableWidthForBinsContentArea / cols
	candidateBinHeight := effectiveAvailableHeightForBinsContentArea / rows

	// Adjust final bin height to maintain aspect ratio (3:1 assumed for content area)
	// and to respect the minimum total height for the div (content + padding + border).
	finalBinHeight := int(math.Min(float64(candidateBinHeight), float64(candidateBinWidth)/3.0))
	if finalBinHeight < minBinDivTotalHeight {
		finalBinHeight = minBinDivTotalHeight
	}

	finalBinWidth := finalBinHeight * 3 // Maintain 3:1 aspect ratio

	if finalBinWidth < 0 {
		finalBinWidth = 0
	}

	if finalBinHeight < 0 {
		finalBinHeight = 0
	}

	marginStyle := lipgloss.NewStyle().
		MarginTop(binDivMargin.Top).
		MarginRight(binDivMargin.Right).
		MarginBottom(binDivMargin.Bottom).
		MarginLeft(binDivMargin.Left)

	var gridRows []string
	binIndex := 0

	// Build the grid row by row
	for r := 0; r < rows; r++ {
		var rowBins []string
		for c := 0; c < cols; c++ {
			var renderedBin string

			if binIndex < len(m.bins) {
				// Get the rendered view string from the bin model
				b := m.bins[binIndex].SetSize(finalBinWidth, finalBinHeight)
				renderedBin = b.View().Content
			} else {
				// If the shelf isn't completely full, render an empty box to preserve grid shape
				renderedBin = lipgloss.NewStyle().
					Width(finalBinWidth).
					Height(finalBinHeight).
					Render("")
			}

			// Wrap the bin string in the margin style and add to the current row
			rowBins = append(rowBins, marginStyle.Render(renderedBin))
			binIndex++
		}

		// Join all bins in this row horizontally
		rowStr := lipgloss.JoinHorizontal(lipgloss.Top, rowBins...)
		gridRows = append(gridRows, rowStr)
	}

	// Join all rows vertically
	grid := lipgloss.JoinVertical(lipgloss.Left, gridRows...)

	// 4. Place the completed grid in the absolute center of the available terminal space
	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, grid))
}

func (m *Model) loadPhysicalShelf() {
	s := m.physicalShelf
	if s == nil {
		return
	}

	m.bins = nil

	// Choose style set based on the `focused` flag
	var binStyles bin.Style
	if m.focused {
		alignedSelected := style.Centered.
			Bold(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(style.LightWhite).
			Foreground(style.LightGreen)
		binStyles = bin.Style{
			EmptySelected:   alignedSelected,
			EmptyUnselected: style.Centered.BorderStyle(lipgloss.NormalBorder()).BorderForeground(style.LightWhite),
			FullSelected:    alignedSelected.Background(style.DarkBlue),
			FullUnselected:  style.Centered.BorderStyle(lipgloss.NormalBorder()).BorderForeground(style.LightWhite).Background(style.DarkBlue).Foreground(style.DarkBlack),
		}
	} else { // Blurred styles
		alignedSelected := style.Centered.Foreground(style.LightGreenDimmed).BorderForeground(style.LightWhiteDimmed)
		binStyles = bin.Style{
			EmptySelected:   alignedSelected,
			EmptyUnselected: style.Centered.Foreground(style.DarkWhiteDimmed).BorderStyle(lipgloss.NormalBorder()).BorderForeground(style.LightWhiteDimmed),
			FullSelected:    alignedSelected.Background(style.DarkBlueDimmed),
			FullUnselected:  style.Centered.Foreground(style.DarkBlackDimmed).Background(style.DarkBlueDimmed).BorderStyle(lipgloss.NormalBorder()).BorderForeground(style.LightWhiteDimmed),
		}
	}

	for _, pb := range s.Bins {
		b := bin.New(pb, binStyles)
		m.bins = append(m.bins, b)
	}
}

func (m Model) SetSize(w, h int) Model {
	m.width = w
	m.height = h
	return m
}

func (m Model) SelectBin(b int) Model {
	numBins := len(m.physicalShelf.Bins)
	if numBins == 0 {
		return m
	}
	m.selectedBin = b % numBins

	bin := m.bins[m.selectedBin]

	m.bins[m.selectedBin] = bin.Select()

	return m
}

func (m Model) SelectNextBin() Model {
	numBins := len(m.physicalShelf.Bins)
	if numBins == 0 {
		return m
	}
	m.bins[m.selectedBin] = m.bins[m.selectedBin].Unselect()

	m.selectedBin = (m.selectedBin + 1) % numBins

	m.bins[m.selectedBin] = m.bins[m.selectedBin].Select()

	return m
}

func (m Model) SelectPrevBin() Model {
	numBins := len(m.physicalShelf.Bins)
	if numBins == 0 {
		return m
	}
	m.bins[m.selectedBin] = m.bins[m.selectedBin].Unselect()

	m.selectedBin = (m.selectedBin - 1 + numBins) % numBins

	m.bins[m.selectedBin] = m.bins[m.selectedBin].Select()
	return m
}

func (m Model) GetSelectedBin() bin.Model {
	return m.bins[m.selectedBin]
}

func (m Model) Title() string {
	if m.physicalShelf == nil {
		return ""
	}
	return m.physicalShelf.Name
}

func (m Model) FilterValue() string {
	if m.physicalShelf == nil {
		return ""
	}
	return m.physicalShelf.Name
}

func (m Model) Description() string {
	if m.physicalShelf == nil {
		return ""
	}

	bins := len(m.physicalShelf.Bins)
	cap := bins * m.physicalShelf.BinSize

	return fmt.Sprintf("%d bins, capacity %d", bins, cap)
}

func (m Model) PhysicalShelf() *shelf.Entity {
	return m.physicalShelf
}

func (m Model) ID() uint {
	return m.id
}
