package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"

	"git.bacardi55.io/bacardi55/gtl/core"
)

// Display Stream for CLI output.
func DisplayStreamCli(data *core.TlData, limit int) error {
	e := data.RefreshFeeds()
	if e != nil {
		return fmt.Errorf("Couldn't refresh feeds")
	}

	stream := data.Stream

	var max int
	max = len(stream.Items)
	if limit > 0 {
		max = limit
	} else if data.Config.Cli_limit > 0 {
		max = data.Config.Cli_limit
	}

	t := time.Now()

	for i := 0; i < max; i++ {
		elapsed := t.Sub(stream.Items[i].Published)

		r := color.New(color.FgRed)
		b := color.New(color.FgBlue)

		r.Print(stream.Items[i].Author)
		fmt.Print(" - ")
		b.Print(formatElapsedTime(elapsed))
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

	return nil
}

// Display help message.
func Help() {
	fmt.Println("gtl is a CLI and TUI (soon) tool for gemini tinylogs")
	CliHelp()
	fmt.Println()
	SubHelp()
	return
}

func CliHelp() {
	fmt.Println("CLI usage:")
	fmt.Println("\t--config configFile\tIndicate a specific config file.")
	fmt.Println("\t--mode {cli,tui}\tSelect the cli or tui mode.")
	fmt.Println("\t--limit X\t\tWhen using cli mode, display only X item.")
	fmt.Println("\t--version\t\tDisplay gtl's current version.")
	fmt.Println("\t--help\t\t\tDisplay this help message.")
}

func SubHelp() {
	fmt.Println("Subscription management usage:")
	fmt.Println("\tadd --url url [--title title]\tIndicate a new tinylog url (and optional title) to subscribe to.")
	fmt.Println("\trm --url url\t\t\tIndicate a tinylog url to be removed from the subscription.")
}
