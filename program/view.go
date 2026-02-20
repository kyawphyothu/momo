package program

import tea "charm.land/bubbletea/v2"

func (m Model) View() tea.View {
	v := tea.NewView("Hello, World!")
	v.AltScreen = true
	return v
}
