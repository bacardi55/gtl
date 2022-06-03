# TUI Theming

Since version 0.7.0, the TUI theme can be customized in `gtl.toml` file. If no custom colors are specified, the default dark theme will apply.

## Theme configuration

This is the configuration for the default dark theme:

**Do not put a starting `#` for the color code**

```toml
# Default (dark) theme:
tui_color_background = "000000"
tui_color_box = "FFFFFF"
tui_color_focus_box = "008000"
tui_color_author_name = "FF0000"
tui_color_links = "87CEEB"
tui_color_elapsed_time = "87CEEB"
tui_color_text = "FFFFFF"
tui_color_selected_background = "FFFFFF"
tui_color_selected_foreground = "000000"
tui_color_highlight = "FF0000"
tui_color_quote = "808080"
tui_color_button_color = "ffffff"
tui_color_button_text = "000000"
tui_color_button_focus = "008000"
tui_color_button_focus_text = "ffffff"
```

![Gtl TUI screenshot](/docs/images/gtl_tui_screenshot.png)

## Theme example

If you create your own theme, feel free to share it with other here or by contacting me!

### Dracula inspired theme

```toml
# Dracula theme example:
tui_color_background = "282a36"
tui_color_text = "f8f8f2"
tui_color_author_name = "ffb86c"
tui_color_links = "8be9fd"
tui_color_box = "44475a"
tui_color_focus_box = "f8f8f2"
tui_color_elapsed_time = "f1fa8c"
tui_color_selected_background = "f8f8f2"
tui_color_selected_foreground = "44475a"
tui_color_highlight = "ff79c6"
tui_color_quote = "f1fa8c"
tui_color_button_color = "44475a"
tui_color_button_text = "f8f8f2"
tui_color_button_focus = "ffb86c"
tui_color_button_focus_text = "282a36"
```

![Gtl TUI screenshot (dracula based theme example)](/docs/images/gtl_tui_screenshot_dracula.png)


### Light theme example

```toml
# Light theme example:
tui_color_background = "ffffff"
tui_color_box = "000000"
tui_color_focus_box = "ff0000"
tui_color_author_name = "ff0000"
tui_color_links = "0000ff"
tui_color_elapsed_time = "0000ff"
tui_color_text = "111111"
tui_color_selected_background = "000000"
tui_color_selected_foreground = "ffffff"
tui_color_highlight = "00f0f0"
tui_color_quote = "ff5555"
tui_color_button_color = "000000"
tui_color_button_text = "ffffff"
tui_color_button_focus = "ff0000"
tui_color_button_focus_text = "111111"
```

![Gtl TUI screenshot (light theme example)](/docs/images/gtl_tui_screenshot_light.png)

