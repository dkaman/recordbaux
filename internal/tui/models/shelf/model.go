package shelf

import (
	"fmt"
	"log/slog"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/shelf"
	"github.com/dkaman/recordbaux/internal/tui/models/bin"
	"github.com/dkaman/recordbaux/internal/tui/style"
	"github.com/dkaman/recordbaux/internal/tui/style/layout"
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
	id            db.ID
	selectedBin   int
	physicalShelf *shelf.Entity

	bins []bin.Model

	// screen dimensions not shelf dimensions
	width  int // Current allocated width for the entire shelf block
	height int // Current allocated height for the entire shelf block

	binDivStyle   lipgloss.Style            // Base style for individual bins' content
	binDivMargin  layout.TopRightBottomLeft // Margin around each bin div
	binDivPadding layout.TopRightBottomLeft // Padding inside each bin div (between border and content)
	binDivBorder  bool                      // Whether each bin div has a border

	logger *slog.Logger
}

func New(p *shelf.Entity, log *slog.Logger) Model {
	logger := log.WithGroup("shelf")

	// Initialize base style for bins (alignments only; border/padding handled by div options)
	baseBinDivStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	// Initialize margin and padding structs forlayout.Div options
	defaultBinDivMargin := layout.TopRightBottomLeft{
		Top:    defaultBinDivMarginVertical,
		Right:  defaultBinDivMarginHorizontal,
		Bottom: defaultBinDivMarginVertical,
		Left:   defaultBinDivMarginHorizontal,
	}

	defaultBinDivPadding := layout.TopRightBottomLeft{
		Top:    defaultBinDivPaddingVertical,
		Right:  defaultBinDivPaddingHorizontal,
		Bottom: defaultBinDivPaddingVertical,
		Left:   defaultBinDivPaddingHorizontal,
	}

	m := Model{
		id:            p.ID,
		selectedBin:   0,
		physicalShelf: p,
		width:         0,
		height:        0,
		logger:        logger,
		binDivStyle:   baseBinDivStyle,
		binDivMargin:  defaultBinDivMargin,
		binDivPadding: defaultBinDivPadding,
		binDivBorder:  defaultBinDivBorder,
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
		m = m.loadPhysicalShelf(msg.phy)

		m.logger.Info("model", slog.Any("m", m))
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

	// if no physical shelf return message (likely won't happen)
	if m.physicalShelf == nil {
		return "no shelves loaded"
	}

	// don't render if no space is available
	if availableWidth <= 0 || availableHeight <= 0 {
		return ""
	}

	if len(m.physicalShelf.AllBins()) == 0 {
		return "shelf has no bins to display"
	}

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

	// Estimate the minimum total width needed for a bin including content, padding, border, and its own margins.
	// This helps in determining how many bins can fit in a row initially.
	// Assume a 3:1 aspect ratio for the content area for this estimate.
	minBinTotalRenderWidthEstimate := (minBinContentHeight * 3) + binInternalHorizontalPaddingAndBorder + totalHorizontalMarginPerBin
	if minBinTotalRenderWidthEstimate < 1 { // Ensure at least 1 to avoid division by zero
		minBinTotalRenderWidthEstimate = 1
	}

	cols := m.physicalShelf.DimX()
	rows := m.physicalShelf.DimY()

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
	// We calculate the content width that would maintain the 3:1 aspect ratio based on candidate height,
	// then add back padding and border.
	finalBinHeight := int(math.Min(float64(candidateBinHeight), float64(candidateBinWidth)/3.0))
	if finalBinHeight < minBinDivTotalHeight {
		finalBinHeight = minBinDivTotalHeight
	}
	finalBinWidth := finalBinHeight * 3 // Maintain 3:1 aspect ratio for content, padding and border area

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

	root, _ := layout.New(layout.Column, style.Centered)
	root.Resize(m.width, m.height)

	binIndex := 0
	for r := 0; r < rows; r++ {
		rowDiv, _ := layout.New(layout.Row, lipgloss.NewStyle())
		rowDiv.Resize(
			(finalBinWidth+m.binDivMargin.Left+m.binDivMargin.Right)*cols,
			finalBinHeight+m.binDivMargin.Top+m.binDivMargin.Bottom,
		)

		for c := 0; c < cols; c++ {
			if binIndex < len(m.bins) {
				b := m.bins[binIndex].
					SetSize(finalBinWidth, finalBinHeight)

				rowDiv.AddChild(&layout.TextNode{
					Body: b.View(),
				})
			}
			binIndex++
		}
		root.AddChild(rowDiv)
	}
	return root.Render()
}

func (m Model) loadPhysicalShelf(s *shelf.Entity) Model {
	m.physicalShelf = s

	m.bins = nil

	alignedSelected := style.Centered.Bold(true).Foreground(style.LightGreen)

	binStyles := bin.Style{
		EmptySelected:   alignedSelected,
		EmptyUnselected: style.Centered,
		FullSelected:    alignedSelected.Background(style.DarkBlue),
		FullUnselected:  style.Centered.Background(style.DarkBlue).Foreground(style.DarkBlack),
	}

	for _, pb := range m.physicalShelf.AllBins() {
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

	bins := len(m.physicalShelf.AllBins())
	cap := bins * m.physicalShelf.BinSize

	return fmt.Sprintf("%d bins, capacity %d", bins, cap)
}

func (m Model) PhysicalShelf() *shelf.Entity {
	return m.physicalShelf
}

func (m Model) ID() db.ID {
	return m.id
}
