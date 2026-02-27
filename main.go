package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/kyawphyothu/momo/program"
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
}
