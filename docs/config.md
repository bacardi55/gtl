# Configuration

## Default configuration file

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
# Ignore Highlights from author:
# Must be the full author displayed name!
# ignore_author_highlights = "🤔 @bacardi55"
# Maximum number of entries showed in cli mode. If --limit is used, it will overide this setting.
# Will be ignored in tui mode.
cli_limit = 10
# Mode: either cli, tui or gemini
mode = "tui"
# If false, standard ascii characters will be used.
tui_status_emoji = false

# Enable tinylog edition from gtl:
# This will use an external editor,
# configured in your EDITOR environment variable.
# You can check with 'echo $EDITOR' to see if it
# is configured correctly.
# ctrl+n is disabled when set to false.
# Settings available since v0.5.0
allow_edit = false
# Path to tinylog file. This option is ignored if
# allow_edit = false.
# If not a valid file, editing will not be possible
# and ctrl+n will be disabled.
tinylog_path = "path/to/tinylog/file.gmi"
# Path to script to be executed after the edition is done.
# This script needs to be executable.
# If not a valid executable script, it will be ignored.
post_edit_script = "path/to/script"
# Auto refresh feeds after editing the tinylog file.
# Only used when allow_edit = true
# Settings available since v0.5.1
post_edit_refresh = false
# Limit the number of entries displayed in TUI.
# Indicate 0 for all entries.
# Settings available since v0.6.0
tui_max_entries = 0
# Copy a pre formatted text to clipboard when creating a new entry
# On linux, requires 'xclip' or 'xsel'
# Settings available since v0.6.0
tui_copy_stub_clipboard = false
# If you are running gtl without X, the copy to clipboard feature
# will not work (Or if you don't have `xclip` or `xsel`).
# In this case, enabling this option will allow gtl to display
# the sub text in a modal for easy copy in tools like tmux
# At this stage, the rendering is ugly because of an issue in cview:
# https://code.rocketnine.space/tslocum/cview/issues/72#issuecomment-3968
tui_show_stub = false
# TUI Theming:
# Please see the TUI and theming documentation page in docs/mode-tui.md and docs/theming.md
# Dracula theme example:
#tui_color_background = "282a36"
#tui_color_text = "f8f8f2"
#tui_color_author_name = "ffb86c"
#tui_color_links = "8be9fd"
#tui_color_box = "44475a"
#tui_color_focus_box = "f8f8f2"
#tui_color_elapsed_time = "f1fa8c"
#tui_color_selected_background = "f8f8f2"
#tui_color_selected_foreground = "282a36"
#tui_color_highlight = "ff79c6"
#tui_color_quote = "f1fa8c"
#tui_color_button_color = "44475a"
#tui_color_button_text = "f8f8f2"
#tui_color_button_focus = "ffb86c"
#tui_color_button_focus_text = "282a36"

# Bookmarks (only available via TUI for now):
# Enable/disable bookmarks feature:
bookmarks_enabled = false
# Path to bookmarks file:
bookmarks_file = "~/.config/gtl/bookmarks.json"
# Reverse order (from newer to older) instead of default (older to newer):
bookmarks_reverse_order = true
```

By default, gtl will look for ~/.config/gtl/gtl.toml . It will create it if needed.

The --config option only look for the file, it will not create it if the file given as argument of --config doesn't exist.

## Subscription management

You can add and remove tinylog entry either manually from the file directly, or use gtl to do it for you:
```
Subscription management usage:
	add --url url [--title title]	Indicate a new tinylog url (and optional title) to subscribe to.
	rm --url url			Indicate a tinylog url to be removed from the subscription.
```

