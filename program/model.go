package program

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type Model struct {
	Width  int
	Height int

	Spinner spinner.Model

	IsRunningYtdlpInstall bool
	YtdlpInstallErr       error

	URLOverlayOpen bool
	URLTextInput   textinput.Model

	FormatsViewport viewport.Model
	FormatsLoaded   bool
	FormatsLoading  bool
}

func InitialModel() Model {
	urlTextInput := textinput.New()
	urlTextInput.Placeholder = "Enter URL"
	urlTextInput.CharLimit = 100
	urlTextInput.SetWidth(URLOverlayWidth - 7) // 2 border, 2 padding, 1 prompt, 1 gap, 1 right cursor visible

	FormatsViewport := viewport.New(viewport.WithWidth(80), viewport.WithHeight(24))
	FormatsViewport.MouseWheelEnabled = true

	s := spinner.New()
	s.Spinner = spinner.Meter
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		Width:                 80,
		Height:                24,
		Spinner:               s,
		IsRunningYtdlpInstall: true,
		URLTextInput:          urlTextInput,
		FormatsViewport:       FormatsViewport,
		FormatsLoaded:         false,
		FormatsLoading:        false,
	}
}

func (m Model) Init() tea.Cmd {
	// TODO: install yt-dlp in the background
	// return InstallYtdlpCmd()
	return tea.Batch(func() tea.Msg {
		return YtdlpInstallDoneMsg{Err: nil}
	}, m.Spinner.Tick)
}
