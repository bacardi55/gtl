package main

import (
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"

	"git.bacardi55.io/bacardi55/gtl/config"
	"git.bacardi55.io/bacardi55/gtl/core"
	"git.bacardi55.io/bacardi55/gtl/ui"
)

var configFile string
var Data = core.TlData{}
var Version string

func main() {
	var configArg, modeArg, urlArg, titleArg string
	var helpArg, versionArg bool
	var cliLimitArg int

	flag.StringVar(&configArg, "config", "", "The path to gtl.toml config file.")
	flag.StringVar(&modeArg, "mode", "cli", "The mode for gtl, either cli or tui (not ready yet).")
	flag.StringVar(&urlArg, "url", "", "The tinylog url you want to (un)subscribe to.")
	flag.StringVar(&titleArg, "title", "", "The optional title of the tinylog you want to subscribe to.")
	flag.IntVar(&cliLimitArg, "limit", 0, "Limit number of items in CLI mode.")
	flag.BoolVar(&helpArg, "help", false, "Display help.")
	flag.BoolVar(&versionArg, "version", false, "Display gtl's current version.")
	flag.Parse()

	if helpArg == true {
		ui.Help()
		os.Exit(0)
	}

	if versionArg == true {
		fmt.Printf("gtl version: %s\n", Version)
		os.Exit(0)
	}

	// Init configuration and subscriptions.
	config.Init(configArg, &Data)

	if len(os.Args) > 1 {
		if os.Args[1] == "add" {
			if urlArg == "" {
				ui.SubHelp()
				log.Fatalln("No url to add!")
			}
			e := config.AddTlSubscription(Data.Config.Subscribed_data, urlArg, titleArg)
			if e != nil {
				fmt.Println(e)
				log.Fatalln("Error adding entry.")
			} else {
				fmt.Println(urlArg, "has been added to the subscribed list.")
			}
			os.Exit(0)
		} else if os.Args[1] == "rm" {
			if urlArg == "" {
				ui.SubHelp()
				log.Fatalln("No url to remove")
			}
			e := config.RemoveTlSubscription(Data.Config.Subscribed_data, urlArg)
			if e != nil {
				fmt.Println(e)
				log.Fatalln("Error removing entry:\n")
			} else {
				fmt.Println(urlArg, "has been removed from the subscribed list.")
			}
			os.Exit(0)
		}
	}

	// Retrieve feeds and create stream.
	e := Data.RefreshFeeds()
	if e != nil {
		log.Fatalln("Couldn't refresh feeds")
	}

	// Display stream and quit.
	if modeArg == "cli" {
		ui.DisplayStreamCli(&Data, cliLimitArg)
	} else if modeArg == "tui" {
		ui.DisplayStreamTui(&Data)
	}
}
