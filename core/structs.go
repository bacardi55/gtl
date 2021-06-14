package core

import (
	"time"
)

type TlData struct {
	Feeds  map[string]TlFeed
	Config *TlConfig
	Stream *TlStream
}

func (Data *TlData) RefreshFeeds() error {
	var err error
	Data.Stream, err = refreshStream(*Data)
	if err != nil {
		return err
	}

	return nil
}

type TlConfig struct {
	Subscribed_data string
	Date_format     string
	Log_file        string
	Highlights      string
	Cli_limit       int
	Mode            string
}

type TlFeed struct {
	Title       string
	Link        string
	DisplayName string
}

type TlFeedItem struct {
	Author    string
	Content   string
	Published time.Time
}

type TlStream struct {
	// Name could be used to manage multiple stream
	// Ex: All, Notification, …
	Name  string
	Items []*TlFeedItem
}

// Implement sort.Interface Len.
func (Stream *TlStream) Len() int {
	return len(Stream.Items)
}

// Implement Interface sort.Interface Less.
func (Stream *TlStream) Less(i, j int) bool {
	return Stream.Items[i].Published.After(Stream.Items[j].Published)
}

// Implement Interface sort.Interface Swap.
func (Stream *TlStream) Swap(i, j int) {
	Stream.Items[i], Stream.Items[j] = Stream.Items[j], Stream.Items[i]
}
