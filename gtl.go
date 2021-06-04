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
  var configArg, modeArg string
  var helpArg bool
  var cliLimitArg int
  flag.StringVar(&configArg, "config", "", "The path to gtl.toml config file.")
  flag.BoolVar(&helpArg, "help", false, "Display help.")
  flag.StringVar(&modeArg, "mode", "cli", "The mode for gtl, either cli or tui (not ready yet)")
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
  if modeArg == "cli" {
    displayStreamCli(Data.Stream, cliLimitArg)
  } else if modeArg == "tui" {
    fmt.Println("TUI is not available yet, please use the cli mode.")
  }
}

// Display help message.
func help() {
  fmt.Println("gtl is a TUI tool for gemini tinylogs")
  fmt.Println("Usage:")
  fmt.Println("\t--config configFile\tIndicate a specific config file.")
  fmt.Println("\t--mode {cli,tui}\tSelect the cli or tui mode.")
  fmt.Println("\t--limit X\t\tWhen using cli mode, display only X item.")
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
