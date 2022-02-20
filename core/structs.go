package core

import (
	"time"
)

const (
	FeedValid       = 1
	FeedUnreachable = 2
	FeedWrongFormat = 3
	FeedSSLError    = 4
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
	// Global config:
	Subscribed_data string
	Date_format     string
	Log_file        string
	Highlights      string
	Mode            string
	// Tinylog file edition:
	Allow_edit        bool
	Tinylog_path      string
	Post_edit_script  string
	Post_edit_refresh bool
	// CLI options:
	Cli_limit int
	// TUI options:
	Tui_status_emoji        bool
	Tui_max_entries         int
	Tui_copy_stub_clipboard bool
	Tui_show_stub           bool
	// TUI theme:
	Tui_color_background                 string
	Tui_color_links                      string
	Tui_color_text                       string
	Tui_color_focus_box                  string
	Tui_color_author_name                string
	Tui_color_elapsed_time               string
	Tui_color_box                        string
	Tui_color_selected_background        string
	Tui_color_selected_foreground        string
	Tui_color_highlight                  string
	Tui_color_quote                      string
	Tui_color_button_text                string
	Tui_color_button_selected_background string
	Tui_color_button_background          string
}

type TlFeed struct {
	Title       string
	Link        string
	DisplayName string
	Status      int
}

type TlFeedItem struct {
	Author    string
	Content   string
	Published time.Time
	Uri       string
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
