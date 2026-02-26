package program

import (
	"context"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/lrstanley/go-ytdlp"
)

// InstallYtdlpCmd runs ytdlp.Install in the background and sends YtdlpInstallDoneMsg when finished.
func InstallYtdlpCmd() tea.Cmd {
	return func() tea.Msg {
		_, err := ytdlp.Install(context.TODO(), nil)
		return YtdlpInstallDoneMsg{Err: err}
	}
}

// fakeFetchFormats enables a 5s fake fetch for testing the loading UI.
const fakeFetchFormats = true

func FetchFormatsCmd(url string) tea.Cmd {
	return func() tea.Msg {
		if fakeFetchFormats {
			time.Sleep(5 * time.Second)
			return FormatsLoadedMsg{Table: "ID  EXT  RESOLUTION  FPS\n--- --- ----------- ---\n137 mp4  1920x1080    30  (fake)\n140 m4a  audio only   -   (fake)"}
		}

		ctx := context.Background()
		dl := ytdlp.New().Print("formats_table").NoWarnings()
		res, err := dl.Run(ctx, url)
		if err != nil {
			return FormatsLoadedMsg{Err: err}
		}
		return FormatsLoadedMsg{Table: res.Stdout}
	}
}
