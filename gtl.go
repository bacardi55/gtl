package main

import (
  "fmt"
  "os"
  "log"

  flag "github.com/spf13/pflag"

  "git.bacardi55.io/bacardi55/gtl/core"
  "git.bacardi55.io/bacardi55/gtl/config"
)

var configFile string
var Data = core.TlData{}

func main() {
  var configArg string
  var helpArg, cliArg bool
  var cliLimitArg int
  flag.StringVar(&configArg, "config", "", "The path to gtl.toml config file.")
  flag.BoolVar(&helpArg, "help", false, "Display help.")
  flag.BoolVar(&cliArg, "cli", false, "Display tinylog stream and quit.")
  flag.IntVar(&cliLimitArg, "limit", 0, "Limit number of items in CLI mode")
  flag.Parse()

  if helpArg == true {
    help()
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
  if cliArg == true {
    displayStreamCli(Data.Stream, cliLimitArg)
    os.Exit(0)
  }

  fmt.Println("TUI is coming, only CLI for now, default-ing to CLI display.\n")
  displayStreamCli(Data.Stream, cliLimitArg)
}

// Display help message.
func help() {
  fmt.Println("gtl is a TUI tool for gemini tinylogs")
  fmt.Println("Usage:")
  fmt.Println("\t--config configFile\tIndicate a specific config file.")
  fmt.Println("\t--help\t\t\tDisplay this help message.")
  return
}

func displayStreamCli(stream *core.TlStream, limit int) {
  if limit < 1 {
    limit = len(stream.Items)
  }

  for i := 0; i < limit; i++ {
    fmt.Println(stream.Items[i].Author, "-", stream.Items[i].Published.Format(Data.Config.Date_format))
    fmt.Println(stream.Items[i].Content, "\n")
  }
}
