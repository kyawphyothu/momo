package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/kyawphyothu/momo/program"
	"github.com/lrstanley/go-ytdlp"
	"github.com/mattn/go-runewidth"
)

func main() {
	var err error

	// dl := ytdlp.New().
	// 	FormatSort("res,ext:mp4:m4a").
	// 	RecodeVideo("mp4").
	// 	Output("%(extractor)s - %(title)s.%(ext)s")

	// dl := ytdlp.New().Print("formats_table")

	// res, err := dl.Run(context.TODO(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("String:")
	// fmt.Println(res.String())
	// fmt.Println("Stdout:")
	// fmt.Println(res.Stdout)
	// fmt.Println("GetExtractedInfo:")
	// os.Exit(0)
	// ext, err := res.GetExtractedInfo()
	// if err != nil {
	// 	panic(err)
	// }
	// for _, format := range ext {
	// 	fmt.Println(format.Format)
	// 	fmt.Println(format.FormatID)
	// 	fmt.Println(format.FormatNote)
	// }
	// ytdlp.MustInstallAll(context.TODO())

	// dl := ytdlp.New().
	// 	PrintJSON().
	// 	NoProgress().
	// 	FormatSort("res,ext:mp4:m4a").
	// 	RecodeVideo("mp4").
	// 	NoPlaylist().
	// 	NoOverwrites().
	// 	Continue().
	// 	ProgressFunc(100*time.Millisecond, func(prog ytdlp.ProgressUpdate) {
	// 		fmt.Printf( //nolint:forbidigo
	// 			"%s @ %s [eta: %s] :: %s\n",
	// 			prog.Status,
	// 			prog.PercentString(),
	// 			prog.ETA(),
	// 			prog.Filename,
	// 		)
	// 	}).
	// 	Output("%(extractor)s - %(title)s.%(ext)s")

	// r, err := dl.Run(context.TODO(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	// if err != nil {
	// 	panic(err)
	// }

	// f, err := os.Create("results.json")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	// enc := json.NewEncoder(f)
	// enc.SetIndent("", "    ")

	// if err = enc.Encode(r); err != nil {
	// 	panic(err)
	// }

	// slog.Info("wrote results to results.json")

	// os.Exit(0)

	m := program.InitialModel()
	p := tea.NewProgram(m)
	_, err = p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	os.Exit(0)

	ctx := context.Background()
	dl := ytdlp.New().DumpJSON().NoWarnings()
	res, err := dl.Run(ctx, "https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	if err != nil {
		panic(err)
	}
	// fmt.Println(res.Stdout)
	// fmt.Println(res.Args)
	info, err := res.GetExtractedInfo()
	if err != nil {
		panic(err)
	}
	if len(info) == 0 {
		panic(fmt.Errorf("no info returned by yt-dlp"))
	}
	if len(info[0].Formats) == 0 {
		panic(fmt.Errorf("no formats found"))
	}
	fmt.Println("Title:")
	fmt.Println(*info[0].Title)
	fmt.Println("ABR:")
	fmt.Println(*info[0].ABR)
	fmt.Println("Formats:")
	formatRows := formatsToRows(info[0].Formats)
	for _, row := range formatRows {
		for _, cell := range row {
			fmt.Println(cell)
		}
	}

	fmt.Println("Formatted Rows:")
	formattedRows := formatFormatsRowsForDisplay(formatRows)
	for _, row := range formattedRows {
		for _, cell := range row {
			fmt.Println(cell)
		}
	}

	os.Exit(0)

}

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
