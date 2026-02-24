package program

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
)

type BoxWithLabel struct {
	BoxStyle   lipgloss.Style
	LabelStyle lipgloss.Style
}

func NewBoxWithLabel(BoxStyle lipgloss.Style, LabelStyle lipgloss.Style) BoxWithLabel {
	return BoxWithLabel{BoxStyle: BoxStyle, LabelStyle: LabelStyle}
}

func (b BoxWithLabel) Render(label, content string, width int) string {
	var (
		// Query the box style for some of its border properties so we can
		// essentially take the top border apart and put it around the label.
		border          lipgloss.Border             = b.BoxStyle.GetBorderStyle()
		topBorderStyler func(strs ...string) string = lipgloss.NewStyle().Foreground(b.BoxStyle.GetBorderTopForeground()).Render
		topLeft         string                      = topBorderStyler(border.TopLeft)
		topRight        string                      = topBorderStyler(border.TopRight)

		renderedLabel string = b.LabelStyle.Render(label)
	)

	// Render top row with the label
	cellsShort := max(0, width-lipgloss.Width(topLeft+topRight+renderedLabel))
	gap := strings.Repeat(border.Top, cellsShort)
	top := topLeft + renderedLabel + topBorderStyler(gap) + topRight

	// Render the rest of the box
	bottomStyle := b.BoxStyle
	bottom := bottomStyle.
		BorderTop(false).
		Width(width).
		Render(content)

	// Stack the pieces
	return top + "\n" + bottom
}
