# GTL: Go Tiny Logs

Goal: A TUI for the [tinylogs]() format on the [gemini]() space.

Status: Work In Progress, only the CLI mode is available at this stage.

# Installation

gtl requires go ≥ 1.5

# From Source
```
git clone http://git.bacardi55.io/bacardi55/gtl.git
cd gtl
go build -o path/to/binary/folder
```

Put the binary in a directory available in your $PATH if you don't want to type the full path to the program.

# From Binaries

You can download binaries for linux [here](https://github.com/bacardi55/gtl/releases).
Binaries are only available for linux 386, amd64, arm and arm64 for now.

# Usage

```
Usage:
	--config configFile	Indicate a specific config file.
	--mode {cli,tui}	Select the cli or tui mode.
	--limit X		When using cli mode, display only X item.
	--help			Display this help message.
```

If you don't provide a config file path, gtl will look for it in `{homepath}/.config/gtl/gtl.toml`

You need a subscription file though with the list of tinylogs to follow. For easier migration, the format is the same as lace except that the second argument is mandatory (for now):
```
<urlOfTinyLog> nameOfTinyLog
<urlOfTinyLog2> nameOfTinyLog2
…
```

This file should be in your configuration file:

```toml
subscribed_data = "path/to/sub/file"
```


# Default config file

```toml
# Default config file:
# Path to subscribed tinylogs:
subscribed_data = "~/.config/gtl/subs"
# Refresh time:
refresh = 10
# Date display format
date_format = "Mon 02 Jan 2006 15:04 MST"
# Log file:
log_file = "/dev/null"
```

By default, gtl will look for ~/.config/gtl/gtl.toml . It will create it if needed.

The --config option only look for the file, it will not create it if the file given as argument of --config doesn't exist.

# TODOLIST:

MVP:
* ~~Load or create Configuration file.~~
* ~~Load subscription from file set in Configuration~~
* Load and parse all tinylogs feeds:
  * ~~Load .gmi file for all feeds~~
  * ~~Parse header for author and avatar (cf @adele & @szczezuja)~~
  * ~~Parse tinylog entry:~~
    * ~~date~~
    * ~~content~~
* ~~Sort feeds items~~
* ~~Display as CLI output feeds items in order~~
* ~~Limit option for CLI output~~
* Create TUI basic
* Auto Refresh based on configuration.
* Extract links from tinylog entry and display them the gemini way.
* ~~Move logs to logfile instead of stdout~~
* Optional name in subscription file.

* Help / Documentation
* ~~Format code according to go standard via gofmt.~~

Others:
* Subscription management: Add / Remove tinylogs.
* Notification view
* Highlight notifications
* Create subscription file if doesn't exit.
* Add option to limit number of entries per tinylog or a maximum number of days of history for TUI
* Structured logs.

+ All todos in the code…
