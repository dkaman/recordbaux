// flist = focusable list, small wrapper around bubble's list model
package flist

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/style"
)

type FocusableDelegate interface {
	list.ItemDelegate
	Focus() FocusableDelegate
	Blur() FocusableDelegate
}

type Model struct {
	list.Model

	focused  bool
	delegate FocusableDelegate
}

func New(items []list.Item, delegate FocusableDelegate) Model {
	l := list.New(items, delegate, 0, 0)
	l.Styles = style.DefaultListStyles()
	return Model{
		Model:    l,
		focused:  false,
		delegate: delegate,
	}
}

func (m Model) Focus() Model {
	m.focused = true
	m.Styles = style.DefaultListStyles()
	m.delegate = m.delegate.Focus()
	m.Model.SetDelegate(m.delegate)
	return m
}

func (m Model) Blur() Model {
	m.focused = false
	m.Styles = style.DefaultListStylesDimmed()
	m.delegate = m.delegate.Blur()
	m.Model.SetDelegate(m.delegate)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Update only processes messages if the list is focused.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// Only handle inputs if focused. Always handle resize messages.
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.Model.SetSize(sizeMsg.Width, sizeMsg.Height)
		return m, nil
	}

	// For key messages, prevent list filtering when not focused
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok && !m.focused {
		if key.Matches(keyMsg, m.KeyMap.Filter) {
			return m, nil
		}
	}

	m.Model, cmd = m.Model.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	return tea.NewView(m.Model.View())
}
