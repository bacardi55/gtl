package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	// TODO: Remove ioutil dependencies since go ≥ 1.16.
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"

	"git.bacardi55.io/bacardi55/gtl/core"
)

var defaultConf = []byte(`# Default config file
## General configuration:
# Path to subscribed tinylogs:
subscribed_data = "~/.config/gtl/subs"
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
# Mode: either cli or tui
mode = "tui"
# If false, standard ascii characters will be used.
tui_status_emoji = false

## TUI options:
# Enable tinylog edition from gtl:
# This will use an external editor,
# configured in your EDITOR environment variable.
# You can check with 'echo $EDITOR' to see if it
# is configured correctly.
# ctrl+n is disabled when set to false.
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
post_edit_refresh = false
# Limit the number of entries displayed in TUI.
# Indicate 0 for all entries.
tui_max_entries = 0
# Copy a pre formatted text to clipboard when creating a new entry
# On linux, requires 'xclip' or 'xsel'
tui_copy_stub_clipboard = false
# If you are running gtl without X, the copy to clipboard feature
# will not work (Or if you don't have xclip or xsel).
# In this case, enabling this option will allow gtl to display
# the sub text in a modal for easy copy in tools like tmux
# At this stage, the rendering is ugly because of an issue in cview:
# https://code.rocketnine.space/tslocum/cview/issues/72#issuecomment-3968
tui_show_stub = false

## TUI Theming:
# Define colors for GTL TUI mode.
# Default (dark) theme:
#tui_color_background = "000000"
#tui_color_box = "FFFFFF"
#tui_color_focus_box = "008000"
#tui_color_author_name = "FF0000"
#tui_color_links = "87CEEB"
#tui_color_elapsed_time = "87CEEB"
#tui_color_text = "FFFFFF"
#tui_color_selected_background = "FFFFFF"
#tui_color_selected_foreground = "000000"
#tui_color_highlight = "FF0000"
#tui_color_quote = "808080"

# Bookmarks (only available via TUI for now):
bookmarks_enabled = false
bookmarks_file_path = "~/.config/gtl/bookmarks.json"
bookmarks_reverse_order = true
`)

func Init(configArg string, Data *core.TlData) {
	Config := getTlConfig(configArg)
	Data.Config = &Config

	if configureLogs(Config) != nil {
		log.Fatalln("Log file couldn't be created")
	}

	Feeds := getFeeds(Config.Subscribed_data)
	Data.Feeds = Feeds

	if Data.Config.Bookmarks_enabled {
		log.Println("Trying to enable bookmarks")
		filepath, err := loadBookmarkFile(Data.Config.Bookmarks_file_path)
		if err != nil {
			log.Println("Couldn't load bookmarks", err)
			log.Println("Disabling bookmarks")
			Data.Config.Bookmarks_enabled = false
		}

		bookmarks, e := loadBookmarks(filepath)
		if e != nil {
			log.Println("Couldn't init bookmarks", e)
			log.Println("Disabling bookmarks")
			Data.Config.Bookmarks_enabled = false
		}

		// Update real path to file in config object for reuse:
		Data.Config.Bookmarks_file_path = filepath
		// Save bookmarks in Data.
		Data.Bookmarks = &bookmarks
	}
}

// Get TlConfig.
func getTlConfig(configArg string) core.TlConfig {
	configFile, err := getConfigFilePath(configArg, "gtl.toml")
	if err != nil {
		log.Fatalln(err)
	}

	Config := core.TlConfig{}
	err = loadConfig(configFile, &Config)
	if err != nil {
		log.Fatalln(err)
	}

	return Config
}

// Load configuration file.
func loadConfig(configFile string, Config *core.TlConfig) error {
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("Couldn't open configuration file (%v)\n%v: ", configFile, err)
	}
	defer func() error {
		if err = file.Close(); err != nil {
			return fmt.Errorf("Couldn't close configuration file (%v)\n%v: ", configFile, err)
		}
		return nil
	}()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Couldn't read configuration file; ", configFile)
	}

	return toml.Unmarshal(b, Config)
}

// Init Config:
// - Try to load the given file if any.
// - Load default configuration file otherwise. Create it if it doesn't exist.
func getConfigFilePath(configArg string, defaultFilename string) (string, error) {
	var configFile string
	if configArg != "" {
		log.Println("Trying to load provided configuration file (%v)", configArg)
		if e := fileExist(configArg); e != nil {
			return "", e
		}

		configFile = configArg
		log.Println("Configuration file found:", configFile, "\n")

	} else {
		homepath, err := homedir.Dir()
		if err != nil {
			return "", fmt.Errorf("Error finding home directory")
		}

		configDir := filepath.Join(homepath, ".config", "gtl")
		if err = fileExist(configDir); err != nil {
			log.Println("GTL directory doesn't exist\n", err)
			os.Mkdir(configDir, 0744)
		}

		configFile = filepath.Join(configDir, defaultFilename)

		// Load or create configFile.
		log.Println("Trying to load default configuration file (%v)", configFile)
		f, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		defer f.Close()
		if err == nil {
			log.Println("Default configuration (" + defaultFilename + ") file does not exist yet, creating it…")
			_, err := f.Write(defaultConf)
			if err != nil {
				return "", fmt.Errorf("Default configuration file (" + defaultFilename + ") couldn't be created in default place")
			}
		}
	}

	return configFile, nil
}

// Check if a file exist and if not return a custom error.
func fileExist(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("File (%v) not found", filename)
	}
	return nil
}

// Configure Log file.
func configureLogs(config core.TlConfig) error {
	var logFile string
	if config.Log_file != "" {
		logFile = config.Log_file
	} else {
		logFile = "gtl.log"
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)

	return nil
}

// Configure Bookmarks.
func loadBookmarks(filepath string) (core.TlBookmarks, error) {
	bookmarks, err := core.LoadBookmarksFromFile(filepath)
	return bookmarks, err
}

func loadBookmarkFile(bookmarkFile string) (string, error) {
	bookmarks_path, err := homedir.Expand(bookmarkFile)
	if err != nil {
		return "", err
	}

	filepath, err := getConfigFilePath(bookmarks_path, "bookmarks.json")

	// Given bookmark file that doesn't exist…
	// Let's create it!
	if len(bookmarkFile) > 0 && err != nil {
		bfp, e := homedir.Expand(bookmarkFile)
		if e != nil {
			return "", e
		}

		f, err := os.Create(bfp)
		defer f.Close()
		if err == nil {
			log.Println("Creating custom bookmarks file: " + bookmarkFile)
			filepath = bfp
		}
	}

	return filepath, nil
}
