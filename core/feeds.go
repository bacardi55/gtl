package core

import (
	"fmt"
	"log"
	"context"
	"time"
	"io"
	"sync"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
)

// Refresh TlData.Stream.
func (Data *TlData) RefreshFeeds() error {
	var wg sync.WaitGroup

	numFeed := len(Data.Feeds)
	chFeedContent := make(chan string, numFeed)
	chFeedError := make(chan error, numFeed)

	for _, feed := range Data.Feeds {
		wg.Add(1)
		go loadTinyLogContent(feed, chFeedContent, chFeedError, &wg)
	}

	wg.Wait()

	for i := 0; i < numFeed; i++ {
		fmt.Println("i", i)
		e := <-chFeedError
		if e != nil {
			<-chFeedContent
			log.Println(e)
		} else {
			content := <-chFeedContent
			//fmt.Println("contents", content)
			feedItems, err := parseTinyLogContent(content)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println(feedItems)
			}
		}
	}

  return nil
}

// Load tinylog page from Feed URL
func loadTinyLogContent(feed TlFeed, chFeedContent chan string, chFeedError chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("Retrieving content from ", feed.Link)
	gemclient := &gemini.Client{}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)

	response, err := gemclient.Get(ctx, feed.Link)
	if err != nil {
		chFeedError <- fmt.Errorf("Error retrieving content from %v", feed.Link)
		chFeedContent <- ""
		return
	}
	defer response.Body.Close()

	// TODO: Add an option to accept gemini feeds with expired certificate.
	// TODO: Add possibility to validate certs?
	if respCert := response.TLS().PeerCertificates;
	(len(respCert) > 0 && time.Now().After(respCert[0].NotAfter)) {
		chFeedError <- fmt.Errorf("Invalid certificate for capsule ", feed.Link, " caspule is ignored.")
		chFeedContent <- ""
		return
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		chFeedError <- fmt.Errorf("Couldn't read response from tinylogs ", feed.Link, ", ignoring feed.")
		chFeedContent <- ""
		return
	}

	chFeedContent <- string(content)
	chFeedError <- nil

	log.Println("log retrieved from ", feed.Link)
	return
}

// Parse gemini content of the tinylog file.
func parseTinyLogContent(content string) ([]*TlFeedItem, error) {
	//fmt.Println(content)
	// TODO:
	var fi []*TlFeedItem

  lines := strings.Split(content, "\n")
	for _, line := range lines {
		log.Println(line)
  }

	return fi, nil
}

