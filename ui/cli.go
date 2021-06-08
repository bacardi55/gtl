package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"

	"git.bacardi55.io/bacardi55/gtl/core"
)

// Display Stream for CLI output.
func DisplayStreamCli(data *core.TlData, limit int) {
	stream := data.Stream
	if limit < 1 {
		limit = len(stream.Items)
	}

	t := time.Now()

	for i := 0; i < limit; i++ {
		elapsed := t.Sub(stream.Items[i].Published)

		r := color.New(color.FgRed)
		b := color.New(color.FgBlue)

		r.Print(stream.Items[i].Author)
		fmt.Print(" - ")
		b.Print(elapsed.Round(time.Second).String(), " ago")
		fmt.Print(" - ")
		fmt.Println(stream.Items[i].Published.Format(data.Config.Date_format))

		f := false
		if len(data.Config.Highlights) > 0 {
			highlights := strings.Split(data.Config.Highlights, ",")

			if len(highlights) > 0 {
				for _, h := range highlights {
					if strings.Contains(stream.Items[i].Content, strings.TrimSpace(h)) {
						bo := color.New(color.Bold).Add(color.Italic)
						bo.Println(stream.Items[i].Content, "\n")
						f = true
						break
					}
				}
			}
		}
		if f == false {
			fmt.Println(stream.Items[i].Content, "\n")
		}
	}
}

// Display help message.
func Help() {
	fmt.Println("gtl is a TUI tool for gemini tinylogs")
	fmt.Println("Usage:")
	fmt.Println("\t--config configFile\tIndicate a specific config file.")
	fmt.Println("\t--mode {cli,tui}\tSelect the cli or tui mode.")
	fmt.Println("\t--limit X\t\tWhen using cli mode, display only X item.")
	fmt.Println("\t--version\t\tDisplay gtl's current version.")
	fmt.Println("\t--help\t\t\tDisplay this help message.")
	return
}
