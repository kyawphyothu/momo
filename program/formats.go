package program

import (
	"strings"

	"charm.land/bubbles/v2/table"
	"github.com/mattn/go-runewidth"
)

// parseFormatsTable parses yt-dlp's formats_table output into table rows.
// Actual format has 2 pipes (3 sections) - see program/stdout.txt for reference:
//
//	ID  EXT   RESOLUTION FPS CH |   FILESIZE    TBR PROTO | VCODEC  VBR ACODEC  ABR ASR MORE INFO
func parseFormatsTable(raw string) []table.Row {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) < 3 {
		return nil
	}

	var rows []table.Row
	// Skip header (line 0) and separator (line 1)
	for i := 2; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}

		row := parseFormatsLine(line)
		if row != nil {
			rows = append(rows, row)
		}
	}

	return rows
}

// parseFormatsLine parses a single line split by | into 3 or 4 sections.
// Real format (stdout.txt) has 2 pipes (3 parts); section 3 contains codec + MORE INFO combined.
// Some variants have 3 pipes (4 parts) with MORE INFO in part 4.
func parseFormatsLine(line string) table.Row {
	parts := splitByPipe(line)
	if len(parts) < 3 {
		return nil
	}

	id, ext, resolution, fps, ch := parseSection1(parts[0])
	filesize, tbr, proto := parseSection2(parts[1])

	var vcodec, vbr, acodec, abr, asr, moreInfo string
	if len(parts) >= 4 {
		vcodec, vbr, acodec, abr, asr = parseSection3(parts[2])
		moreInfo = strings.TrimSpace(parts[3])
	} else {
		vcodec, vbr, acodec, abr, asr, moreInfo = parseSection3WithMoreInfo(parts[2])
	}

	// Column widths must match model.go. Right-align numeric columns.
	return table.Row{
		id, ext, resolution, fps, ch,
		filesize, tbr, proto,
		vcodec, vbr, acodec, abr, asr,
		moreInfo,
	}
}

// formatFormatsRowsForDisplay applies right-alignment to numeric columns.
// Call this before SetRows when displaying the table.
func formatFormatsRowsForDisplay(rows []table.Row) []table.Row {
	// Column widths and indices must match model.go
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

func splitByPipe(s string) []string {
	// Support both Unicode box-drawing pipe (U+2502) and ASCII pipe
	s = strings.ReplaceAll(s, "│", "|")
	return strings.Split(s, "|")
}

// parseSection1: ID EXT RESOLUTION FPS CH
// e.g. "sb3 mhtml 48x27        0    " or "139 m4a   audio only      2 "
func parseSection1(s string) (id, ext, resolution, fps, ch string) {
	fields := strings.Fields(s)
	if len(fields) < 2 {
		return "", "", "", "", ""
	}
	id = fields[0]
	ext = fields[1]

	// Find last 1-2 numeric tokens for FPS and CH
	var numIndices []int
	for i := len(fields) - 1; i >= 2; i-- {
		if isNumeric(fields[i]) {
			numIndices = append([]int{i}, numIndices...)
			if len(numIndices) >= 2 {
				break
			}
		}
	}

	if len(numIndices) == 2 {
		fps = fields[numIndices[0]]
		ch = fields[numIndices[1]]
		resolution = strings.Join(fields[2:numIndices[0]], " ")
	} else if len(numIndices) == 1 {
		idx := numIndices[0]
		resolution = strings.Join(fields[2:idx], " ")
		// For "audio only", the trailing number is CH; for "NxN" it's FPS
		if strings.Contains(strings.ToLower(resolution), "audio") {
			ch = fields[idx]
		} else {
			fps = fields[idx]
		}
	} else {
		resolution = strings.Join(fields[2:], " ")
	}
	return
}

// parseSection2: FILESIZE TBR PROTO
func parseSection2(s string) (filesize, tbr, proto string) {
	fields := strings.Fields(s)

	// Merge standalone approximate prefix with the next token.
	// yt-dlp uses "~" (tilde) or "≈" (U+2248) before estimated sizes.
	if len(fields) >= 2 && (fields[0] == "~" || fields[0] == "≈") {
		fields = append([]string{fields[0] + " " + fields[1]}, fields[2:]...)
	}

	if len(fields) >= 3 {
		return fields[0], fields[1], fields[2]
	}
	if len(fields) == 2 {
		return fields[0], fields[1], ""
	}
	if len(fields) == 1 {
		if looksLikeFilesize(fields[0]) {
			return fields[0], "", ""
		}
		return "", "", fields[0]
	}
	return "", "", ""
}

func looksLikeFilesize(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] == '~' || (s[0] >= '0' && s[0] <= '9')
}

// parseSection3WithMoreInfo parses section 3 which contains VCODEC, VBR, ACODEC, ABR, ASR, and MORE INFO.
// MORE INFO is extracted from the end (e.g. "storyboard", "low, m4a_dash", "144p, mp4_dash", "[en]").
func parseSection3WithMoreInfo(s string) (vcodec, vbr, acodec, abr, asr, moreInfo string) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return "", "", "", "", "", ""
	}

	// Extract MORE INFO from the end: tokens that contain comma, brackets, or don't look like codec values
	var moreInfoFields []string
	i := len(fields) - 1
	for i >= 0 {
		f := fields[i]
		if looksLikeMoreInfo(f) {
			moreInfoFields = append([]string{f}, moreInfoFields...)
			i--
		} else {
			break
		}
	}
	moreInfo = strings.Join(moreInfoFields, " ")
	codecFields := fields[:i+1]

	vc, vbr, ac, ab, as := parseSection3(strings.Join(codecFields, " "))
	return vc, vbr, ac, ab, as, moreInfo
}

func looksLikeMoreInfo(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Contains comma: "low, m4a_dash", "144p, mp4_dash"
	if strings.Contains(s, ",") {
		return true
	}
	// Brackets: [en]
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return true
	}
	// Resolution-like: 144p, 240p, 1080p
	if len(s) > 1 && s[len(s)-1] == 'p' && isNumeric(s[:len(s)-1]) {
		return true
	}
	// Format tags: m4a_dash, webm_dash, mp4_dash
	if strings.Contains(s, "_") && !strings.Contains(s, ".") {
		return true
	}
	// Single word: storyboard (images is VCODEC, not MORE INFO)
	if s == "storyboard" {
		return true
	}
	return false
}

// mergeCompoundTokens joins adjacent "audio"/"video" + "only" pairs into single tokens.
func mergeCompoundTokens(fields []string) []string {
	var result []string
	for i := 0; i < len(fields); i++ {
		if i+1 < len(fields) && (fields[i] == "audio" || fields[i] == "video") && fields[i+1] == "only" {
			result = append(result, fields[i]+" "+fields[i+1])
			i++
		} else {
			result = append(result, fields[i])
		}
	}
	return result
}

// parseSection3: VCODEC VBR ACODEC ABR ASR
// VCODEC or ACODEC can be multi-word ("audio only", "video only", "images")
func parseSection3(s string) (vcodec, vbr, acodec, abr, asr string) {
	fields := mergeCompoundTokens(strings.Fields(s))
	if len(fields) == 0 {
		return "", "", "", "", ""
	}

	// "audio only" / "video only" as VCODEC — VBR is empty (no video bitrate).
	// Only applies when remaining tokens fit acodec/abr/asr (<=3 tokens).
	// With more tokens (e.g. 4-pipe format including language tags), fall through
	// to the standard positional handler.
	if fields[0] == "audio only" || fields[0] == "video only" {
		rest := fields[1:]
		if len(rest) <= 3 {
			vcodec = fields[0]
			switch len(rest) {
			case 3:
				return vcodec, "", rest[0], rest[1], rest[2]
			case 2:
				return vcodec, "", rest[0], rest[1], ""
			case 1:
				return vcodec, "", rest[0], "", ""
			default:
				return vcodec, "", "", "", ""
			}
		}
	}

	switch len(fields) {
	case 1:
		return fields[0], "", "", "", ""
	case 2:
		return fields[0], fields[1], "", "", ""
	case 3:
		return fields[0], fields[1], fields[2], "", ""
	case 4:
		return fields[0], fields[1], fields[2], fields[3], ""
	default:
		vcodec = strings.Join(fields[:len(fields)-4], " ")
		vbr = fields[len(fields)-4]
		acodec = fields[len(fields)-3]
		abr = fields[len(fields)-2]
		asr = fields[len(fields)-1]
		return
	}
}

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}
