package program

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case YtdlpInstallDoneMsg:
		m.IsRunningYtdlpInstall = false
		m.YtdlpInstallErr = msg.Err
		return m, nil

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.FormatsViewport.SetWidth(msg.Width)
		m.FormatsViewport.SetHeight(msg.Height)
		return m, nil

	case tea.PasteMsg:
		return m.handlePaste(msg)

	case FormatsLoadedMsg:
		if msg.Err != nil {
			// handle error (e.g. store in model, show in view)
			return m, nil
		}
		m.FormatsLoaded = true
		m.FormatsLoading = false
		m.FormatsViewport.SetWidth(m.Width)
		m.FormatsViewport.SetHeight(m.Height)
		m.FormatsViewport.SetContent(msg.Table)
		return m, nil

	case tea.KeyMsg:
		if m.URLOverlayOpen {
			return m.handleURLOverlay(msg)
		}
		if m.FormatsLoaded {
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "l":
				m.URLOverlayOpen = true
				m.URLTextInput.Focus()
				return m, nil
			}

			var cmd tea.Cmd
			m.FormatsViewport, cmd = m.FormatsViewport.Update(msg)
			return m, cmd
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
	default:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	// return m, nil
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
	case "esc":
		m.URLOverlayOpen = false
		m.URLTextInput.Blur()
		return m, nil
	case "enter":
		m.URLOverlayOpen = false
		m.FormatsLoading = true
		m.URLTextInput.Blur()
		return m, FetchFormatsCmd(m.URLTextInput.Value())
	}

	var cmd tea.Cmd
	m.URLTextInput, cmd = m.URLTextInput.Update(msg)
	return m, cmd
}
