package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"io/ioutil"

	"github.com/mitchellh/go-homedir"
)

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
