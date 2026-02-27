package program

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/kyawphyothu/momo/table"
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
	FormatsErr     error
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
		{Title: "ID", Width: FormatsTableIDWidth},
		{Title: "EXT", Width: FormatsTableEXTWidth},
		{Title: "RES", Width: FormatsTableRESWidth},
		{Title: rightAlign("FPS", FormatsTableFPSWidth), Width: FormatsTableFPSWidth},
		{Title: rightAlign("CH", FormatsTableCHWidth), Width: FormatsTableCHWidth},
		{Title: rightAlign("SIZE", FormatsTableSIZEWidth), Width: FormatsTableSIZEWidth},
		{Title: rightAlign("TBR", FormatsTableTBRWidth), Width: FormatsTableTBRWidth},
		{Title: "PROTO", Width: FormatsTablePROTOWidth},
		{Title: "VCODEC", Width: FormatsTableVCODECWidth},
		{Title: rightAlign("VBR", FormatsTableVBRWidth), Width: FormatsTableVBRWidth},
		{Title: "ACODEC", Width: FormatsTableACODECWidth},
		{Title: rightAlign("ABR", FormatsTableABRWidth), Width: FormatsTableABRWidth},
		{Title: rightAlign("ASR", FormatsTableASRWidth), Width: FormatsTableASRWidth},
		{Title: "MORE INFO", Width: FormatsTableMOREINFOWidth},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(FormatsTableHeight),
		table.WithWidth(FormatsTableWidth),
	)
	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		Foreground(Primary).
		Background(Background).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(Text).
		BorderBottom(true).
		Bold(true)
	tableStyle.Selected = tableStyle.Selected.
		Foreground(Text).
		Background(Highlight).
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
		FormatsErr:            nil,
	}
}

func (m Model) Init() tea.Cmd {
	// TODO: install yt-dlp in the background
	// return InstallYtdlpCmd()
	return tea.Batch(func() tea.Msg {
		return YtdlpInstallDoneMsg{Err: nil}
	}, m.Spinner.Tick)
}
