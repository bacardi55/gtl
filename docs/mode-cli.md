# Use gtl CLI

![Gtl CLI screenshot](docs/images/gtl_screenshot.png)

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
