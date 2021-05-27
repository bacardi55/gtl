package core

import (
  "time"
)

type TlData struct {
  Feeds map[string] TlFeed
  Config *TlConfig
}

type TlConfig struct {
  Subscribed_data string
  Refresh int
}

type TlFeed struct {
  Title string
  Link string
  Items []*TlFeedItem
}

type TlFeedItem struct {
  Content string
  Published time.Time
}

