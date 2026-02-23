package program

import (
	"image"
	"strings"

	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
)

// ANSI SGR for primary (#FF8AD8) and dark box bg (#1a1a1a).
const (
	sgrPrimary = "\033[38;2;255;138;216m"
	sgrBoxBg   = "\033[48;2;26;26;26m"
	sgrReset   = "\033[0m"
)

// overlayLayer implements tea.Layer to draw a dark overlay with a centered input box.
type overlayLayer struct {
	inputText string
	width     int
	height    int
}

var _ tea.Layer = (*overlayLayer)(nil)

func newOverlayLayer(inputText string, width, height int) *overlayLayer {
	return &overlayLayer{inputText: inputText, width: width, height: height}
}

func (o *overlayLayer) buildInputBox() (string, int, int) {
	const boxInnerWidth = 40
	inner := o.inputText + "|"
	if len(inner) > boxInnerWidth {
		inner = inner[len(inner)-boxInnerWidth:]
	}
	padded := inner + strings.Repeat(" ", boxInnerWidth-len(inner))
	top := sgrPrimary + "╭" + strings.Repeat("─", boxInnerWidth) + "╮" + sgrReset
	mid := sgrPrimary + "│" + sgrBoxBg + sgrPrimary + " " + padded + " " + sgrReset + sgrPrimary + "│" + sgrReset
	bot := sgrPrimary + "╰" + strings.Repeat("─", boxInnerWidth) + "╯" + sgrReset
	boxStr := top + "\n" + mid + "\n" + bot
	boxW := boxInnerWidth + 2
	boxH := 3
	return boxStr, boxW, boxH
}

// Draw renders the overlay: dark backdrop then centered input box.
func (o *overlayLayer) Draw(s tea.Screen, r tea.Rectangle) {
	// Dark backdrop
	backdrop := &uv.Cell{
		Content: " ",
		Width:   1,
		Style:   uv.Style{Bg: Background},
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			s.SetCell(x, y, backdrop)
		}
	}

	boxStr, boxW, boxH := o.buildInputBox()
	centerX := r.Min.X + (r.Dx()-boxW)/2
	centerY := r.Min.Y + (r.Dy()-boxH)/2
	if centerX < r.Min.X {
		centerX = r.Min.X
	}
	if centerY < r.Min.Y {
		centerY = r.Min.Y
	}

	boxArea := image.Rect(centerX, centerY, centerX+boxW, centerY+boxH)
	uv.NewStyledString(boxStr).Draw(s, boxArea)
}
