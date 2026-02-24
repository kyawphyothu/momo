package program

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.PasteMsg:
		return m.handlePaste(msg)

	case tea.KeyMsg:
		if m.URLOverlayOpen {
			return m.handleURLOverlay(msg)
		}

		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "l":
			m.URLOverlayOpen = true
			m.URLTextInput.Focus()
			return m, nil
		default:
			return m, nil
		}
	}

	return m, nil
}

func (m Model) handlePaste(msg tea.PasteMsg) (Model, tea.Cmd) {
	if m.URLOverlayOpen {
		return m.handlePasteURLOverlay(msg)
	}

	return m, nil
}
func (m Model) handlePasteURLOverlay(msg tea.PasteMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.URLTextInput, cmd = m.URLTextInput.Update(msg)
	return m, cmd
}

func (m Model) handleURLOverlay(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter":
		m.URLOverlayOpen = false
		m.URLTextInput.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.URLTextInput, cmd = m.URLTextInput.Update(msg)
	return m, cmd
}
