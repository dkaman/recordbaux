package shelf

import (
	"fmt"
	"log/slog"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/physical"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/div" // Import the div package
)

// Define constants for bin styling
const (
	// Minimum height for bin content (e.g., for "ID\n0/0" label)
	minBinContentHeight = 3 // ID (1) + count/size (1) + extra buffer (1) for internal text node

	// Default internal padding for each bin's content (applied to div.Div)
	defaultBinDivPaddingHorizontal = 0 // 1 char left + 1 char right
	defaultBinDivPaddingVertical   = 0 // 0 char top + 0 char bottom

	// Default margin around each bin div (applied to div.Div)
	defaultBinDivMarginHorizontal = 0 // 1 char space horizontally
	defaultBinDivMarginVertical   = 0 // 1 char space vertically

	// Default border setting for each bin div
	defaultBinDivBorder = true // Whether each bin div has a border
)

type Model struct {
	selectedBin    int
	physicalShelf  *physical.Shelf
	width          int // Current allocated width for the entire shelf block
	height         int // Current allocated height for the entire shelf block
	logger         *slog.Logger
	binDivStyle    lipgloss.Style           // Base style for individual bins' content
	binDivMargin   div.TopRightBottomLeft   // Margin around each bin div
	binDivPadding  div.TopRightBottomLeft   // Padding inside each bin div (between border and content)
	binDivBorder   bool                     // Whether each bin div has a border
}

func New(p *physical.Shelf, log *slog.Logger) Model {
	logger := log.WithGroup("shelf")

	// Initialize base style for bins (alignments only; border/padding handled by div options)
	baseBinDivStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	// Initialize margin and padding structs for div.Div options
	defaultBinDivMargin := div.TopRightBottomLeft{
		Top:    defaultBinDivMarginVertical,
		Right:  defaultBinDivMarginHorizontal,
		Bottom: defaultBinDivMarginVertical,
		Left:   defaultBinDivMarginHorizontal,
	}

	defaultBinDivPadding := div.TopRightBottomLeft{
		Top:    defaultBinDivPaddingVertical,
		Right:  defaultBinDivPaddingHorizontal,
		Bottom: defaultBinDivPaddingVertical,
		Left:   defaultBinDivPaddingHorizontal,
	}

	return Model{
		selectedBin:    0,
		physicalShelf:  p,
		width:          0,
		height:         0,
		logger:         logger,
		binDivStyle:    baseBinDivStyle,
		binDivMargin:   defaultBinDivMargin,
		binDivPadding:  defaultBinDivPadding,
		binDivBorder:   defaultBinDivBorder,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.logger.Info("shelf message received", slog.Any("msg", msg))

	m.logger.Info("model dimensions",
		slog.Int("width", m.width),
		slog.Int("height", m.height),
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.logger.Info("updated shelf dimensions from tea",
			slog.Int("new_width", m.width),
			slog.Int("new_height", m.height),
		)
	case LoadShelfMsg:
		m.physicalShelf = msg.phy
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	availableWidth := m.width
	availableHeight := m.height

	m.logger.Info("shelf model dimensions in View()",
		slog.Int("width", availableWidth),
		slog.Int("height", availableHeight),
	)

	if m.physicalShelf == nil {
		return "no shelves loaded"
	}
	if availableWidth <= 0 || availableHeight <= 0 {
		return "" // Don't render if no space is available
	}

	shape := m.physicalShelf.Shape
	if shape == nil || len(m.physicalShelf.Bins) == 0 {
		return "shelf has no shape or bins to display"
	}

	numBins := len(m.physicalShelf.Bins)

	// Calculate the combined horizontal space consumed by padding and border for a single bin div
	binInternalHorizontalPaddingAndBorder := (m.binDivPadding.Left + m.binDivPadding.Right)
	if m.binDivBorder {
		binInternalHorizontalPaddingAndBorder += 2 // 1 char for left border, 1 for right
	}

	// Calculate the combined vertical space consumed by padding and border for a single bin div
	binInternalVerticalPaddingAndBorder := (m.binDivPadding.Top + m.binDivPadding.Bottom)
	if m.binDivBorder {
		binInternalVerticalPaddingAndBorder += 2 // 1 char for top border, 1 for bottom
	}

	// Calculate the total horizontal margin consumed by *one* bin div (left + right)
	totalHorizontalMarginPerBin := (m.binDivMargin.Left + m.binDivMargin.Right)

	// Calculate the total vertical margin consumed by *one* bin div (top + bottom)
	totalVerticalMarginPerBin := (m.binDivMargin.Top + m.binDivMargin.Bottom)

	// Minimum height a bin div must have to show content + padding + border
	minBinDivTotalHeight := minBinContentHeight + binInternalVerticalPaddingAndBorder

	// Estimate the *minimum content width* for a bin based on minBinContentHeight and aspect ratio.
	// This is the ideal minimum width for the text content itself.
	estimatedMinContentWidth := minBinContentHeight * 3

	// Calculate the total horizontal space required for a *single bin* if it were rendered at its minimum size,
	// including its content, internal padding, border, and surrounding margins.
	minTotalBinOccupiedWidth := estimatedMinContentWidth + binInternalHorizontalPaddingAndBorder + totalHorizontalMarginPerBin
	if minTotalBinOccupiedWidth < 1 { // Ensure at least 1 to prevent division by zero
		minTotalBinOccupiedWidth = 1
	}

	cols := 1
	rows := 1

	if rect, ok := shape.(*physical.Rectangular); ok {
		cols = rect.X
		rows = rect.Y
	} else {
		// For irregular shelves, determine columns based on available width
		// and the estimated total width a single bin would occupy.
		maxBinsPerRowBasedOnAvailableWidth := availableWidth / minTotalBinOccupiedWidth
		if maxBinsPerRowBasedOnAvailableWidth < 1 {
			maxBinsPerRowBasedOnAvailableWidth = 1
		}
		cols = maxBinsPerRowBasedOnAvailableWidth
		rows = int(math.Ceil(float64(numBins) / float64(cols)))
	}

	if cols <= 0 || rows <= 0 {
		return "shelf shape has invalid dimensions or insufficient space"
	}

	// Now that cols and rows are determined, calculate the effective available space
	// that can be distributed among the actual bin content, padding, and border *for the chosen cols/rows*.
	effectiveAvailableWidthForBinsContentArea := availableWidth - (totalHorizontalMarginPerBin * cols)
	effectiveAvailableHeightForBinsContentArea := availableHeight - (totalVerticalMarginPerBin * rows) // totalVerticalMarginPerBin includes top/bottom margins for each row

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

	// Adjust final bin dimensions to maintain aspect ratio (3:1 assumed for content area)
	// and to respect the minimum total height for the div (content + padding + border).
	// The finalBinWidth and finalBinHeight will be the explicit dimensions passed to div.WithDimensions.
	finalBinHeight := int(math.Min(float64(candidateBinHeight), float64(candidateBinWidth)/3.0))
	if finalBinHeight < minBinDivTotalHeight {
		finalBinHeight = minBinDivTotalHeight
	}
	finalBinWidth := finalBinHeight * 3 // Maintain 3:1 aspect ratio based on the final height

	// Ensure finalBinWidth/Height are not negative
	if finalBinWidth < 0 {
		finalBinWidth = 0
	}
	if finalBinHeight < 0 {
		finalBinHeight = 0
	}

	m.logger.Info("calculated bin dimensions",
		slog.Int("finalBinWidth", finalBinWidth),
		slog.Int("finalBinHeight", finalBinHeight),
		slog.Int("cols", cols),
		slog.Int("rows", rows),
		slog.Int("effectiveAvailableWidthForBinsContentArea", effectiveAvailableWidthForBinsContentArea),
		slog.Int("effectiveAvailableHeightForBinsContentArea", effectiveAvailableHeightForBinsContentArea),
		slog.Int("binInternalHorizontalPaddingAndBorder", binInternalHorizontalPaddingAndBorder),
		slog.Int("binInternalVerticalPaddingAndBorder", binInternalVerticalPaddingAndBorder),
	)

	// Create the main shelf div (column-wise arrangement of rows)
	shelfRootDiv, err := div.New(div.Column, lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center))
	if err != nil {
		m.logger.Error("error creating shelf root div", slog.Any("error", err))
		return "Error rendering shelf."
	}
	// The root div should explicitly take up the available space.
	shelfRootDiv.Resize(availableWidth, availableHeight)

	binIndex := 0
	for r := 0; r < rows; r++ {
		rowDiv, err := div.New(div.Row, lipgloss.NewStyle()) // Row-wise arrangement of bins
		if err != nil {
			m.logger.Error("error creating row div", slog.Any("error", err))
			return "Error rendering shelf."
		}
		// Rows are children of shelfRootDiv (Column direction).
		// Their width should be exactly `availableWidth` so they don't overflow.
		// Their height will be fixed based on the `finalBinHeight` and vertical margins.
		// NOTE: Passing availableWidth to rowDiv.Resize implies that the rowDiv will try
		// to distribute this width to its children (bins). If rowDiv's children (the bins)
		// collectively ask for more space than rowDiv.width, then the internal lipgloss
		// JoinHorizontal will cause wrapping. The goal is to constrain the *number* of
		// bins in a row such that they fit within availableWidth.
		rowDiv.Resize(availableWidth, finalBinHeight + m.binDivMargin.Top + m.binDivMargin.Bottom)


		for c := 0; c < cols; c++ {
			var binNode *div.Div
			var binLabel string
			currentBinSelected := false

			if binIndex < numBins { // Use numBins (total bins) to check if a physical bin exists
				b := m.physicalShelf.Bins[binIndex]
				count := len(b.Records)
				binLabel = fmt.Sprintf("%s\n%d/%d", b.ID, count, b.Size)
				currentBinSelected = (binIndex == m.selectedBin)
			} else {
				// Empty placeholder bin if there are fewer bins than calculated slots
				binLabel = ""
			}

			binInnerStyle := m.binDivStyle.Copy()
			if binIndex < numBins && len(m.physicalShelf.Bins[binIndex].Records) > 0 {
				binInnerStyle = binInnerStyle.BorderBackground(lipgloss.Color("62"))
			}
			if currentBinSelected {
				binInnerStyle = binInnerStyle.BorderForeground(lipgloss.Color("5"))
			}

			binNode, err = div.New(div.Column, binInnerStyle,
				div.WithName(fmt.Sprintf("bin-%d", binIndex)),
				div.WithMargin(m.binDivMargin.Top, m.binDivMargin.Right, m.binDivMargin.Bottom, m.binDivMargin.Left),
				div.WithPadding(m.binDivPadding.Top, m.binDivPadding.Right, m.binDivPadding.Bottom, m.binDivPadding.Left),
			)
			if err != nil {
				m.logger.Error("error creating bin div", slog.Any("error", err))
				return "Error rendering shelf."
			}

			binNode.Resize(finalBinWidth, finalBinHeight)

			if m.binDivBorder {
				binNode.ApplyOption(div.WithBorder(true)) // Apply border if configured
			}

			binNode.AddChild(&div.TextNode{Body: binLabel}) // Add content to the bin div

			rowDiv.AddChild(binNode)
			binIndex++
		}
		shelfRootDiv.AddChild(rowDiv)
	}

	return shelfRootDiv.Render()
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
	return m
}

func (m Model) SelectNextBin() Model {
	numBins := len(m.physicalShelf.Bins)
	if numBins == 0 {
		return m
	}
	m.selectedBin = (m.selectedBin + 1) % numBins
	return m
}

func (m Model) SelectPrevBin() Model {
	numBins := len(m.physicalShelf.Bins)
	if numBins == 0 {
		return m
	}
	m.selectedBin = (m.selectedBin - 1 + numBins) % numBins
	return m
}

func (m Model) GetSelectedBin() bin.Model {
	if m.physicalShelf == nil || len(m.physicalShelf.Bins) == 0 {
		return bin.Model{}
	}
	b := m.physicalShelf.Bins[m.selectedBin]
	return bin.New(b, style.ActiveTextStyle)
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
	cap := bins * m.physicalShelf.Shape.BinSize()
	return fmt.Sprintf("%d bins, capacity %d", bins, cap)
}

func (m Model) PhysicalShelf() *physical.Shelf {
	return m.physicalShelf
}
