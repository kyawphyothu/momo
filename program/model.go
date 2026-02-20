package program

import (
	tea "charm.land/bubbletea/v2"
)

type Model struct {
}

func InitialModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}
