package program

import (
	"charm.land/lipgloss/v2"
)

var (
	// Primary is the main brand color ( #FF8AD8).
	Primary = lipgloss.Color("#FF8AD8")
	// Text is the text color ( #FFFFFF).
	Text = lipgloss.Color("#FFFFFF")
	// Background is the dark background color ( #121212).
	Background = lipgloss.Color("#121212")
	// Muted is used for secondary or subdued text ( #6B6B6B).
	Muted = lipgloss.Color("#6B6B6B")
	// Border is a subtle border color ( #2A2A2A).
	Border = lipgloss.Color("#2A2A2A")
	// Error is the error color ( #FF0000).
	Error = lipgloss.Color("#FF0000")
	// Success is the success color ( #00FF00).
	Success = lipgloss.Color("#00FF00")
	// Highlight is the highlight color ( #2ED6A3).
	Highlight = lipgloss.Color("#2ED6A3")

	// URLOverlayBoxStyle defines the visual style of the URL overlay container.
	URLOverlayBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Padding(0, 1, 0, 1)

	// URLOverlayLabelStyle defines the visual style of the URL overlay label.
	URLOverlayLabelStyle = lipgloss.NewStyle().
				Foreground(Highlight).
				Bold(true).
				PaddingTop(0).
				PaddingBottom(0).
				PaddingLeft(1).
				PaddingRight(1)
)
