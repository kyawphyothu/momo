package program

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
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

	FormatsTable   table.Model
	FormatsLoaded  bool
	FormatsLoading bool
}

func InitialModel() Model {
	urlTextInput := textinput.New()
	urlTextInput.Placeholder = "Enter URL"
	urlTextInput.CharLimit = 100
	urlTextInput.SetWidth(URLOverlayWidth - 7) // 2 border, 2 padding, 1 prompt, 1 gap, 1 right cursor visible

	s := spinner.New()
	s.Spinner = spinner.Meter
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "EXT", Width: 6},
		{Title: "RES", Width: 12},
		{Title: rightAlign("FPS", 4), Width: 4},
		{Title: rightAlign("CH", 3), Width: 3},
		{Title: rightAlign("SIZE", 10), Width: 10},
		{Title: rightAlign("TBR", 6), Width: 6},
		{Title: "PROTO", Width: 8},
		{Title: "VCODEC", Width: 12},
		{Title: rightAlign("VBR", 6), Width: 6},
		{Title: "ACODEC", Width: 12},
		{Title: rightAlign("ABR", 6), Width: 6},
		{Title: rightAlign("ASR", 6), Width: 6},
		{Title: "MORE INFO", Width: 18},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(115),
	)
	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		Foreground(Primary).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(Border).
		BorderBottom(true).
		Bold(true)
	tableStyle.Selected = tableStyle.Selected.
		Foreground(Text).
		Background(Primary).
		Bold(true)
	t.SetStyles(tableStyle)

	return Model{
		Width:                 80,
		Height:                24,
		Spinner:               s,
		IsRunningYtdlpInstall: true,
		URLTextInput:          urlTextInput,
		FormatsTable:          t,
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
