package program

import (
	"context"
	"fmt"
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

// fakeFetchFormats enables a 1s fake fetch for testing the loading UI.
const fakeFetchFormats = true

func FetchFormatsCmd(url string) tea.Cmd {
	if url == "" {
		url = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	}
	return func() tea.Msg {
		if fakeFetchFormats {
			time.Sleep(1 * time.Second)
			return FormatsLoadedMsg{Formats: fakeFormats()}
		}

		ctx := context.Background()
		dl := ytdlp.New().DumpJSON().NoWarnings()
		res, err := dl.Run(ctx, url)
		if err != nil {
			return FormatsLoadedMsg{Err: err}
		}
		info, err := res.GetExtractedInfo()
		if err != nil {
			return FormatsLoadedMsg{Err: err}
		}
		if len(info) == 0 {
			return FormatsLoadedMsg{Err: fmt.Errorf("no info returned by yt-dlp")}
		}
		if len(info[0].Formats) == 0 {
			return FormatsLoadedMsg{Err: fmt.Errorf("no formats found")}
		}
		return FormatsLoadedMsg{Formats: info[0].Formats}
	}
}

func fakeFormats() []*ytdlp.ExtractedFormat {
	return []*ytdlp.ExtractedFormat{
		{FormatID: ptr("sb3"), Extension: ptr("mhtml"), Resolution: ptr("48x27"), FPS: ptr(0.0), Protocol: ptr("mhtml"), FormatNote: ptr("storyboard")},
		{FormatID: ptr("sb0"), Extension: ptr("mhtml"), Resolution: ptr("320x180"), FPS: ptr(1.0), Protocol: ptr("mhtml"), FormatNote: ptr("storyboard")},
		{FormatID: ptr("139"), Extension: ptr("m4a"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(1300234), TBR: ptr(49.0), Protocol: ptr("https"), ACodec: ptr("mp4a.40.5"), ABR: ptr(49.0), ASR: ptr(22050.0), FormatNote: ptr("low"), Language: ptr("en")},
		{FormatID: ptr("249"), Extension: ptr("webm"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(1226833), TBR: ptr(46.0), Protocol: ptr("https"), ACodec: ptr("opus"), ABR: ptr(46.0), ASR: ptr(48000.0), FormatNote: ptr("low"), Language: ptr("en")},
		{FormatID: ptr("140"), Extension: ptr("m4a"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(3450060), TBR: ptr(130.0), Protocol: ptr("https"), ACodec: ptr("mp4a.40.2"), ABR: ptr(130.0), ASR: ptr(44100.0), FormatNote: ptr("medium"), Language: ptr("en")},
		{FormatID: ptr("251"), Extension: ptr("webm"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(3429236), TBR: ptr(129.0), Protocol: ptr("https"), ACodec: ptr("opus"), ABR: ptr(129.0), ASR: ptr(48000.0), FormatNote: ptr("medium"), Language: ptr("en")},
		{FormatID: ptr("160"), Extension: ptr("mp4"), Resolution: ptr("256x144"), FPS: ptr(25.0), FileSize: ptr(2055208), TBR: ptr(77.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d400c"), VBR: ptr(77.0), FormatNote: ptr("144p")},
		{FormatID: ptr("278"), Extension: ptr("webm"), Resolution: ptr("256x144"), FPS: ptr(25.0), FileSize: ptr(1541406), TBR: ptr(58.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(58.0), FormatNote: ptr("144p")},
		{FormatID: ptr("134"), Extension: ptr("mp4"), Resolution: ptr("640x360"), FPS: ptr(25.0), FileSize: ptr(8388608), TBR: ptr(315.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401e"), VBR: ptr(315.0), FormatNote: ptr("360p")},
		{FormatID: ptr("18"), Extension: ptr("mp4"), Resolution: ptr("640x360"), FPS: ptr(25.0), AudioChannels: ptr(2.0), FileSizeApprox: ptr(11826381), TBR: ptr(444.0), Protocol: ptr("https"), VCodec: ptr("avc1.42001E"), ACodec: ptr("mp4a.40.2"), ASR: ptr(44100.0), FormatNote: ptr("360p"), Language: ptr("en")},
		{FormatID: ptr("135"), Extension: ptr("mp4"), Resolution: ptr("854x480"), FPS: ptr(25.0), FileSize: ptr(14104780), TBR: ptr(530.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401e"), VBR: ptr(530.0), FormatNote: ptr("480p")},
		{FormatID: ptr("136"), Extension: ptr("mp4"), Resolution: ptr("1280x720"), FPS: ptr(25.0), FileSize: ptr(26455245), TBR: ptr(993.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401f"), VBR: ptr(993.0), FormatNote: ptr("720p")},
		{FormatID: ptr("137"), Extension: ptr("mp4"), Resolution: ptr("1920x1080"), FPS: ptr(25.0), FileSize: ptr(80899727), TBR: ptr(3038.0), Protocol: ptr("https"), VCodec: ptr("avc1.640028"), VBR: ptr(3038.0), FormatNote: ptr("1080p")},
		{FormatID: ptr("248"), Extension: ptr("webm"), Resolution: ptr("1920x1080"), FPS: ptr(25.0), FileSize: ptr(30849163), TBR: ptr(1158.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(1158.0), FormatNote: ptr("1080p")},
		{FormatID: ptr("271"), Extension: ptr("webm"), Resolution: ptr("2560x1440"), FPS: ptr(25.0), FileSize: ptr(151080550), TBR: ptr(5674.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(5674.0), FormatNote: ptr("1440p")},
		{FormatID: ptr("401"), Extension: ptr("mp4"), Resolution: ptr("3840x2160"), FPS: ptr(25.0), FileSize: ptr(240338534), TBR: ptr(9025.0), Protocol: ptr("https"), VCodec: ptr("av01.0.12M.08"), VBR: ptr(9025.0), FormatNote: ptr("2160p")},
		{FormatID: ptr("sb3"), Extension: ptr("mhtml"), Resolution: ptr("48x27"), FPS: ptr(0.0), Protocol: ptr("mhtml"), FormatNote: ptr("storyboard")},
		{FormatID: ptr("sb0"), Extension: ptr("mhtml"), Resolution: ptr("320x180"), FPS: ptr(1.0), Protocol: ptr("mhtml"), FormatNote: ptr("storyboard")},
		{FormatID: ptr("139"), Extension: ptr("m4a"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(1300234), TBR: ptr(49.0), Protocol: ptr("https"), ACodec: ptr("mp4a.40.5"), ABR: ptr(49.0), ASR: ptr(22050.0), FormatNote: ptr("low"), Language: ptr("en")},
		{FormatID: ptr("249"), Extension: ptr("webm"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(1226833), TBR: ptr(46.0), Protocol: ptr("https"), ACodec: ptr("opus"), ABR: ptr(46.0), ASR: ptr(48000.0), FormatNote: ptr("low"), Language: ptr("en")},
		{FormatID: ptr("140"), Extension: ptr("m4a"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(3450060), TBR: ptr(130.0), Protocol: ptr("https"), ACodec: ptr("mp4a.40.2"), ABR: ptr(130.0), ASR: ptr(44100.0), FormatNote: ptr("medium"), Language: ptr("en")},
		{FormatID: ptr("251"), Extension: ptr("webm"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(3429236), TBR: ptr(129.0), Protocol: ptr("https"), ACodec: ptr("opus"), ABR: ptr(129.0), ASR: ptr(48000.0), FormatNote: ptr("medium"), Language: ptr("en")},
		{FormatID: ptr("160"), Extension: ptr("mp4"), Resolution: ptr("256x144"), FPS: ptr(25.0), FileSize: ptr(2055208), TBR: ptr(77.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d400c"), VBR: ptr(77.0), FormatNote: ptr("144p")},
		{FormatID: ptr("278"), Extension: ptr("webm"), Resolution: ptr("256x144"), FPS: ptr(25.0), FileSize: ptr(1541406), TBR: ptr(58.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(58.0), FormatNote: ptr("144p")},
		{FormatID: ptr("134"), Extension: ptr("mp4"), Resolution: ptr("640x360"), FPS: ptr(25.0), FileSize: ptr(8388608), TBR: ptr(315.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401e"), VBR: ptr(315.0), FormatNote: ptr("360p")},
		{FormatID: ptr("18"), Extension: ptr("mp4"), Resolution: ptr("640x360"), FPS: ptr(25.0), AudioChannels: ptr(2.0), FileSizeApprox: ptr(11826381), TBR: ptr(444.0), Protocol: ptr("https"), VCodec: ptr("avc1.42001E"), ACodec: ptr("mp4a.40.2"), ASR: ptr(44100.0), FormatNote: ptr("360p"), Language: ptr("en")},
		{FormatID: ptr("135"), Extension: ptr("mp4"), Resolution: ptr("854x480"), FPS: ptr(25.0), FileSize: ptr(14104780), TBR: ptr(530.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401e"), VBR: ptr(530.0), FormatNote: ptr("480p")},
		{FormatID: ptr("136"), Extension: ptr("mp4"), Resolution: ptr("1280x720"), FPS: ptr(25.0), FileSize: ptr(26455245), TBR: ptr(993.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401f"), VBR: ptr(993.0), FormatNote: ptr("720p")},
		{FormatID: ptr("137"), Extension: ptr("mp4"), Resolution: ptr("1920x1080"), FPS: ptr(25.0), FileSize: ptr(80899727), TBR: ptr(3038.0), Protocol: ptr("https"), VCodec: ptr("avc1.640028"), VBR: ptr(3038.0), FormatNote: ptr("1080p")},
		{FormatID: ptr("248"), Extension: ptr("webm"), Resolution: ptr("1920x1080"), FPS: ptr(25.0), FileSize: ptr(30849163), TBR: ptr(1158.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(1158.0), FormatNote: ptr("1080p")},
		{FormatID: ptr("271"), Extension: ptr("webm"), Resolution: ptr("2560x1440"), FPS: ptr(25.0), FileSize: ptr(151080550), TBR: ptr(5674.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(5674.0), FormatNote: ptr("1440p")},
		{FormatID: ptr("401"), Extension: ptr("mp4"), Resolution: ptr("3840x2160"), FPS: ptr(25.0), FileSize: ptr(240338534), TBR: ptr(9025.0), Protocol: ptr("https"), VCodec: ptr("av01.0.12M.08"), VBR: ptr(9025.0), FormatNote: ptr("2160p")},
		{FormatID: ptr("sb3"), Extension: ptr("mhtml"), Resolution: ptr("48x27"), FPS: ptr(0.0), Protocol: ptr("mhtml"), FormatNote: ptr("storyboard")},
		{FormatID: ptr("sb0"), Extension: ptr("mhtml"), Resolution: ptr("320x180"), FPS: ptr(1.0), Protocol: ptr("mhtml"), FormatNote: ptr("storyboard")},
		{FormatID: ptr("139"), Extension: ptr("m4a"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(1300234), TBR: ptr(49.0), Protocol: ptr("https"), ACodec: ptr("mp4a.40.5"), ABR: ptr(49.0), ASR: ptr(22050.0), FormatNote: ptr("low"), Language: ptr("en")},
		{FormatID: ptr("249"), Extension: ptr("webm"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(1226833), TBR: ptr(46.0), Protocol: ptr("https"), ACodec: ptr("opus"), ABR: ptr(46.0), ASR: ptr(48000.0), FormatNote: ptr("low"), Language: ptr("en")},
		{FormatID: ptr("140"), Extension: ptr("m4a"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(3450060), TBR: ptr(130.0), Protocol: ptr("https"), ACodec: ptr("mp4a.40.2"), ABR: ptr(130.0), ASR: ptr(44100.0), FormatNote: ptr("medium"), Language: ptr("en")},
		{FormatID: ptr("251"), Extension: ptr("webm"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(3429236), TBR: ptr(129.0), Protocol: ptr("https"), ACodec: ptr("opus"), ABR: ptr(129.0), ASR: ptr(48000.0), FormatNote: ptr("medium"), Language: ptr("en")},
		{FormatID: ptr("160"), Extension: ptr("mp4"), Resolution: ptr("256x144"), FPS: ptr(25.0), FileSize: ptr(2055208), TBR: ptr(77.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d400c"), VBR: ptr(77.0), FormatNote: ptr("144p")},
		{FormatID: ptr("278"), Extension: ptr("webm"), Resolution: ptr("256x144"), FPS: ptr(25.0), FileSize: ptr(1541406), TBR: ptr(58.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(58.0), FormatNote: ptr("144p")},
		{FormatID: ptr("134"), Extension: ptr("mp4"), Resolution: ptr("640x360"), FPS: ptr(25.0), FileSize: ptr(8388608), TBR: ptr(315.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401e"), VBR: ptr(315.0), FormatNote: ptr("360p")},
		{FormatID: ptr("18"), Extension: ptr("mp4"), Resolution: ptr("640x360"), FPS: ptr(25.0), AudioChannels: ptr(2.0), FileSizeApprox: ptr(11826381), TBR: ptr(444.0), Protocol: ptr("https"), VCodec: ptr("avc1.42001E"), ACodec: ptr("mp4a.40.2"), ASR: ptr(44100.0), FormatNote: ptr("360p"), Language: ptr("en")},
		{FormatID: ptr("135"), Extension: ptr("mp4"), Resolution: ptr("854x480"), FPS: ptr(25.0), FileSize: ptr(14104780), TBR: ptr(530.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401e"), VBR: ptr(530.0), FormatNote: ptr("480p")},
		{FormatID: ptr("136"), Extension: ptr("mp4"), Resolution: ptr("1280x720"), FPS: ptr(25.0), FileSize: ptr(26455245), TBR: ptr(993.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401f"), VBR: ptr(993.0), FormatNote: ptr("720p")},
		{FormatID: ptr("137"), Extension: ptr("mp4"), Resolution: ptr("1920x1080"), FPS: ptr(25.0), FileSize: ptr(80899727), TBR: ptr(3038.0), Protocol: ptr("https"), VCodec: ptr("avc1.640028"), VBR: ptr(3038.0), FormatNote: ptr("1080p")},
		{FormatID: ptr("248"), Extension: ptr("webm"), Resolution: ptr("1920x1080"), FPS: ptr(25.0), FileSize: ptr(30849163), TBR: ptr(1158.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(1158.0), FormatNote: ptr("1080p")},
		{FormatID: ptr("271"), Extension: ptr("webm"), Resolution: ptr("2560x1440"), FPS: ptr(25.0), FileSize: ptr(151080550), TBR: ptr(5674.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(5674.0), FormatNote: ptr("1440p")},
		{FormatID: ptr("401"), Extension: ptr("mp4"), Resolution: ptr("3840x2160"), FPS: ptr(25.0), FileSize: ptr(240338534), TBR: ptr(9025.0), Protocol: ptr("https"), VCodec: ptr("av01.0.12M.08"), VBR: ptr(9025.0), FormatNote: ptr("2160p")},
		{FormatID: ptr("sb3"), Extension: ptr("mhtml"), Resolution: ptr("48x27"), FPS: ptr(0.0), Protocol: ptr("mhtml"), FormatNote: ptr("storyboard")},
		{FormatID: ptr("sb0"), Extension: ptr("mhtml"), Resolution: ptr("320x180"), FPS: ptr(1.0), Protocol: ptr("mhtml"), FormatNote: ptr("storyboard")},
		{FormatID: ptr("139"), Extension: ptr("m4a"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(1300234), TBR: ptr(49.0), Protocol: ptr("https"), ACodec: ptr("mp4a.40.5"), ABR: ptr(49.0), ASR: ptr(22050.0), FormatNote: ptr("low"), Language: ptr("en")},
		{FormatID: ptr("249"), Extension: ptr("webm"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(1226833), TBR: ptr(46.0), Protocol: ptr("https"), ACodec: ptr("opus"), ABR: ptr(46.0), ASR: ptr(48000.0), FormatNote: ptr("low"), Language: ptr("en")},
		{FormatID: ptr("140"), Extension: ptr("m4a"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(3450060), TBR: ptr(130.0), Protocol: ptr("https"), ACodec: ptr("mp4a.40.2"), ABR: ptr(130.0), ASR: ptr(44100.0), FormatNote: ptr("medium"), Language: ptr("en")},
		{FormatID: ptr("251"), Extension: ptr("webm"), Resolution: ptr("audio only"), AudioChannels: ptr(2.0), FileSize: ptr(3429236), TBR: ptr(129.0), Protocol: ptr("https"), ACodec: ptr("opus"), ABR: ptr(129.0), ASR: ptr(48000.0), FormatNote: ptr("medium"), Language: ptr("en")},
		{FormatID: ptr("160"), Extension: ptr("mp4"), Resolution: ptr("256x144"), FPS: ptr(25.0), FileSize: ptr(2055208), TBR: ptr(77.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d400c"), VBR: ptr(77.0), FormatNote: ptr("144p")},
		{FormatID: ptr("278"), Extension: ptr("webm"), Resolution: ptr("256x144"), FPS: ptr(25.0), FileSize: ptr(1541406), TBR: ptr(58.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(58.0), FormatNote: ptr("144p")},
		{FormatID: ptr("134"), Extension: ptr("mp4"), Resolution: ptr("640x360"), FPS: ptr(25.0), FileSize: ptr(8388608), TBR: ptr(315.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401e"), VBR: ptr(315.0), FormatNote: ptr("360p")},
		{FormatID: ptr("18"), Extension: ptr("mp4"), Resolution: ptr("640x360"), FPS: ptr(25.0), AudioChannels: ptr(2.0), FileSizeApprox: ptr(11826381), TBR: ptr(444.0), Protocol: ptr("https"), VCodec: ptr("avc1.42001E"), ACodec: ptr("mp4a.40.2"), ASR: ptr(44100.0), FormatNote: ptr("360p"), Language: ptr("en")},
		{FormatID: ptr("135"), Extension: ptr("mp4"), Resolution: ptr("854x480"), FPS: ptr(25.0), FileSize: ptr(14104780), TBR: ptr(530.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401e"), VBR: ptr(530.0), FormatNote: ptr("480p")},
		{FormatID: ptr("136"), Extension: ptr("mp4"), Resolution: ptr("1280x720"), FPS: ptr(25.0), FileSize: ptr(26455245), TBR: ptr(993.0), Protocol: ptr("https"), VCodec: ptr("avc1.4d401f"), VBR: ptr(993.0), FormatNote: ptr("720p")},
		{FormatID: ptr("137"), Extension: ptr("mp4"), Resolution: ptr("1920x1080"), FPS: ptr(25.0), FileSize: ptr(80899727), TBR: ptr(3038.0), Protocol: ptr("https"), VCodec: ptr("avc1.640028"), VBR: ptr(3038.0), FormatNote: ptr("1080p")},
		{FormatID: ptr("248"), Extension: ptr("webm"), Resolution: ptr("1920x1080"), FPS: ptr(25.0), FileSize: ptr(30849163), TBR: ptr(1158.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(1158.0), FormatNote: ptr("1080p")},
		{FormatID: ptr("271"), Extension: ptr("webm"), Resolution: ptr("2560x1440"), FPS: ptr(25.0), FileSize: ptr(151080550), TBR: ptr(5674.0), Protocol: ptr("https"), VCodec: ptr("vp9"), VBR: ptr(5674.0), FormatNote: ptr("1440p")},
		{FormatID: ptr("401"), Extension: ptr("mp4"), Resolution: ptr("3840x2160"), FPS: ptr(25.0), FileSize: ptr(240338534), TBR: ptr(9025.0), Protocol: ptr("https"), VCodec: ptr("av01.0.12M.08"), VBR: ptr(9025.0), FormatNote: ptr("2160p")},
	}
}
