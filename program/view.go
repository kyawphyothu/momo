package program

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	if !IsWindowSizeValid(m.Width, m.Height) {
		return newView(renderSmallWindowScreen(m.Width, m.Height))
	}

	if m.IsRunningYtdlpInstall {
		return newView(lipgloss.NewStyle().Foreground(Primary).Render("Checking for yt-dlp..."))
	}

	if m.YtdlpInstallErr != nil {
		errMsg := lipgloss.NewStyle().Foreground(Error).Render("yt-dlp install failed: " + m.YtdlpInstallErr.Error())
		return newView(errMsg)
	}

	var mainContent string
	if m.FormatsLoading {
		mainContent = m.Spinner.View() + " " + lipgloss.NewStyle().Foreground(Primary).Render("Dancing Machine Momo on standby...")
		mainContent = lipgloss.NewStyle().Width(m.Width).Height(m.Height).Align(lipgloss.Center, lipgloss.Center).Render(mainContent)
	} else if m.FormatsLoaded {
		mainContent = m.FormatsTable.View()
	} else if m.FormatsErr != nil {
		mainContent = lipgloss.NewStyle().Foreground(Error).Render("Formats loading failed: " + m.FormatsErr.Error())
	} else {
		mainContent = lipgloss.NewStyle().Foreground(Primary).Render("Hello, World!")
	}

	if m.URLOverlayOpen {
		overlayContent := renderURLOverlay(m.URLTextInput)
		overlayLayer := lipgloss.
			NewLayer(overlayContent).
			X((m.Width - URLOverlayWidth) / 2).
			Y(URLOverlayYAxis).
			Z(1)

		mainLayer := lipgloss.
			NewLayer(mainContent).
			X(0).
			Y(0).
			Z(0)

		canvas := lipgloss.NewCanvas(mainLayer, overlayLayer)
		finalContent := canvas.Render()

		return newView(finalContent)
	}

	return newView(mainContent)
}

func newView(content string) tea.View {
	v := tea.NewView(content)
	v.AltScreen = true
	v.BackgroundColor = Background
	v.ForegroundColor = Text
	return v
}

func renderURLOverlay(textInput textinput.Model) string {
	b := NewBoxWithLabel(URLOverlayBoxStyle, URLOverlayLabelStyle)
	return b.Render("YouTubeURL", textInput.View(), URLOverlayWidth)
}

func renderDimension(label string, current, minimum int) string {
	muted := lipgloss.NewStyle().Foreground(Muted)
	color := Success
	if current < minimum {
		color = Error
	}
	value := lipgloss.NewStyle().Foreground(color).Render(strconv.Itoa(current))
	return muted.Render(label+": ") + value + muted.Render(fmt.Sprintf("/%d", minimum))
}

func renderSmallWindowScreen(width, height int) string {
	var content strings.Builder

	content.WriteString(lipgloss.NewStyle().Foreground(Primary).Bold(true).Render("Window size is too small"))
	content.WriteString("\n")
	content.WriteString(lipgloss.NewStyle().Foreground(Muted).Bold(false).Render("Please resize your window to continue"))
	content.WriteString("\n")
	content.WriteString(renderDimension("Width", width, minWindowWidth))
	content.WriteString("\n")
	content.WriteString(renderDimension("Height", height, minWindowHeight))
	content.WriteString("\n")

	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render(content.String())
}
