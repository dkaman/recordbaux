package shelf

import (
	"fmt"
	"log/slog"
	"math"

	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"

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

type Model struct {
	id            uint
	selectedBin   int
	physicalShelf *shelf.Entity
	bins          []bin.Model

	width  int
	height int
	logger *slog.Logger
}

func New(p *shelf.Entity, log *slog.Logger) Model {
	logger := log.WithGroup("shelf")

	m := Model{
		id:            p.ID,
		selectedBin:   0,
		physicalShelf: p,
		width:         0,
		height:        0,
		logger:        logger,
	}

	if p != nil {
		m = m.loadPhysicalShelf(p)
	}

	return m
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
		m = m.loadPhysicalShelf(msg.phy)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.id == 0 {
		return ""
	}

	availableWidth := m.width
	availableHeight := m.height
	// if no physical shelf return message (likely won't happen)
	if m.physicalShelf == nil {
		return "no shelves loaded"
	}
	// don't render if no space is available
	if availableWidth <= 0 || availableHeight <= 0 {
		return ""
	}
	if len(m.physicalShelf.Bins) == 0 {
		return "shelf has no bins to display"
	}

	// Local layout values, previously stored in the struct for the old layout system.
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
		return "error: could not read shelf shape"
	}

	cols := s.X
	rows := s.Y
	if cols <= 0 || rows <= 0 {
		return "shelf shape has invalid dimensions or insufficient space"
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

	// --- New Rendering Logic using lipgloss.Canvas ---

	// Calculate total grid dimensions for centering
	totalGridWidth := cols * (finalBinWidth + totalHorizontalMarginPerBin)
	totalGridHeight := rows * (finalBinHeight + totalVerticalMarginPerBin)

	// Calculate offsets to center the entire grid of bins
	offsetX := (m.width - totalGridWidth) / 2
	if offsetX < 0 {
		offsetX = 0
	}
	offsetY := (m.height - totalGridHeight) / 2
	if offsetY < 0 {
		offsetY = 0
	}

	// Create a canvas to draw on.
	canvas := lipgloss.NewCanvas()

	binIndex := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if binIndex < len(m.bins) {
				// Calculate the top-left position for the current bin's layer
				xPos := offsetX + c*(finalBinWidth+totalHorizontalMarginPerBin)
				yPos := offsetY + r*(finalBinHeight+totalVerticalMarginPerBin)

				// Get the rendered view string from the bin model
				b := m.bins[binIndex].SetSize(finalBinWidth, finalBinHeight)
				binView := b.View()

				// Create a new layer with the bin's content and place it on the canvas
				binLayer := lipgloss.NewLayer(binView)
				canvas.AddLayers(binLayer.X(xPos).Y(yPos))
			}
			binIndex++
		}
	}

	return canvas.Render()
}

func (m Model) loadPhysicalShelf(s *shelf.Entity) Model {
	m.physicalShelf = s

	m.bins = nil

	alignedSelected := style.Centered.
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		Foreground(style.LightGreen)

	binStyles := bin.Style{
		EmptySelected:   alignedSelected,
		EmptyUnselected: style.Centered.BorderStyle(lipgloss.NormalBorder()),
		FullSelected:    alignedSelected.Background(style.DarkBlue),
		FullUnselected:  style.Centered.BorderStyle(lipgloss.NormalBorder()).Background(style.DarkBlue).Foreground(style.DarkBlack),
	}

	for _, pb := range m.physicalShelf.Bins {
		b := bin.New(pb, binStyles)
		m.bins = append(m.bins, b)
	}

	return m
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
