package program

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	URLOverlayOpen bool
	URLTextInput   textinput.Model

	Width  int
	Height int
}

func InitialModel() Model {
	urlTextInput := textinput.New()
	urlTextInput.Placeholder = "Enter URL"
	urlTextInput.CharLimit = 100
	urlTextInput.SetWidth(URLOverlayWidth - 7) // 2 border, 2 padding, 1 prompt, 1 gap, 1 right cursor visible

	return Model{Width: 80, Height: 24, URLTextInput: urlTextInput}
}

func (m Model) Init() tea.Cmd {
	return nil
}
