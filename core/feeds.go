package core

import (
	"context"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"git.sr.ht/~adnano/go-gemini"
)

var supportedTimeFormat = []string{
	"ANSIC",
	"UnixDate",
	"RubyDate",
	"RFC822",
	"RFC822Z",
	"RFC850",
	"RFC1123",
	"RFC1123Z",
	"RFC3339",
	"2006-01-02 15:04:05 MST",
	"2006-01-02 15:04 MST",
	"Mon 02 Jan 2006 03:04:05 PM MST",
	"Mon 02 Jan 2006 03:04 PM MST",
	"Mon 02 Jan 2006 15:04 MST",
	"Mon Jan  2 15:04:05 MST 2006",
	"Mon Jan 02 15:04:05 MST 2006",
	"Mon Jan 02 03:04:05 PM MST 2006",
}

type TlRawFeed struct {
	Name    string
	Content string
}

// Refresh Stream.
func refreshStream(data TlData) (*TlStream, error) {
	var wg sync.WaitGroup

	numFeed := len(data.Feeds)
	chFeedContent := make(chan TlRawFeed, numFeed)
	chFeedError := make(chan error, numFeed)

	for _, feed := range data.Feeds {
		wg.Add(1)
		go loadTinyLogContent(feed, chFeedContent, chFeedError, &wg)
	}

	wg.Wait()

	var tlfi []*TlFeedItem
	for i := 0; i < numFeed; i++ {
		e := <-chFeedError
		if e != nil {
			<-chFeedContent
			log.Println(e)
		} else {
			rf := <-chFeedContent
			feedItems, err := parseTinyLogContent(rf)
			if err != nil {
				log.Println(err)
			} else {
				tlfi = append(tlfi, feedItems...)
			}
		}
	}

	s := TlStream{
		Name:  "main",
		Items: tlfi,
	}
	sort.Sort(&s)
	return &s, nil
}

// Load tinylog page from Feed URL
func loadTinyLogContent(feed TlFeed, chFeedContent chan TlRawFeed, chFeedError chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("Retrieving content from ", feed.Link)

	// Fallback title if not within tinylog response page.
	rf := TlRawFeed{Name: feed.Title}

	gemclient := &gemini.Client{}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)

	response, err := gemclient.Get(ctx, feed.Link)
	if err != nil {
		chFeedError <- fmt.Errorf("Error retrieving content from %v", feed.Link)
		chFeedContent <- rf
		return
	}
	defer response.Body.Close()

	// TODO: Add an option to accept gemini feeds with expired certificate.
	// TODO: Add possibility to validate certs?
	if respCert := response.TLS().PeerCertificates; len(respCert) > 0 && time.Now().After(respCert[0].NotAfter) {
		chFeedError <- fmt.Errorf("Invalid certificate for capsule", feed.Link, " caspule is ignored.")
		chFeedContent <- rf
		return
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		chFeedError <- fmt.Errorf("Couldn't read response from tinylogs", feed.Link, ", ignoring feed.")
		chFeedContent <- rf
		return
	}

	rf.Content = string(content)
	chFeedContent <- rf
	chFeedError <- nil

	log.Println("Tiny log retrieved from", feed.Link)
	return
}

// Parse gemini content of the tinylog file.
func parseTinyLogContent(rawFeed TlRawFeed) ([]*TlFeedItem, error) {
	author := rawFeed.Name

	entries := strings.Split(rawFeed.Content, "\n\n")
	nbEntries := len(entries)

	var fi []*TlFeedItem

	if nbEntries < 1 {
		return fi, fmt.Errorf("Invalid tinylog format")
	}

	if nbEntries > 1 {
		foundMeta := false
		for i := 0; i < nbEntries; i++ {
			l := strings.TrimSpace(entries[i])
			if strings.HasPrefix(l, "## ") {
				f, e := parseTinyLogItem(l, author)
				if e != nil {
					// Ignoring the entry but continuing in case other entries of this feed are in a known format.
					log.Println(e)
				} else {
					fi = append(fi, &f)
				}
			} else if foundMeta == false {
				a := parseTinyLogHeaderForAuthor(entries[i])
				if a != "" {
					author = a
					foundMeta = true
				} else {
					log.Println("Ignoring malformed entry", author, l)
				}
			} else {
				log.Println("Ignoring malformed entry", author, l)
			}
		}
	}

	return fi, nil
}

// Parse tinylog Header.
func parseTinyLogHeaderForAuthor(header string) string {
	var metaAuthor, metaAvatar string

	lines := strings.Split(header, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "author:") {
			metaAuthor = strings.TrimSpace(line[len("author:"):])
		} else if strings.HasPrefix(line, "avatar:") {
			// TODO: If avatar is more than 1 emoji, cut.
			metaAvatar = strings.TrimSpace(line[len("avatar:"):])
			if n := strings.Split(metaAvatar, " "); len(n) > 1 {
				metaAvatar = n[0]
			}
		}
	}

	if metaAuthor != "" {
		return strings.TrimSpace(metaAvatar + " " + metaAuthor)
	}

	return ""
}

// Parse tinylog Item.
func parseTinyLogItem(content string, author string) (TlFeedItem, error) {
	ft := TlFeedItem{Author: author}

	lines := strings.Split(content, "\n")

	if len(lines) < 2 {
		return ft, fmt.Errorf("Ignoring malformed entry", author, content)
	}

	pubDate, err := parseTinyLogItemForDate(lines[0])
	if err != nil {
		return ft, err
	}
	entry := strings.Join(lines[1:], "\n")

	ft.Content = strings.TrimSpace(entry)
	ft.Published = pubDate

	return ft, nil
}

// Get date from entry.
func parseTinyLogItemForDate(content string) (time.Time, error) {
	stringDate := content[3:]
	var date time.Time

	valid := false
	for _, format := range supportedTimeFormat {
		d, e := time.Parse(format, stringDate)
		if e == nil {
			valid = true
			date = d
			break
		}
	}

	if valid == false {
		return time.Time{}, fmt.Errorf("No date format found for this entry: %v", stringDate)
	}

	return date, nil
}
