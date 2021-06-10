package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"

	"git.bacardi55.io/bacardi55/gtl/core"
)

var defaultConf = []byte(`# Default config file
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
`)

func Init(configArg string, Data *core.TlData) {
	Config := getTlConfig(configArg)
	Data.Config = &Config

	if configureLogs(Config) != nil {
		log.Fatalln("Log file couldn't be created")
	}

	Feeds := getFeeds(Config.Subscribed_data)
	Data.Feeds = Feeds
}

// Get TlConfig.
func getTlConfig(configArg string) core.TlConfig {
	configFile, err := getConfigFilePath(configArg)
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

	toml.Unmarshal(b, Config)
	return nil
}

// Init Config:
// - Try to load the given file if any.
// - Load default configuration file otherwise. Create it if it doesn't exist.
func getConfigFilePath(configArg string) (string, error) {
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

		// TODO: Use filepath.Join() instead:
		configFile = homepath + "/.config/gtl/gtl.toml"

		// Load or create configFile.
		f, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		defer f.Close()
		if err == nil {
			log.Println("Default configuration file does not exist yet, creating it…")
			// Todo: Fix when directory doesn't exist. Doesn't work atm.
			_, err := f.Write(defaultConf)
			if err != nil {
				return "", fmt.Errorf("Default configuration file couldn't be created in default place")
			}
		}
	}

	return configFile, nil
}

func getFeeds(subFile string) map[string]core.TlFeed {
	subFilePath, e := homedir.Expand(subFile)
	if e != nil {
		log.Fatalln(e)
	}

	file, err := os.Open(subFilePath)
	if err != nil {
		os.Create(subFilePath)
	}
	defer func() error {
		if err = file.Close(); err != nil {
			log.Fatalln("Couldn't close subsciption file (%v)\n%v: ", subFile, err)
		}
		return nil
	}()

	F, err := parseSubscriptions(file)
	if err != nil {
		log.Fatalln(err)
	}

	return F
}

// Parse the subscription file.
func parseSubscriptions(content io.Reader) (map[string]core.TlFeed, error) {
	var Feeds map[string]core.TlFeed
	Feeds = make(map[string]core.TlFeed)

	scanner := bufio.NewScanner(content)
  for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		v := strings.Fields(strings.TrimSpace(line))
		if lv := len(v); lv == 2 {
			Feed := core.TlFeed{
				Title: strings.TrimSpace(v[1]),
				Link:  strings.TrimSpace(v[0]),
			}
			Feeds[Feed.Title] = Feed
		} else if lv == 1 {
      tmpTitle := "Anonymous_" + string(i+1)
			Feed := core.TlFeed{
				Title: tmpTitle,
				Link:  v[0],
			}
			Feeds[Feed.Title] = Feed
		} else {
			return Feeds, fmt.Errorf("Ignoring malformated entry: ", line)
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		return Feeds, fmt.Errorf("reading standard input:", err)
	}

	return Feeds, nil
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
