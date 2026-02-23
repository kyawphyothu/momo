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

	case tea.KeyMsg:
		if m.OverlayOpen {
			switch msg.String() {
			case "esc", "enter":
				m.OverlayOpen = false
				return m, nil
			case "backspace":
				if len(m.InputText) > 0 {
					m.InputText = m.InputText[:len(m.InputText)-1]
				}
				return m, nil
			default:
				if k := msg.Key(); k.Text != "" {
					m.InputText += k.Text
				}
				return m, nil
			}
		}

		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "l":
			m.OverlayOpen = true
			return m, nil
		default:
			return m, nil
		}
	}

	return m, nil
}
