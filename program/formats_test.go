package program

import (
	"os"
	"reflect"
	"testing"

	"charm.land/bubbles/v2/table"
)

func TestParseFormatsTable(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []table.Row
	}{
		{
			name:     "empty string returns nil",
			input:    "",
			expected: nil,
		},
		{
			name:     "single line returns nil",
			input:    "ID  EXT  RES",
			expected: nil,
		},
		{
			name: "only header and separator returns nil",
			input: `ID  EXT   RESOLUTION FPS CH |   FILESIZE    TBR PROTO | VCODEC  VBR ACODEC  ABR ASR | MORE INFO
---`,
			expected: nil,
		},
		{
			name: "parses audio format with CH",
			input: `ID  EXT   RESOLUTION FPS CH |   FILESIZE    TBR PROTO | VCODEC           VBR ACODEC      ABR ASR | MORE INFO
---
139 m4a   audio only      2 |    1.24MiB    49k https | audio only           mp4a.40.5   49k 22k [en] | low, m4a_dash`,
			expected: []table.Row{
				{"139", "m4a", "audio only", "", "2", "1.24MiB", "49k", "https", "audio only", "mp4a.40.5", "49k", "22k", "[en]", "low, m4a_dash"},
			},
		},
		{
			name: "parses multiple rows",
			input: `ID  EXT   RESOLUTION FPS CH |   FILESIZE    TBR PROTO | VCODEC  VBR ACODEC  ABR ASR | MORE INFO
---
139 m4a   audio only      2 |    1.24MiB    49k https | audio only  mp4a.40.5 49k 22k | low
140 m4a   audio only      2 |    3.29MiB   130k https | audio only  mp4a.40.2 130k 44k | medium`,
			expected: []table.Row{
				{"139", "m4a", "audio only", "", "2", "1.24MiB", "49k", "https", "audio only", "", "mp4a.40.5", "49k", "22k", "low"},
				{"140", "m4a", "audio only", "", "2", "3.29MiB", "130k", "https", "audio only", "", "mp4a.40.2", "130k", "44k", "medium"},
			},
		},
		{
			name: "skips empty lines",
			input: `ID  EXT   RESOLUTION FPS CH |   FILESIZE    TBR PROTO | VCODEC  VBR ACODEC  ABR ASR | MORE INFO
---

139 m4a   audio only      2 |    1.24MiB    49k https | audio only  mp4a.40.5 49k 22k | low

140 m4a   audio only      2 |    3.29MiB   130k https | audio only  mp4a.40.2 130k 44k | medium`,
			expected: []table.Row{
				{"139", "m4a", "audio only", "", "2", "1.24MiB", "49k", "https", "audio only", "", "mp4a.40.5", "49k", "22k", "low"},
				{"140", "m4a", "audio only", "", "2", "3.29MiB", "130k", "https", "audio only", "", "mp4a.40.2", "130k", "44k", "medium"},
			},
		},
		{
			name: "parses video format with FPS",
			input: `ID  EXT   RESOLUTION FPS CH |   FILESIZE    TBR PROTO | VCODEC  VBR ACODEC  ABR ASR | MORE INFO
---
137 mp4  1920x1080       30    |    ~5MiB    128k https | avc1.64001f 128k mp4a.40.2 128k 44k | 1080p`,
			expected: []table.Row{
				{"137", "mp4", "1920x1080", "30", "", "~5MiB", "128k", "https", "avc1.64001f", "128k", "mp4a.40.2", "128k", "44k", "1080p"},
			},
		},
		{
			name: "Unicode pipe works",
			input: `ID  EXT   RESOLUTION FPS CH │   FILESIZE    TBR PROTO │ VCODEC  VBR ACODEC  ABR ASR │ MORE INFO
---
139 m4a   audio only      2 │    1.24MiB    49k https │ audio only  mp4a.40.5 49k 22k │ low`,
			expected: []table.Row{
				{"139", "m4a", "audio only", "", "2", "1.24MiB", "49k", "https", "audio only", "", "mp4a.40.5", "49k", "22k", "low"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFormatsTable(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("parseFormatsTable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseSection1(t *testing.T) {
	tests := []struct {
		name                                      string
		input                                     string
		wantID, wantExt, wantRes, wantFPS, wantCH string
	}{
		{"video with FPS", "sb3 mhtml 48x27        0    ", "sb3", "mhtml", "48x27", "0", ""},
		{"audio with CH", "139 m4a   audio only      2 ", "139", "m4a", "audio only", "", "2"},
		{"video 1920x1080 30fps", "137 mp4 1920x1080 30", "137", "mp4", "1920x1080", "30", ""},
		{"too few fields", "id", "", "", "", "", ""},
		{"resolution only", "sb0 mhtml 320x180 1", "sb0", "mhtml", "320x180", "1", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, ext, res, fps, ch := parseSection1(tt.input)
			if id != tt.wantID || ext != tt.wantExt || res != tt.wantRes || fps != tt.wantFPS || ch != tt.wantCH {
				t.Errorf("parseSection1(%q) = (%q, %q, %q, %q, %q), want (%q, %q, %q, %q, %q)",
					tt.input, id, ext, res, fps, ch, tt.wantID, tt.wantExt, tt.wantRes, tt.wantFPS, tt.wantCH)
			}
		})
	}
}

func TestParseSection2(t *testing.T) {
	tests := []struct {
		name                             string
		input                            string
		wantFilesize, wantTBR, wantProto string
	}{
		{"three fields", "   1.24MiB    49k https ", "1.24MiB", "49k", "https"},
		{"empty section", "                   mhtml ", "", "", "mhtml"},
		{"two fields", "  ~5MiB 128k", "~5MiB", "128k", ""},
		{"one field", "  mhtml", "", "", "mhtml"},
		{"empty", "", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, tbr, proto := parseSection2(tt.input)
			if fs != tt.wantFilesize || tbr != tt.wantTBR || proto != tt.wantProto {
				t.Errorf("parseSection2(%q) = (%q, %q, %q), want (%q, %q, %q)",
					tt.input, fs, tbr, proto, tt.wantFilesize, tt.wantTBR, tt.wantProto)
			}
		})
	}
}

func TestParseSection3(t *testing.T) {
	tests := []struct {
		name                                              string
		input                                             string
		wantVcodec, wantVBR, wantAcodec, wantABR, wantASR string
	}{
		{"single word", " images ", "images", "", "", "", ""},
		{"audio only with codecs", " audio only           mp4a.40.5   49k 22k ", "audio only", "", "mp4a.40.5", "49k", "22k"},
		{"five fields", " avc1.64001f 128k mp4a.40.2 128k 44k ", "avc1.64001f", "128k", "mp4a.40.2", "128k", "44k"},
		{"six fields (vcodec multi-word)", " audio only mp4a.40.5 49k 22k [en] ", "audio only", "mp4a.40.5", "49k", "22k", "[en]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc, vbr, ac, abr, asr := parseSection3(tt.input)
			if vc != tt.wantVcodec || vbr != tt.wantVBR || ac != tt.wantAcodec || abr != tt.wantABR || asr != tt.wantASR {
				t.Errorf("parseSection3(%q) = (%q, %q, %q, %q, %q), want (%q, %q, %q, %q, %q)",
					tt.input, vc, vbr, ac, abr, asr, tt.wantVcodec, tt.wantVBR, tt.wantAcodec, tt.wantABR, tt.wantASR)
			}
		})
	}
}

func TestSplitByPipe(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"Unicode pipe", "a│b│c", []string{"a", "b", "c"}},
		{"ASCII pipe", "a|b|c", []string{"a", "b", "c"}},
		{"mixed", "a│b|c", []string{"a", "b", "c"}},
		{"no pipe", "a b c", []string{"a b c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitByPipe(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("splitByPipe(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestParseFormatsTableFromStdout verifies parsing of actual yt-dlp output (stdout.txt).
func TestParseFormatsTableFromStdout(t *testing.T) {
	raw, err := os.ReadFile("stdout.txt")
	if err != nil {
		t.Skip("stdout.txt not found, skipping")
	}
	rows := parseFormatsTable(string(raw))
	if len(rows) == 0 {
		t.Error("expected rows from stdout.txt, got none")
	}
	// Verify first row (storyboard)
	if len(rows[0]) < 14 {
		t.Errorf("row 0 has %d columns, want 14", len(rows[0]))
	}
	if rows[0][0] != "sb3" || rows[0][2] != "48x27" || rows[0][5] != "" || rows[0][13] != "storyboard" {
		t.Errorf("row 0: got ID=%q RES=%q MORE=%q", rows[0][0], rows[0][2], rows[0][13])
	}
	if rows[4][12] != "22k" || rows[4][13] != "[en] low, m4a_dash" {
		t.Errorf("row 4: got ASR=%q", rows[4][12])
	}
	if rows[9][10] != "video only" {
		t.Errorf("row 9: got ACODEC=%q", rows[9][10])
	}
	if rows[18][5] != "≈ 11.28MiB" {
		t.Errorf("row 18: got Filesize=%q", rows[18][5])
	}
	// Find and verify 139 audio row (comes after sb3, sb2, sb1, sb0)
	var found139 bool
	for _, row := range rows {
		if len(row) >= 14 && row[0] == "139" {
			if row[2] != "audio only" || row[13] != "[en] low, m4a_dash" {
				t.Errorf("row 139: got RES=%q MORE=%q", row[2], row[13])
			}
			found139 = true
			break
		}
	}
	if !found139 {
		t.Error("expected to find format 139 in parsed rows")
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"0", true},
		{"123", true},
		{"2", true},
		{"", false},
		{"12k", false},
		{"1.24", false},
		{"abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isNumeric(tt.input)
			if got != tt.expected {
				t.Errorf("isNumeric(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
