package program

import "image/color"

// Theme colors for the app: primary #FF8AD8, dark background, and accents.
var (
	// Primary is the main brand color (#FF8AD8).
	Primary = color.RGBA{R: 0xff, G: 0x8a, B: 0xd8, A: 0xff}
	// Text is the text color.
	Text = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	// Background is the dark background color.
	Background = color.RGBA{R: 0x12, G: 0x12, B: 0x12, A: 0xff}
	// Muted is used for secondary or subdued text.
	Muted = color.RGBA{R: 0x6b, G: 0x6b, B: 0x6b, A: 0xff}
	// Border is a subtle border color.
	Border = color.RGBA{R: 0x2a, G: 0x2a, B: 0x2a, A: 0xff}
	// Error is the error color.
	Error = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	// Success is the success color.
	Success = color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}
)
