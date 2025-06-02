package main

// import (
// 	"log"

// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"

// 	"github.com/dkaman/recordbaux/internal/tui/style/layout"
// )

// type model struct {
// 	root        *layout.Div
// 	width, height int
// }

// func initialModel() model {
// 	return model{}
// }

// func (m model) Init() tea.Cmd {
// 	// No initial command—wait for a WindowSizeMsg
// 	return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		if msg.Type == tea.KeyCtrlC {
// 			return m, tea.Quit
// 		}
// 	case tea.WindowSizeMsg:
// 		// Rebuild the Div tree on every resize
// 		m.width = msg.Width
// 		m.height = msg.Height
// 		m.root = buildDivTree(msg.Width, msg.Height)
// 		return m, nil
// 	}
// 	return m, nil
// }

// func (m model) View() string {
// 	if m.root == nil {
// 		return ""
// 	}
// 	return m.root.Render()
// }

// func LayoutTest() {
// 	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func buildDivTree(width, height int) *layout.Div {
// 	headerStyle := lipgloss.NewStyle().
// 		Background(lipgloss.Color("62")) // faint blue

// 	header, _ := layout.New(layout.Row, headerStyle,
// 		layout.WithBorder(true),
// 		layout.WithMargin(0, 0, 0, 0),
// 	)

// 	// 2) Footer (1/3 of height)
// 	footerStyle := lipgloss.NewStyle().
// 		Background(lipgloss.Color("52")) // faint red

// 	footer, _ := layout.New(layout.Row, footerStyle,
// 		layout.WithBorder(true),
// 		layout.WithMargin(0, 0, 0, 0),
// 	)

// 	// 3) Sidebar and Main go into “body” (which is the middle third)
// 	// ── Sidebar
// 	sidebarStyle := lipgloss.NewStyle().
// 		Background(lipgloss.Color("22")) // faint green

// 	sidebar, _ := layout.New(layout.Row, sidebarStyle,
// 		layout.WithBorder(false),
// 	)

// 	// ── Main, which itself splits into two stacked sections (A above B)
// 	sectionAStyle := lipgloss.NewStyle().
// 		Background(lipgloss.Color("94")) // faint magenta
// 	sectionA, _ := layout.New(layout.Row, sectionAStyle,
// 		layout.WithBorder(false),
// 	)

// 	sectionBStyle := lipgloss.NewStyle().
// 		Background(lipgloss.Color("130")) // faint yellow
// 	sectionB, _ := layout.New(layout.Row, sectionBStyle,
// 		layout.WithBorder(false),
// 	)

// 	mainInner, _ := layout.New(layout.Column, lipgloss.NewStyle(),
// 		layout.WithChild(sectionA),
// 		layout.WithChild(sectionB),
// 		layout.WithBorder(false),
// 	)

// 	body, _ := layout.New(layout.Row, lipgloss.NewStyle(),
// 		layout.WithChild(sidebar),
// 		layout.WithChild(mainInner),
// 		layout.WithBorder(true),
// 		layout.WithMargin(1, 0, 1, 0),
// 	)

// 	// 4) Put Header, Body, Footer into one big vertical container
// 	outer, _ := layout.New(layout.Column, lipgloss.NewStyle(),
// 		layout.WithChild(header),
// 		layout.WithChild(body),
// 		layout.WithChild(footer),
// 		layout.WithBorder(true),
// 		layout.WithPadding(1, 1, 1, 1),
// 	)

// 	// Now that the tree is built, run Resize on the root to propagate down:
// 	outer.Resize(width, height)
// 	return outer
// }
