package program

import (
	"testing"

	"github.com/lrstanley/go-ytdlp"
)

func assertEqual(t *testing.T, field, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", field, got, want)
	}
}

func TestFormatsToRows_AudioOnly(t *testing.T) {
	formats := []*ytdlp.ExtractedFormat{
		{
			FormatID:      ptr("139"),
			Extension:     ptr("m4a"),
			Resolution:    ptr("audio only"),
			AudioChannels: ptr(2.0),
			FileSize:      ptr(1300234),
			TBR:           ptr(49.0),
			Protocol:      ptr("https"),
			ACodec:        ptr("mp4a.40.5"),
			ABR:           ptr(49.0),
			ASR:           ptr(22050.0),
			FormatNote:    ptr("low"),
			Language:      ptr("en"),
		},
	}

	rows := formatsToRows(formats)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	row := rows[0]
	assertEqual(t, "ID", row[0], "139")
	assertEqual(t, "EXT", row[1], "m4a")
	assertEqual(t, "RES", row[2], "audio only")
	assertEqual(t, "FPS", row[3], "")
	assertEqual(t, "CH", row[4], "2")
	assertEqual(t, "SIZE", row[5], "1.24MiB")
	assertEqual(t, "TBR", row[6], "49k")
	assertEqual(t, "PROTO", row[7], "https")
	assertEqual(t, "VCODEC", row[8], "audio only")
	assertEqual(t, "VBR", row[9], "")
	assertEqual(t, "ACODEC", row[10], "mp4a.40.5")
	assertEqual(t, "ABR", row[11], "49k")
	assertEqual(t, "ASR", row[12], "22k")
	assertEqual(t, "MORE INFO", row[13], "[en] low")
}

func TestFormatsToRows_VideoOnly(t *testing.T) {
	formats := []*ytdlp.ExtractedFormat{
		{
			FormatID:   ptr("137"),
			Extension:  ptr("mp4"),
			Resolution: ptr("1920x1080"),
			FPS:        ptr(25.0),
			FileSize:   ptr(80899727),
			TBR:        ptr(3038.0),
			Protocol:   ptr("https"),
			VCodec:     ptr("avc1.640028"),
			VBR:        ptr(3038.0),
			FormatNote: ptr("1080p"),
		},
	}

	rows := formatsToRows(formats)
	row := rows[0]
	assertEqual(t, "VCODEC", row[8], "avc1.640028")
	assertEqual(t, "ACODEC", row[10], "video only")
	assertEqual(t, "FPS", row[3], "25")
	assertEqual(t, "SIZE", row[5], "77.15MiB")
	assertEqual(t, "MORE INFO", row[13], "1080p")
}

func TestFormatsToRows_Combined(t *testing.T) {
	formats := []*ytdlp.ExtractedFormat{
		{
			FormatID:       ptr("18"),
			Extension:      ptr("mp4"),
			Resolution:     ptr("640x360"),
			FPS:            ptr(25.0),
			AudioChannels:  ptr(2.0),
			FileSizeApprox: ptr(11826381),
			TBR:            ptr(444.0),
			Protocol:       ptr("https"),
			VCodec:         ptr("avc1.42001E"),
			ACodec:         ptr("mp4a.40.2"),
			ASR:            ptr(44100.0),
			FormatNote:     ptr("360p"),
			Language:       ptr("en"),
		},
	}

	rows := formatsToRows(formats)
	row := rows[0]
	assertEqual(t, "VCODEC", row[8], "avc1.42001E")
	assertEqual(t, "ACODEC", row[10], "mp4a.40.2")
	assertEqual(t, "SIZE", row[5], "~ 11.28MiB")
	assertEqual(t, "ASR", row[12], "44k")
	assertEqual(t, "MORE INFO", row[13], "[en] 360p")
}

func TestFormatsToRows_Storyboard(t *testing.T) {
	formats := []*ytdlp.ExtractedFormat{
		{
			FormatID:   ptr("sb3"),
			Extension:  ptr("mhtml"),
			Resolution: ptr("48x27"),
			FPS:        ptr(0.0),
			Protocol:   ptr("mhtml"),
			FormatNote: ptr("storyboard"),
		},
	}

	rows := formatsToRows(formats)
	row := rows[0]
	assertEqual(t, "FPS", row[3], "0")
	assertEqual(t, "VCODEC", row[8], "")
	assertEqual(t, "ACODEC", row[10], "")
	assertEqual(t, "MORE INFO", row[13], "storyboard")
}

func TestFormatsToRows_Empty(t *testing.T) {
	rows := formatsToRows(nil)
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestHumanizeBytes(t *testing.T) {
	tests := []struct {
		bytes  int
		approx bool
		want   string
	}{
		{1300234, false, "1.24MiB"},
		{11826381, true, "~ 11.28MiB"},
		{80899727, false, "77.15MiB"},
		{1073741824, false, "1.00GiB"},
		{512000, false, "500.00KiB"},
		{500, false, "500B"},
		{0, false, "0B"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := humanizeBytes(tt.bytes, tt.approx)
			if got != tt.want {
				t.Errorf("humanizeBytes(%d, %v) = %q, want %q", tt.bytes, tt.approx, got, tt.want)
			}
		})
	}
}

func TestFormatBitrate(t *testing.T) {
	assertEqual(t, "nil", formatBitrate(nil), "")
	assertEqual(t, "zero", formatBitrate(ptr(0.0)), "")
	assertEqual(t, "49", formatBitrate(ptr(49.0)), "49k")
	assertEqual(t, "3038", formatBitrate(ptr(3038.0)), "3038k")
}

func TestFormatSampleRate(t *testing.T) {
	assertEqual(t, "nil", formatSampleRate(nil), "")
	assertEqual(t, "zero", formatSampleRate(ptr(0.0)), "")
	assertEqual(t, "22050", formatSampleRate(ptr(22050.0)), "22k")
	assertEqual(t, "44100", formatSampleRate(ptr(44100.0)), "44k")
	assertEqual(t, "48000", formatSampleRate(ptr(48000.0)), "48k")
}

func TestFormatMoreInfo(t *testing.T) {
	assertEqual(t, "both", formatMoreInfo(ptr("low"), ptr("en")), "[en] low")
	assertEqual(t, "note only", formatMoreInfo(ptr("1080p"), nil), "1080p")
	assertEqual(t, "lang only", formatMoreInfo(nil, ptr("en")), "[en]")
	assertEqual(t, "neither", formatMoreInfo(nil, nil), "")
}
