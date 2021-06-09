# GTL: Go Tiny Logs

Goal: A TUI for the [tinylogs](https://codeberg.org/bacardi55/gemini-tinylog-rfc/src/branch/main) format on the [gemini](gemini.circumlunar.space/) space.

Status: Work In Progress, only the CLI mode is available at this stage.

# Installation

gtl requires go ≥ 1.5

# From Source
```
git clone http://git.bacardi55.io/bacardi55/gtl.git
cd gtl
make dependencies
make build
```

Put the binary in a directory available in your $PATH if you don't want to type the full path to the program.

# From Binaries

You can download binaries for linux [here](https://github.com/bacardi55/gtl/releases).
Binaries are only available for linux 386, amd64, arm and arm64 for now.

# Usage

## Quick start:

*Assuming you put the binary in ~/bin*:
```
mkdir ~/.config/gtl/ # Not created automatically, known issue.
~/bin/gtl # will create the configuration and subscription files in ~/.config/gtl/
~/bin/gtl add --url gemini://gmi.bacardi55.io/tinylog.gmi # Adding an entry will create the sub file.
# Repeat add command for all the feeds.
~/bin/gtl --mode cli --limit 10
```

## Global commands:
```
gtl --help
gtl --version
```

## Use gtl
```
Usage:
	--config configFile	Indicate a specific config file.
	--mode {cli,tui}	Select the cli or tui mode.
	--limit X		When using cli mode, display only X item.
	--help			Display this help message.
```

Example:
```bash
gtl --mode cli --limit 10
gtl --limit 10 # cli mode is default, so this is the same as above.
gtl --mode cli --limit 10 --config path/to/config/file # with specific path for config file.
```

If you don't provide a config file path, gtl will look for it in `{homepath}/.config/gtl/gtl.toml`

You need a subscription file though with the list of tinylogs to follow. For easier migration, the format is the same as [lace](https://friendo.monster/log/lace.html):
```
<urlOfTinyLog> nameOfTinyLog
<urlOfTinyLog2> nameOfTinyLog2
…
```

**Warning**: The `nameOfTinyLog` is optional. But if you don't indicate one and the tinylog doesn't have an `author: @authorName` metadata, gtl will not no what to display for the author and will indicate "unknown"

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
# Optional: Highlight when text is found in content.
# Separate values by a coma, eg:
# highlights = "@bacardi55, @bacardi, anything"
```

By default, gtl will look for ~/.config/gtl/gtl.toml . It will create it if needed.

The --config option only look for the file, it will not create it if the file given as argument of --config doesn't exist.

**Known Bug**: If you intend to let gtl create the default configuration file, you need to create the ~/.config/gtl directory first as it won't be created automatically.


## Subscription management

You can add and remove tinylog entry either manually from the file directly, or use gtl to do it for you:
```
Subscription management usage:
	add --url url [--title title]	Indicate a new tinylog url (and optional title) to subscribe to.
	rm --url url			Indicate a tinylog url to be removed from the subscription.
```


# Screenshots

![Gtl screenshot](docs/images/gtl_screenshot.png)

# TODOLIST:

MVP:
* ~~Load or create Configuration file.~~
* ~~Load subscription from file set in Configuration~~
* ~~Load and parse all tinylogs feeds:~~
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
* ~~Optional name in subscription file.~~

* ~~Help~~
* ~~Documentation~~
* ~~Format code according to go standard via gofmt.~~

Others:
* ~~Subscription management: Add / Remove tinylogs.~~
* Notification view
* ~~Highlight notifications in CLI~~
* Highlight notifications in TUI
* ~~Create subscription file if doesn't exit.~~
* Add option to limit number of entries per tinylog or a maximum number of days of history for TUI
* Structured logs.

+ All todos in the code…
