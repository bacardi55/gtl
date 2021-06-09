package ui

import (
	"fmt"
	"strconv"
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
}

func formatElapsedTime(elapsed time.Duration) string {
	ret := elapsed.Round(time.Second).String()

	if d := int(elapsed.Hours() / 24); d > 0 {
		ret = strconv.Itoa(d) + " day"
		if d > 1 {
			ret = ret + "s"
		}
	} else if h := int(elapsed.Hours()); h > 0 {
		ret = strconv.Itoa(h) + " hour"
		if h > 1 {
			ret = ret + "s"
		}
	} else if m := int(elapsed.Minutes()); m > 0 {
		ret = strconv.Itoa(m) + " minute"
		if m > 1 {
			ret = ret + "s"
		}
	} else {
		ret = strconv.Itoa(int(elapsed.Round(time.Second).Seconds())) + " seconds"
	}

	return ret + " ago"
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