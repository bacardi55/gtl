package core

import (
	"log"
	"sort"
	"sync"
)

type TlRawFeed struct {
	Name    string
	Content string
	Status  int
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
		var rf TlRawFeed
		var displayName string

		e := <-chFeedError
		if e != nil {
			rf = <-chFeedContent
			displayName = rf.Name
			log.Println(e)
		} else {
			rf = <-chFeedContent
			dn, feedItems, err := parseTinyLogContent(rf)
			displayName = dn

			if err != nil {
				log.Println(err)
			} else {
				tlfi = append(tlfi, feedItems...)
			}
		}
		f := data.Feeds[rf.Name]
		if displayName != "" {
			f.DisplayName = displayName
		} else {
			f.DisplayName = rf.Name
		}
		f.Status = rf.Status
		data.Feeds[rf.Name] = f
	}

	s := TlStream{
		Name:  "main",
		Items: tlfi,
	}
	sort.Sort(&s)
	return &s, nil
}
