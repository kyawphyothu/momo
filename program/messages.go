package program

import "github.com/lrstanley/go-ytdlp"

type FormatsLoadedMsg struct {
	Formats []*ytdlp.ExtractedFormat
	Err     error
}

// YtdlpInstallDoneMsg is sent when yt-dlp install completes (success or error).
type YtdlpInstallDoneMsg struct {
	Err error
}
