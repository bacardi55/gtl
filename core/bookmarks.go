package core

import (
	"encoding/json"
	"log"
	"os"
)

func SaveBookmarksToFile(file string, bookmarks TlBookmarks) error {
	bookmarksJson, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(file, bookmarksJson, 0666)
	if err != nil {
		return err
	}

	return nil
}

func LoadBookmarksFromFile(file string) (TlBookmarks, error) {
	var bookmarks TlBookmarks

	bookmarksFile, err := os.ReadFile(file)
	if err != nil {
		log.Println("Couldn't read bookmark file", file, err)
		return bookmarks, err
	}

	err = json.Unmarshal(bookmarksFile, &bookmarks)

	return bookmarks, err
}
