package core

import (
  "time"
)

type TlData struct {
  Feeds map[string] TlFeed
  Config *TlConfig
  Stream []*TlFeedItem
}

type TlConfig struct {
  Subscribed_data string
  Refresh int
}

type TlFeed struct {
  Title string
  Link string
}

type TlFeedItem struct {
  Author string
  Content string
  Published time.Time
}

