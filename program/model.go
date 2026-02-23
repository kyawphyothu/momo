package program

import (
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	OverlayOpen bool
	InputText   string
	Width       int
	Height      int
}

func InitialModel() Model {
	return Model{Width: 80, Height: 24}
}

func (m Model) Init() tea.Cmd {
	return nil
}
