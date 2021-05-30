package main

import (
  "fmt"
  "os"
  "log"
  //"sort"

  flag "github.com/spf13/pflag"

  "git.bacardi55.io/bacardi55/gtl/core"
  "git.bacardi55.io/bacardi55/gtl/config"
)

var configFile string
var Data = core.TlData{}

func main() {
  log.Println("Received arguments: ", os.Args)

  var configArg string
  var helpArg, cliArg bool
  flag.StringVar(&configArg, "config", "", "The path to gtl.toml config file.")
  flag.BoolVar(&helpArg, "help", false, "Display help.")
  flag.BoolVar(&cliArg, "cli", false, "Display tinylog stream and quit.")
  flag.Parse()

  if helpArg == true {
    help()
    os.Exit(0)
  }

  // Init configuration and subscriptions.
  config.Init(configArg, &Data)

  //fmt.Println(Data)
  //fmt.Println(Data.Config)
  //log.Println(Data.Feeds)

  // Retrieve feeds and create stream.
  e := Data.RefreshFeeds()
  if e != nil {
    log.Fatalln("Couldn't refresh feeds")
  }

  // Display stream and quit.
  if cliArg == true {
    displayStreamCli(Data.Stream)
  }

  fmt.Println("TUI is coming, only CLI for now, default-ing to CLI display.\n")
  displayStreamCli(Data.Stream)
}

// Display help message.
func help() {
  fmt.Println("gtl is a TUI tool for gemini tinylogs")
  fmt.Println("Usage:")
  fmt.Println("\t--config configFile\tIndicate a specific config file.")
  fmt.Println("\t--help\t\t\tDisplay this help message.")
  return
}

func displayStreamCli(stream *core.TlStream) {
  for _, s := range stream.Items {
    //fmt.Println(i, s)
    fmt.Println(s.Author, "-", s.Published.Format(Data.Config.Date_format))
    fmt.Println(s.Content, "\n")
  }
}
