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
	var configArg, modeArg string
	var helpArg, versionArg bool
	var cliLimitArg int
	flag.StringVar(&configArg, "config", "", "The path to gtl.toml config file.")
	flag.StringVar(&modeArg, "mode", "cli", "The mode for gtl, either cli or tui (not ready yet)")
	flag.IntVar(&cliLimitArg, "limit", 0, "Limit number of items in CLI mode")
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
