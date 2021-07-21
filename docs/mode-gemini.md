# Use gtl gemini mode

![Gtl Gemini output screenshot](docs/images/gtl_gmi_screenshot.png)

You can use gtl to generate a valid text/gemini output that could then be place in a capsule and read via a gemini browser:

```bash
~/bin/gtl --config ~/.config/gtl/test.gtl.toml --mode gemini --limit 55
```

You can see an example used here:
gemini://tinylogs.gmi.bacardi55.io

Or see a screenshot below.

![Gtl Gemini output in a browser screenshot](docs/images/gtl_gmi_screenshot_browser.png)

Funny thing, the format is compatible with the tinylog RFC, so you can subscribe to it via gtl (there is a screenshot of this below too).

![Gtl Gemini output in gtl screenshot](docs/images/gtl_gmi_screenshot_gtl.png)
