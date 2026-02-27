package program

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/table"
	"github.com/lrstanley/go-ytdlp"
	"github.com/mattn/go-runewidth"
)

// formatsToRows converts structured yt-dlp format data into table rows.
func formatsToRows(formats []*ytdlp.ExtractedFormat) []table.Row {
	rows := make([]table.Row, 0, len(formats))
	for _, f := range formats {
		rows = append(rows, formatToRow(f))
	}
	return rows
}

// formatToRow maps a single ExtractedFormat to a table row.
// Column order must match the columns defined in model.go.
func formatToRow(f *ytdlp.ExtractedFormat) table.Row {
	return table.Row{
		derefStr(f.FormatID),
		derefStr(f.Extension),
		derefStr(f.Resolution),
		formatFPS(f.FPS),
		formatChannels(f.AudioChannels),
		formatFileSize(f.FileSize, f.FileSizeApprox),
		formatBitrate(f.TBR),
		derefStr(f.Protocol),
		formatVCodec(f),
		formatBitrate(f.VBR),
		formatACodec(f),
		formatBitrate(f.ABR),
		formatSampleRate(f.ASR),
		formatMoreInfo(f.FormatNote, f.Language),
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func formatFPS(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%.0f", *f)
}

func formatChannels(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%.0f", *f)
}

func formatFileSize(exact *int, approx *int) string {
	if exact != nil {
		return humanizeBytes(*exact, false)
	}
	if approx != nil {
		return humanizeBytes(*approx, true)
	}
	return ""
}

func humanizeBytes(b int, approx bool) string {
	prefix := ""
	if approx {
		prefix = "~ "
	}

	size := float64(b)
	switch {
	case size >= 1<<30:
		return fmt.Sprintf("%s%.2fGiB", prefix, size/(1<<30))
	case size >= 1<<20:
		return fmt.Sprintf("%s%.2fMiB", prefix, size/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%s%.2fKiB", prefix, size/(1<<10))
	default:
		return fmt.Sprintf("%s%dB", prefix, b)
	}
}

func formatBitrate(br *float64) string {
	if br == nil || *br == 0 {
		return ""
	}
	return fmt.Sprintf("%dk", int(*br))
}

func formatSampleRate(sr *float64) string {
	if sr == nil || *sr == 0 {
		return ""
	}
	return fmt.Sprintf("%dk", int(*sr/1000))
}

func formatVCodec(f *ytdlp.ExtractedFormat) string {
	if f.VCodec != nil && *f.VCodec != "" {
		return *f.VCodec
	}
	if f.ACodec != nil && *f.ACodec != "" {
		return "audio only"
	}
	return ""
}

func formatACodec(f *ytdlp.ExtractedFormat) string {
	if f.ACodec != nil && *f.ACodec != "" {
		return *f.ACodec
	}
	if f.VCodec != nil && *f.VCodec != "" {
		return "video only"
	}
	return ""
}

func formatMoreInfo(note, lang *string) string {
	var parts []string
	if lang != nil && *lang != "" {
		parts = append(parts, fmt.Sprintf("[%s]", *lang))
	}
	if note != nil && *note != "" {
		parts = append(parts, *note)
	}
	return strings.Join(parts, " ")
}

// formatFormatsRowsForDisplay applies right-alignment to numeric columns.
// Call this before SetRows when displaying the table.
func formatFormatsRowsForDisplay(rows []table.Row) []table.Row {
	rightAlignCols := []struct{ idx, width int }{
		{3, 4}, {4, 3}, {5, 10}, {6, 6}, {9, 4}, {11, 6}, {12, 6},
	}
	out := make([]table.Row, len(rows))
	for i, row := range rows {
		r := make(table.Row, len(row))
		copy(r, row)
		for _, c := range rightAlignCols {
			if c.idx < len(r) {
				r[c.idx] = rightAlign(r[c.idx], c.width)
			}
		}
		out[i] = r
	}
	return out
}

// rightAlign pads s with leading spaces so it appears right-aligned in a
// fixed-width column. Uses runewidth for correct display width.
func rightAlign(s string, width int) string {
	w := runewidth.StringWidth(s)
	if w >= width {
		return runewidth.Truncate(s, width, "…")
	}
	return strings.Repeat(" ", width-w) + s
}
