package overlay

import (
	"charm.land/bubbles/v2/key"

	tea "charm.land/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/tui/style"
)

type TeaActiveCloser interface {
	tea.Model
	Active() bool
	Close() tea.Cmd
}

type TeaFocusBlurer interface {
	tea.Model
	Focus() TeaFocusBlurer
	Blur() TeaFocusBlurer
}

type Model[M TeaFocusBlurer] struct {
	width, height int
	escape        key.Binding

	Base  M
	Modal TeaActiveCloser
}

func New[M TeaFocusBlurer](base M) Model[M] {
	m := Model[M]{
		Base: base,
	}

	m.escape = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close modal"),
	)

	return m
}

func (m Model[M]) Init() tea.Cmd {
	return m.Base.Init()
}

func (m Model[M]) Update(msg tea.Msg) (Model[M], tea.Cmd) {
	var cmds []tea.Cmd

	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = sizeMsg.Width
		m.height = sizeMsg.Height
	}

	if m.Modal != nil {
		if m.Modal.Active() {
			if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
				if key.Matches(keyMsg, m.escape) {
					closeCmd := m.Modal.Close()
					m.Modal = nil
					return m, closeCmd
				}
			}

			updatedModal, cmd := m.Modal.Update(msg)
			if modalModel, ok := updatedModal.(TeaActiveCloser); ok {
				m.Modal = modalModel
			}
			cmds = append(cmds, cmd)

			return m, tea.Batch(cmds...)

		} else {
			modal := m.Modal
			m.Modal = nil
			return m, modal.Close()
		}
	}

	updatedBase, cmd := m.Base.Update(msg)
	if baseModel, ok := updatedBase.(M); ok {
		m.Base = baseModel
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model[M]) View() tea.View {
	baseView := m.Base.View().Content

	if m.Modal == nil || !m.Modal.Active() {
		return tea.NewView(baseView)
	}

	modalView := m.Modal.View().Content
	composited := style.RenderOverlay(baseView, modalView, m.width, m.height)

	return tea.NewView(composited)
}

func (m *Model[M]) SetModal(modal TeaActiveCloser) {
	m.Modal = modal
}

func (m *Model[M]) ClearModal() tea.Cmd {
	return m.Modal.Close()
}
