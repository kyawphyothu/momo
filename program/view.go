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
		v := tea.NewView(renderSmallWindowScreen(m.Width, m.Height))
		v.AltScreen = true
		v.BackgroundColor = Background
		v.ForegroundColor = Text
		return v
	}

	mainContent := lipgloss.NewStyle().Foreground(Primary).Render("Hello, World!")

	if m.URLOverlayOpen {
		overlayContent := renderURLOverlay(m.URLTextInput)
		overlayLayer := lipgloss.
			NewLayer(overlayContent).
			X((m.Width - URLOverlayWidth) / 2).
			Y(URLOverlayYAxis).
			Z(1)

		mainLayer := lipgloss.NewLayer(mainContent).X(0).Y(0).Z(0)

		canvas := lipgloss.NewCanvas(mainLayer, overlayLayer)
		finalContent := canvas.Render()

		v := tea.NewView(finalContent)
		v.AltScreen = true
		v.BackgroundColor = Background
		v.ForegroundColor = Text
		return v
	}

	v := tea.NewView(mainContent)
	v.AltScreen = true
	v.BackgroundColor = Background
	v.ForegroundColor = Text
	return v
}

func renderURLOverlay(textInput textinput.Model) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(0, 1, 0, 1)
	labelStyle := lipgloss.NewStyle().
		Foreground(Highlight).
		Bold(true).
		PaddingTop(0).
		PaddingBottom(0).
		PaddingLeft(1).
		PaddingRight(1)
	b := NewDefaultBoxWithLabel(boxStyle, labelStyle)
	return b.Render("YouTubeURL", textInput.View(), URLOverlayWidth)

	// const body = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam imperdiet ultrices felis, et condimentum odio fringilla et. Etiam porttitor volutpat diam non gravida. Donec et condimentum nulla. Sed consectetur viverra fermentum. Pellentesque eget metus ut felis faucibus consectetur nec ac mi."

	// b := NewDefaultBoxWithLabel()
	// return b.Render("Title", body, 40)

	// lipgloss.NewStyle().BorderTitle("Hello")

	// return style.Render(content.String())
}

func renderSmallWindowScreen(width, height int) string {
	var content strings.Builder

	content.WriteString(lipgloss.NewStyle().Foreground(Primary).Bold(true).Render("Window size is too small"))
	content.WriteString("\n")
	content.WriteString(lipgloss.NewStyle().Foreground(Muted).Bold(false).Render("Please resize your window to continue"))
	content.WriteString("\n")
	currentWidth := strconv.Itoa(width)
	if width < minWindowWidth {
		currentWidth = lipgloss.NewStyle().Foreground(Error).Bold(false).Render(currentWidth)
	} else {
		currentWidth = lipgloss.NewStyle().Foreground(Success).Bold(false).Render(currentWidth)
	}
	muted := lipgloss.NewStyle().Foreground(Muted).Bold(false)
	content.WriteString(muted.Render("Width: "))
	content.WriteString(currentWidth)
	content.WriteString(muted.Render(fmt.Sprintf("/%d", minWindowWidth)))
	content.WriteString("\n")
	currentHeight := strconv.Itoa(height)
	if height < minWindowHeight {
		currentHeight = lipgloss.NewStyle().Foreground(Error).Bold(false).Render(currentHeight)
	} else {
		currentHeight = lipgloss.NewStyle().Foreground(Success).Bold(false).Render(currentHeight)
	}
	content.WriteString(muted.Render("Height: "))
	content.WriteString(currentHeight)
	content.WriteString(muted.Render(fmt.Sprintf("/%d", minWindowHeight)))
	content.WriteString("\n")

	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render(content.String())
}
