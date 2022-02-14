# Use gtl TUI

![Gtl TUI screenshot](/docs/images/gtl_tui_screenshot.png)

![Gtl TUI screenshot (light theme example)](/docs/images/gtl_tui_screenshot_light.png)

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
Arrow keys up/down or j/k: scroll timeline or feeds list
J/K: Navigate tinylog entries (to select entries) # Available only since v0.6.0
N: Open tinylog in $EDITOR and optionaly copy a new entry stub to clipboard. (See configuration below.) # since v0.6.0
R: Open tinylog in $EDITOR and optionaly copy a response stub to the specific entry. (See configuration below.) # since v0.6.0
O: Open link(s) in selected entry. # since v0.6.0
T: If the selected entry is a response to another tinylog entry, will open the original entry in a popup.
q or Ctrl-C: Quit
```
You can navigate on the subscription list and:
* left click or press enter: Will filter only entries from this tinylog and hide all entries from other tinylogs. A Status `F` or `🔎` is indicated.
* right click or press alt+enter: Will open a menu to mute / unmute a tinylog. A tinylog muted means no entry from this tinylog are displayed. A Status `M` or `🔕` is displayed.

**TUI Theming**

Please read the [theming documentation](/docs/mode-tui-theming.md) to customize colors in TUI mode.

**TUI Emoji Status:**

If `tui_status_emoji` is set to true in the configuration file (see below), emoji will be used for the status. Otherwise, simple ASCII characters will be used.

* `V` or `✔`: All good :)
* `X` or `❌`: Indicates that the feed format is wrong or that no entries has been found.
* `D` or `☠️ `: Indicates that the capsule/page is unreachable.
* `S` or `🔓`: Indicates an error with the SSL certificate.
* `F` or `🔎`: Indicates that the feed is selected. It means only entries from this tinylog are displayed.
* `M` or `🔕`: Indicates that the feed is muted. It means no entry of this tinylog will be displayed.
