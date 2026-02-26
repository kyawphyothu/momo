package program

type FormatsLoadedMsg struct {
	Table string // raw table text from yt-dlp
	Err   error
}

// YtdlpInstallDoneMsg is sent when yt-dlp install completes (success or error).
type YtdlpInstallDoneMsg struct {
	Err error
}
