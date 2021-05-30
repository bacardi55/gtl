# GTL: Go Tiny Logs

Goal: A TUI for the [tinylogs]() format on the [gemini]() space.

Status: Work In Progress, unusable at this stage.

# TODOLIST:

MVP:
* ~~Load or create Configuration file.~~
* ~~Load subscription from file set in Configuration~~
* Load and parse all tinylogs feeds:
  * ~~Load .gmi file for all feeds~~
  * ~~Parse header for author and avatar (cf @adele & @szczezuja)~~
  * Parse tinylog entry:
    * ~~date~~
    * ~~content~~
* ~~Sort feeds items~~
* ~~Display as CLI output feeds items in order~~
* Create TUI basic
* Auto Refresh based on configuration.
* Extract links from tinylog entry and display them the gemini way.
* Non verbose mode for logs.
* Move logs to logfile instead of stdout

* Help / Documentation
* Format code according to go standard via gofmt.

Others:
* Subscription management: Add / Remove tinylogs.
* Notification view
* Highlight notifications
* Create subscription file if doesn't exit.
* Add option to limit number of entries per tinylog or a maxium number of days of history

+ All todos in the code…
