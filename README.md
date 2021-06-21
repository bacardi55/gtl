# GTL: Go Tiny Logs

Goal: A TUI for the [tinylogs](https://codeberg.org/bacardi55/gemini-tinylog-rfc/src/branch/main) format on the [gemini](gemini.circumlunar.space/) space.

This is a early version that still miss a lot of things, from features to tests, so use at your own risk :)

# Installation

gtl requires go ≥ 1.16

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

*PS*: I don't have a Mac or Windows box easily accessible, so any help/patch on this is appreciated :).

# Usage

## Quick start:

*Assuming you put the binary in ~/bin*:
```
~/bin/gtl # will create the configuration and subscription files in ~/.config/gtl/ and a subscription path.
[~/bin/gtl add --url gemini://capsule.tld/tinylog.gmi # Adding other entries than just my default one.
[# Repeat add command for all the feeds.]
~/bin/gtl --mode cli --limit 10
# Or use the TUI mode:
~/bin/gtl --mode tui
```

You can still use gtl "continuously":
```
while true; do clear && ~/bin/gtl --mode cli --limit 10 && sleep 1800; done;
```

Or use the TUI and refresh the timeline when you want.

## Global commands:
```
gtl --help
gtl --version
```

## Use gtl TUI

[Screenshot of the TUI below](#tui-mode).

```
gtl --mode tui
```
or configure the `gtl.toml` to set `mode = tui` (see config below).

PS: TUI is the default mode since v0.4.8


**TUI Shortcuts:**
```
?: Display help
r: Refresh timeline (refresh all tinylogs)
t: Display timeline (remove all filters like highlights of specific tinylog)
h: Toggle Highligts only / all entries (keep tinylog filter).
s: Toggle hide/show subscription sidebar (left).
Tab: Switch between timeline and subscription list.
Arrow keys / hjkl: navigate
q or Ctrl-C: Quit
```
You can navigate on the subscription list and:
* left click or press enter: Will filter only entries from this tinylog and hide all entries from other tinylogs. A Status `F` or `🔎` is indicated.
* right click or press alt+enter: Will open a menu to mute / unmute a tinylog. A tinylog muted means no entry from this tinylog are displayed. A Status `M` or `🔕` is displayed.

**TUI Emoji Status:**

If `tui_status_emoji` is set to true in the configuration file (see below), emoji will be used for the status. Otherwise, simple ASCII characters will be used.

* `V` or `✔`: All good :)
* `X` or `❌`: Indicates that the feed format is wrong or that no entries has been found.
* `D` or `☠️ `: Indicates that the capsule/page is unreachable.
* `S` or `🔓`: Indicates an error with the SSL certificate.
* `F` or `🔎`: Indicates that the feed is selected. It means only entries from this tinylog are displayed.
* `M` or `🔕`: Indicates that the feed is muted. It means no entry of this tinylog will be displayed.

## Use gtl CLI

[Screenshot of the CLI below](#cli-mode).

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

Screenshot of the CLI below.

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
# Maximum number of entries showed in cli mode. If --limit is used, it will overide this setting.
# Will be ignored in tui mode.
cli_limit = 10
# Mode: either cli or tui
mode = "tui"
# If false, standard ascii characters will be used.
tui_status_emoji = false
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

## TUI mode

![Gtl TUI screenshot](docs/images/gtl_tui_screenshot.png)

## CLI mode

![Gtl CLI screenshot](docs/images/gtl_screenshot.png)

