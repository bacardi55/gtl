package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"io/ioutil"

	"github.com/mitchellh/go-homedir"

	"git.bacardi55.io/bacardi55/gtl/core"
)

func getFeeds(subFile string) map[string]core.TlFeed {
	subFilePath, e := homedir.Expand(subFile)
	if e != nil {
		log.Fatalln(e)
	}

	var file *os.File
	file, err := os.Open(subFilePath)
	if err != nil {
		log.Println("Creating the subscription file")
		os.Create(subFilePath)
		// Add 1 entry to avoid bugs with empty subscription file.
		if e = AddTlSubscription(subFilePath, "gemini://gmi.bacardi55.io/tinylog.gmi", "bacardi55"); e != nil {
			log.Fatalln(e)
		}
		file, _ = os.Open(subFilePath)
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
			tmpTitle := fmt.Sprintf("Anonymous_%v", (i + 1))
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

func AddTlSubscription(filepath string, sub string, title string) error {
	subFilePath, e := homedir.Expand(filepath)
	if e != nil {
		return e
	}

	f, err := os.OpenFile(subFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Couldn't open subscription file:\n%v", err)
	}
	defer f.Close()

	newEntry := sub
	if strings.TrimSpace(title) != "" {
		newEntry = newEntry + " " + title
	}
	newEntry = newEntry + "\n"
	if _, err := f.WriteString(newEntry); err != nil {
		return fmt.Errorf("Couldn't write new entry in subscription file:\n%v", err)
	}

	return nil
}

func RemoveTlSubscription(filepath string, sub string) error {

	subFilePath, e := homedir.Expand(filepath)
	if e != nil {
		return e
	}

	f, err := os.OpenFile(subFilePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Couldn't open subscription file:\n%v", err)
	}
	defer f.Close()

	fb, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("Couldn't read subscription file; ", filepath)
	}

	re := regexp.MustCompile("(?m)^(" + sub + ".*)\n")
	res := re.ReplaceAllString(string(fb), "")

	err = ioutil.WriteFile(subFilePath, []byte(res), 0644)
	if err != nil {
		return err
	}

	return nil
}
