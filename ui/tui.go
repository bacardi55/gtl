package ui

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"

	"git.bacardi55.io/bacardi55/gtl/core"
)

var shortcuts = []TlShortcut{
	{"Navigate", "↑∕↓j/k", "Refresh timeline but keep active filters."},
	{"Refresh", "r", "Refresh timeline but keep active filters."},
	{"Timeline", "t", "Display timeline, remove active filters."},
	{"Highlights", "h", "Display only entries containing highlights, keep tinylog filters active."},
	{"Focus", "TAB", "Switch focus between the timeline and the subsciption list."},
	{"Sidebar toggle", "s", "Hide/Show TinyLogs sidebar."},
	{"Fitler tinylog", "Enter/Left click", "Only display entries from this tinylog."},
	{"(Un)Mute tinylog", "Alt-Enter/Right click", "Hide entries from this tinylog."},
	{"Select tinylog entry", "J/K", "Select prev/next entries."},
	{"New tinylog entry", "N", "Open tinylog editor with optional stub in clipboard."},
	{"Reply to tinylog entry", "R", "On a selected entry, open tinylog editor with optional stub in clipboard."},
	{"Open link(s) in tinylog entry", "O", "On a selected entry, open link(s) in browser."},
	{"Open thread", "T", "If the selected entry is a response to another tinylog entry, will open the original entry in a popup."},
	{"Open action modal", "Alt-Enter", "Open modal action for a selected entry"},
	{"Search", "/", "Search entries containing a specific text. Search will keep the active filters ((un)muted tinylog(s), highlights) to only search on already filtered entries"},
	{"Help", "?", "Toggle displaying this help."},
	{"Quit", "q / Ctrl-c", "Quit GTL."},
}

var TlTui TlTUI

func displayStreamTui(data *core.TlData) error {
	e := data.RefreshFeeds()
	if e != nil {
		return fmt.Errorf("Couldn't refresh feeds")
	}

	TlTui.InitApp()
	TlTui.SetAppUI(data)
	TlTui.SetShortcuts()

	if data.Config.Allow_edit == true && data.Config.Tinylog_path != "" {
		if err := TlTui.InitTlEditor(data.Config.Tinylog_path, data.Config.Post_edit_script, data.Config.Post_edit_refresh); err != nil {
			log.Println("Error while enabling tinylog edition:\n", err)
		} else {
			log.Println("Tinylog edition enabled.")
		}
	} else {
		log.Println("Tinylog edition not enabled.")
	}

	TlTui.RefreshStream = func(refresh bool) {
		if TlTui.Filter == "All TinyLogs" {
			TlTui.Filter = ""
			TlTui.Muted = []string{}
		}

		if refresh == true {
			e := data.RefreshFeeds()
			if e != nil {
				log.Fatalln("Couldn't refresh TinyLogs")
			}
			TlTui.TlStream = data.Stream
		}
		TlTui.ListTl = createListTl(data)
		TlTui.SideBarBox.AddPanel("subscriptions", TlTui.ListTl, true, true)

		TlTui.TimelineTV = getContentTextView(data)
		TlTui.ContentBox.AddPanel("timeline", TlTui.TimelineTV, true, true)

		// Needs to happen after the getContentTextView function for displaying
		// a seperator between new and old entries.
		if refresh == true {
			TlTui.LastRefresh = time.Now()
		}
		// Needs to happen after TlTui.LastRefresh is updated.
		TlTui.ContentBox.SetTitle(createTimelineTitle(TlTui.LastRefresh, TlTui.FilterHighlights, TlTui.Filter, data.Config))
	}

	TlTui.App.SetRoot(TlTui.Layout, true)
	if err := TlTui.App.Run(); err != nil {
		panic(err)
	}

	return nil
}

func sideBarBox(data *core.TlData) *cview.Panels {
	p := cview.NewPanels()
	p.SetTitle(" [::u]Subscribed Authors[::-]: ")
	p.SetBorder(true)
	p.SetPadding(1, 1, 0, 0)

	// Box border color, default to white:
	p.SetBorderColor(tcell.ColorWhite.TrueColor())
	// If in config:
	if data.Config.Tui_color_box != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_box, 16, 32)
		if e != nil {
			log.Println("Focus color isn't valid (tui_color_box)")
		} else {
			p.SetBorderColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Focus color, default to green:
	p.SetBorderColorFocused(tcell.ColorGreen.TrueColor())
	// If in config:
	if data.Config.Tui_color_focus_box != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_focus_box, 16, 32)
		if e != nil {
			log.Println("Focus color isn't valid (tui_color_focus_box)")
		} else {
			p.SetBorderColorFocused(tcell.NewHexColor(int32(h)))
		}
	}

	// Background color:
	if data.Config.Tui_color_background != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_background, 16, 32)
		if e != nil {
			log.Println("Background color isn't valid (tui_color_background)")
		} else {
			p.SetBackgroundColor(tcell.NewHexColor(int32(h)))
		}
	}

	// DefaultColor, default is white (use default text color):
	if data.Config.Tui_color_text != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_text, 16, 32)
		if e != nil {
			log.Println("Text color isn't valid (tui_color_text)")
		} else {
			p.SetTitleColor(tcell.NewHexColor(int32(h)))
		}
	}

	TlTui.ListTl = createListTl(data)
	p.AddPanel("subscriptions", TlTui.ListTl, true, true)
	return p
}

func contentBox(data *core.TlData) *cview.Panels {
	p := cview.NewPanels()
	p.SetBorder(true)
	p.SetTitle(createTimelineTitle(TlTui.LastRefresh, false, "", data.Config))
	p.SetPadding(0, 0, 1, 0)

	// Box border color, default to white:
	p.SetBorderColor(tcell.ColorWhite.TrueColor())
	// If in config:
	if data.Config.Tui_color_box != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_box, 16, 32)
		if e != nil {
			log.Println("Focus color isn't valid (tui_color_box)")
		} else {
			p.SetBorderColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Focus color, default to green:
	p.SetBorderColorFocused(tcell.ColorGreen.TrueColor())
	// If in config:
	if data.Config.Tui_color_focus_box != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_focus_box, 16, 32)
		if e != nil {
			log.Println("Focus color isn't valid (tui_color_focus_box)")
		} else {
			p.SetBorderColorFocused(tcell.NewHexColor(int32(h)))
		}
	}

	// Background color:
	if data.Config.Tui_color_background != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_background, 16, 32)
		if e != nil {
			log.Println("Background color isn't valid (tui_color_background)")
		} else {
			p.SetBackgroundColor(tcell.NewHexColor(int32(h)))
		}
	}

	TlTui.TimelineTV = getContentTextView(data)
	p.AddPanel("timeline", TlTui.TimelineTV, true, true)

	return p
}

func getContentTextView(data *core.TlData) *cview.TextView {
	var content string
	t := time.Now()
	separator := false
	nbEntries := 0
	for _, i := range data.Stream.Items {
		// If a limit is set and has been reached.
		// In case of search, we don't care about max entries limit:
		if len(TlTui.FilterSearch) == 0 && (data.Config.Tui_max_entries > 0 && nbEntries >= data.Config.Tui_max_entries) {
			break
		}

		// If not the wrong author to filter on.
		if TlTui.Filter != "" && TlTui.Filter != i.Author {
			continue
		}

		// If this author is muted
		if len(TlTui.Muted) > 0 {
			if found, _ := isMuted(i.Author); found == true {
				continue
			}
		}

		// Highlight color, default is red (ff0000):
		hightlightColor := "ff0000"
		if data.Config.Tui_color_highlight != "" {
			h, e := strconv.ParseInt(data.Config.Tui_color_highlight, 16, 32)
			if e != nil || cview.ColorHex(tcell.NewHexColor(int32(h))) == "" {
				log.Println("Author name color isn't valid (Tui_color_highlight)")
			} else {
				hightlightColor = data.Config.Tui_color_highlight
			}
		}

		f := false
		if len(data.Config.Highlights) > 0 {
			if highlights := strings.Split(data.Config.Highlights, ","); len(highlights) > 0 {
				i.Content, f = hightlightContent(i.Content, highlights, hightlightColor)
			}
		}

		// Active search:
		if len(TlTui.FilterSearch) > 0 {
			highlights := make([]string, 1)
			highlights[0] = TlTui.FilterSearch
			i.Content, f = hightlightContent(i.Content, highlights, hightlightColor)

			// If active search and not found in entry, ignore it:
			if f == false {
				continue
			}
		}

		// Highlights only:
		var c string
		if TlTui.FilterHighlights == true && f == true {
			// No bold because all would be bold.
			c = gemtextFormat(i.Content, false, TlTui.TlConfig.Tui_status_emoji, data.Config)
		} else if TlTui.FilterHighlights == false {
			c = gemtextFormat(i.Content, f, TlTui.TlConfig.Tui_status_emoji, data.Config)
			if f == true {
				c = "[:-:b]" + c + "[:-:-]"
			}
		} else {
			continue
		}

		// Entry not ignored:
		// Author color, default is red (ff0000):
		authorColor := "ff0000"
		if data.Config.Tui_color_author_name != "" {
			h, e := strconv.ParseInt(data.Config.Tui_color_author_name, 16, 32)
			if e != nil || cview.ColorHex(tcell.NewHexColor(int32(h))) == "" {
				log.Println("Author name color isn't valid (Tui_color_author_name)")
			} else {
				authorColor = data.Config.Tui_color_author_name
			}
		}
		a := fmt.Sprintf("[#" + authorColor + "]" + i.Author + "[-::]")

		// ElapsedColor, default is Skyblue (87CEEB):
		elapsedColor := "87ceeb"
		if data.Config.Tui_color_elapsed_time != "" {
			h, e := strconv.ParseInt(data.Config.Tui_color_elapsed_time, 16, 32)
			if e != nil || cview.ColorHex(tcell.NewHexColor(int32(h))) == "" {
				log.Println("Elapsed time color isn't valid (Tui_color_elapsed_time)")
			} else {
				elapsedColor = data.Config.Tui_color_elapsed_time
			}
		}
		d := "[#" + elapsedColor + "::]" + formatElapsedTime(t.Sub(i.Published)) + "[-::]"

		// Separator:
		if isTlEntryNew(i, TlTui.LastRefresh) != true && separator != true {
			if nbEntries > 0 {
				content = content + "            --------------------------- \n"
			}
			separator = true
		}

		content = content + fmt.Sprintf("\n[\"entry-"+strconv.Itoa(nbEntries)+"\"]%v - %v\n%v\n%v\n", d, i.Published.Format(data.Config.Date_format), a, c)
		nbEntries++
	}

	// Clean search filter after displaying search results.
	// This will avoid always having previous search still
	// highlighting text:
	TlTui.FilterSearch = ""

	tv := cview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetRegions(true)
	tv.SetToggleHighlights(false)
	tv.SetText(content)

	// Selected background color, default white:
	if data.Config.Tui_color_selected_background != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_selected_background, 16, 32)
		if e != nil {
			log.Println("Selected background color isn't valid (Tui_color_selected_background)")
		} else {
			tv.SetHighlightBackgroundColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Selected foreground color, default white:
	if data.Config.Tui_color_selected_foreground != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_selected_foreground, 16, 32)
		if e != nil {
			log.Println("Selected foreground color isn't valid (Tui_color_selected_foreground)")
		} else {
			tv.SetHighlightForegroundColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Background color:
	if data.Config.Tui_color_background != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_background, 16, 32)
		if e != nil {
			log.Println("Background color isn't valid (tui_color_background)")
		} else {
			tv.SetBackgroundColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Default text color:
	if data.Config.Tui_color_text != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_text, 16, 32)
		if e != nil {
			log.Println("Text color isn't valid (tui_color_text)")
		} else {
			tv.SetTextColor(tcell.NewHexColor(int32(h)))
		}
	}

	return tv
}

func isTlEntryNew(tlfi *core.TlFeedItem, lastRefresh time.Time) bool {
	return tlfi.Published.After(lastRefresh)
}

func createListTl(data *core.TlData) *cview.List {
	tl := data.Feeds

	list := createList("", false)
	list.ShowSecondaryText(true)
	list.SetSelectedAlwaysCentered(false)

	// Author color, default is red:
	list.SetMainTextColor(tcell.ColorRed.TrueColor())
	if data.Config.Tui_color_author_name != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_author_name, 16, 32)
		if e != nil {
			log.Println("Author name color isn't valid (tui_color_author_name)")
		} else {
			list.SetMainTextColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Link color, default is skyblue (87CEEB):
	list.SetSecondaryTextColor(tcell.ColorSkyblue.TrueColor())
	if data.Config.Tui_color_links != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_links, 16, 32)
		if e != nil {
			log.Println("Link color isn't valid (tui_color_links)")
		} else {
			list.SetSecondaryTextColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Background color:
	if data.Config.Tui_color_background != "" {
		h, e := strconv.ParseInt(data.Config.Tui_color_background, 16, 32)
		if e != nil {
			log.Println("Background color isn't valid (tui_color_background)")
		} else {
			list.SetBackgroundColor(tcell.NewHexColor(int32(h)))
		}
	}

	orderedTl := getOrderedSubscriptions(tl)
	sort.Sort(&orderedTl)

	list.AddContextItem("(Un)Mute TinyLog", 'm', func(index int) {
		// index == 0 means All Feeds.
		if index > 0 {
			author := orderedTl.Items[index-1].DisplayName
			found, foundIndex := isMuted(author)

			if found == false {
				TlTui.Muted = append(TlTui.Muted, author)
			} else if found == true && foundIndex >= 0 {
				TlTui.Muted[foundIndex] = TlTui.Muted[len(TlTui.Muted)-1]
				TlTui.Muted = TlTui.Muted[:len(TlTui.Muted)-1]
			}
			TlTui.RefreshStream(false)
		}
	})

	i := createListItem("All TinyLogs", "> Press '?' for help")
	i.SetSelectedFunc(func() {
		TlTui.Filter = TlTui.ListTl.GetCurrentItem().GetMainText()
		TlTui.RefreshStream(false)
	})
	list.AddItem(i)

	for _, f := range orderedTl.Items {
		t := getStatusIcon(f) + " - " + f.DisplayName
		it := createListItem(t, "=> "+f.Link)
		list.AddItem(it)
		it.SetSelectedFunc(func() {
			tmp := strings.Split(TlTui.ListTl.GetCurrentItem().GetMainText(), "-")
			if len(tmp) < 2 {
				TlTui.Filter = ""
			} else {
				TlTui.Filter = strings.TrimSpace(tmp[1])
			}
			TlTui.RefreshStream(false)
		})
	}

	return list
}

func getOrderedSubscriptions(tl map[string]core.TlFeed) TlTuiSubs {
	var sub TlTuiSubs
	for _, t := range tl {
		sub.Items = append(sub.Items, t)
	}

	return sub
}

func createListItem(title string, subtitle string) *cview.ListItem {
	item := cview.NewListItem(title)

	if subtitle != "" {
		item.SetSecondaryText(subtitle)
	}

	return item
}

func createList(title string, border bool) *cview.List {
	list := cview.NewList()
	if title != "" {
		list.SetTitle(title)
	}
	list.SetBorder(border)
	list.SetWrapAround(true)

	return list
}

func createTimelineTitle(t time.Time, highlights bool, filter string, config *core.TlConfig) string {
	// DefaultColor, default is white (use default text color):
	defaultColor := "ffffff"
	if config.Tui_color_text != "" {
		h, e := strconv.ParseInt(config.Tui_color_text, 16, 32)
		if e != nil || cview.ColorHex(tcell.NewHexColor(int32(h))) == "" {
			log.Println("Text color isn't valid (Tui_color_text)")
		} else {
			defaultColor = config.Tui_color_text
		}
	}

	// Author color, default is red (ff0000):
	authorColor := "ff0000"
	if config.Tui_color_author_name != "" {
		h, e := strconv.ParseInt(config.Tui_color_author_name, 16, 32)
		if e != nil || cview.ColorHex(tcell.NewHexColor(int32(h))) == "" {
			log.Println("Author name color isn't valid (Tui_color_author_name)")
		} else {
			authorColor = config.Tui_color_author_name
		}
	}

	start := ""

	if highlights == true {
		start = start + "[#" + defaultColor + "::bu]Highlights[::-]"
	} else {
		start = start + "[#" + defaultColor + "::bu]Timeline[::-]"
	}

	if filter != "" {
		start = start + " from [#" + authorColor + "::]" + filter + "[#" + defaultColor + "::]"
		return fmt.Sprintf("  %v  ", start)
	} else {
		return fmt.Sprintf("  %v - [#"+defaultColor+"::i]Refreshed at %v[::-]  ", start, t.Format("15:04 MST"))
	}
}

func createHelpBox(config *core.TlConfig) *cview.Panels {
	p := cview.NewPanels()
	p.SetBorder(true)
	p.SetTitle(" [::u]Help[::-]: ")
	p.SetPadding(2, 0, 5, 0)

	helpTable := cview.NewTable()
	helpTable.SetBorders(false)
	helpTable.SetFixed(1, 0)
	helpTable.SetSelectable(false, false)
	helpTable.SetSortClicked(false)

	// Background color, default is black:
	p.SetBackgroundColor(tcell.ColorBlack.TrueColor())
	helpTable.SetBackgroundColor(tcell.ColorBlack.TrueColor())
	if config.Tui_color_background != "" {
		h, e := strconv.ParseInt(config.Tui_color_background, 16, 32)
		if e != nil {
			log.Println("Background color isn't valid (tui_color_background)")
		} else {
			p.SetBackgroundColor(tcell.NewHexColor(int32(h)))
			helpTable.SetBackgroundColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Box border color, default to green:
	p.SetBorderColor(tcell.ColorGreen.TrueColor())
	p.SetBorderColorFocused(tcell.ColorGreen.TrueColor())
	// If in config:
	if config.Tui_color_focus_box != "" {
		h, e := strconv.ParseInt(config.Tui_color_focus_box, 16, 32)
		if e != nil {
			log.Println("Focus color isn't valid (tui_color_focus_box)")
		} else {
			p.SetBorderColor(tcell.NewHexColor(int32(h)))
			p.SetBorderColorFocused(tcell.NewHexColor(int32(h)))
		}
	}

	// 3 rows: Name, Command, Description.
	c := getNewTableCell("Name", config)
	c.SetAttributes(tcell.AttrBold | tcell.AttrUnderline)
	helpTable.SetCell(0, 0, c)

	c = getNewTableCell("Command", config)
	c.SetAttributes(tcell.AttrBold | tcell.AttrUnderline)
	helpTable.SetCell(0, 1, c)

	c = getNewTableCell("Description", config)
	c.SetAttributes(tcell.AttrBold | tcell.AttrUnderline)
	helpTable.SetCell(0, 2, c)

	for i, shortcut := range shortcuts {
		tc := getNewTableCell(shortcut.Name, config)
		helpTable.SetCell(i+1, 0, tc)

		tc = getNewTableCell(shortcut.Command, config)
		tc.SetAttributes(tcell.AttrBold)
		helpTable.SetCell(i+1, 1, tc)

		tc = getNewTableCell(shortcut.Description, config)
		helpTable.SetCell(i+1, 2, tc)
	}

	helpTable.InsertRow(1)

	p.AddPanel("help", helpTable, true, true)
	return p
}

func createRefreshBox(config *core.TlConfig) *cview.Panels {
	p := cview.NewPanels()
	p.SetBorder(true)
	p.SetTitle(" [::bu]Refreshing stream[::-]: ")
	p.SetPadding(2, 0, 5, 0)

	tv := cview.NewTextView()
	tv.SetTextAlign(cview.AlignCenter)
	tv.SetText("TinyLogs are being refreshed, please wait…")

	// Background color, default is black:
	tv.SetBackgroundColor(tcell.ColorBlack.TrueColor())
	p.SetBackgroundColor(tcell.ColorBlack.TrueColor())
	if config.Tui_color_background != "" {
		h, e := strconv.ParseInt(config.Tui_color_background, 16, 32)
		if e != nil {
			log.Println("Background color isn't valid (tui_color_background)")
		} else {
			tv.SetBackgroundColor(tcell.NewHexColor(int32(h)))
			p.SetBackgroundColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Text color, default is white:
	tv.SetTextColor(tcell.ColorWhite.TrueColor())
	p.SetTitleColor(tcell.ColorWhite.TrueColor())
	if config.Tui_color_text != "" {
		h, e := strconv.ParseInt(config.Tui_color_text, 16, 32)
		if e != nil {
			log.Println("Text color isn't valid (tui_color_text)")
		} else {
			tv.SetTextColor(tcell.NewHexColor(int32(h)))
			p.SetTitleColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Border color, default is green:
	p.SetBorderColor(tcell.ColorGreen.TrueColor())
	p.SetBorderColorFocused(tcell.ColorGreen.TrueColor())
	if config.Tui_color_focus_box != "" {
		h, e := strconv.ParseInt(config.Tui_color_focus_box, 16, 32)
		if e != nil {
			log.Println("Text color isn't valid (tui_color_text)")
		} else {
			p.SetBorderColor(tcell.NewHexColor(int32(h)))
			p.SetTitleColor(tcell.NewHexColor(int32(h)))
		}
	}

	p.AddPanel("refreshing", tv, true, true)
	TlTui.SideBarBox.AddPanel("subscriptions", TlTui.ListTl, true, true)

	return p
}

func getHelpContent(field string) string {
	text := field + ":\n\n"

	for _, s := range shortcuts {
		switch field {
		case "Names":
			text = text + s.Name + "\n"
		case "Descriptions":
			text = text + s.Description + "\n"
		case "Commands":
			text = text + s.Command + "\n"
		}
	}

	return text
}

func getStatusIcon(f core.TlFeed) string {
	status := f.Status
	r := ""
	if TlTui.TlConfig.Tui_status_emoji == true {
		if status == core.FeedValid {
			r = "✔"
		} else if status == core.FeedUnreachable {
			r = "💀"
		} else if status == core.FeedWrongFormat {
			r = "❌"
		} else if status == core.FeedSSLError {
			r = "🔓"
		}
		if f.DisplayName == TlTui.Filter {
			r = "🔎" + r
		}
		if f, _ := isMuted(f.DisplayName); f == true {
			r = "🔕" + r
		}
	} else {
		if status == core.FeedValid {
			r = "V"
		} else if status == core.FeedUnreachable {
			r = "D"
		} else if status == core.FeedWrongFormat {
			r = "X"
		} else if status == core.FeedSSLError {
			r = "S"
		}
		if f.DisplayName == TlTui.Filter {
			r = "F" + r
		}
		if f, _ := isMuted(f.DisplayName); f == true {
			r = "M" + r
		}
	}

	return r
}

func gemtextFormat(s string, isHighlighted bool, emoji bool, config *core.TlConfig) string {
	closeFormat := "[-:-:-]"
	if isHighlighted == true {
		closeFormat = "[-:-:b]"
	}

	// Format quotes:
	re := regexp.MustCompile("(?im)^(> .*)($)")
	// Quoted text, default color is grey (808080):
	quoteColor := "808080"
	if config.Tui_color_quote != "" {
		h, e := strconv.ParseInt(config.Tui_color_quote, 16, 32)
		if e != nil || cview.ColorHex(tcell.NewHexColor(int32(h))) == "" {
			log.Println("Link color isn't valid (Tui_color_quote)")
		} else {
			quoteColor = config.Tui_color_quote
		}
	}
	if isHighlighted == true {
		s = re.ReplaceAllString(s, "[#"+quoteColor+":-:bi] $1"+closeFormat+"$2")
	} else {
		s = re.ReplaceAllString(s, "[#"+quoteColor+":-:i] $1"+closeFormat+"$2")
	}

	// Format links, default color is skyblue (87CEEB):
	linkColor := "87ceeb"
	if config.Tui_color_links != "" {
		h, e := strconv.ParseInt(config.Tui_color_elapsed_time, 16, 32)
		if e != nil || cview.ColorHex(tcell.NewHexColor(int32(h))) == "" {
			log.Println("Link color isn't valid (Tui_color_links)")
		} else {
			linkColor = config.Tui_color_links
		}
	}
	re = regexp.MustCompile("(?im)^(=>)( [^\n]*[\n]*)")
	s = re.ReplaceAllString(s, "[#"+linkColor+":-:b]→$2"+closeFormat)

	// Format responses:
	// Must be after link format.
	re = regexp.MustCompile(`(?im)^(\[skyblue:-:b\]→ [^ ]* ){0,1}(re: )(.*)$`)
	startFormat := ""
	if emoji == true {
		startFormat = "💬 "
	} else {
		startFormat = "↳ "
	}
	startFormat = startFormat + "$3 $1"
	s = re.ReplaceAllString(s, startFormat+closeFormat)

	// Format lists:
	re = regexp.MustCompile("(?im)^([*] .*$)")
	if isHighlighted == true {
		s = re.ReplaceAllString(s, "  [-:-:bd]$1"+closeFormat)
	} else {
		s = re.ReplaceAllString(s, "  [-:-:d]$1"+closeFormat)
	}

	// Format headers
	re = regexp.MustCompile("(?im)^(####* )(.*$)")
	// If highlighted, already bold anyway.
	if isHighlighted == false {
		s = re.ReplaceAllString(s, "[grey:-:b]$1[-:-:u]$2"+closeFormat)
	} else {
		s = re.ReplaceAllString(s, "[grey:-:b]$1[-:-:bu]$2"+closeFormat)
	}

	return s
}

func isMuted(author string) (bool, int) {
	found := false
	foundIndex := -1

	for i, m := range TlTui.Muted {
		if author == m {
			found = true
			foundIndex = i
			break
		}
	}

	return found, foundIndex
}

func createFormModal(config *core.TlConfig) *cview.Modal {
	m := cview.NewModal()

	// Background color, default is black:
	m.SetBackgroundColor(tcell.ColorBlack.TrueColor())
	if config.Tui_color_background != "" {
		h, e := strconv.ParseInt(config.Tui_color_background, 16, 32)
		if e != nil {
			log.Println("Background color isn't valid (tui_color_background)")
		} else {
			m.SetBackgroundColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Border color, default is white:
	frame := m.GetFrame()
	frame.SetBorderColorFocused(tcell.ColorWhite.TrueColor())
	frame.SetBorderColor(tcell.ColorWhite.TrueColor())
	if config.Tui_color_focus_box != "" {
		h, e := strconv.ParseInt(config.Tui_color_focus_box, 16, 32)
		if e != nil {
			log.Println("Focus box color isn't valid (Tui_color_focus_box)")
		} else {
			frame.SetBorderColorFocused(tcell.NewHexColor(int32(h)))
			frame.SetBorderColor(tcell.NewHexColor(int32(h)))
		}
	}

	// Default text color:
	if config.Tui_color_text != "" {
		h, e := strconv.ParseInt(config.Tui_color_text, 16, 32)
		if e != nil {
			log.Println("Background color isn't valid (tui_color_text)")
		} else {
			m.SetTextColor(tcell.NewHexColor(int32(h)))
		}
	}

	return m
}

func updateFormModalContent(message string, mainButtonName string, buttonName string, execFunc func()) {
	m := TlTui.FormModal
	f := m.GetForm()
	f.Clear(true)

	// To override this, use TlTui.FormModal.SetTextAlign() after calling this function.
	m.SetTextAlign(cview.AlignCenter)

	m.SetText(message)

	if buttonName != "" {
		f.AddButton(buttonName, execFunc)
	}

	f.AddButton(mainButtonName, func() {
		toggleFormModal()
	})
}

func toggleFormModal() {
	if TlTui.DisplayFormModal == false {
		TlTui.DisplayFormModal = true
		TlTui.ContentBox.SendToBack("timeline")
		TlTui.ContentBox.SendToFront("newEntryModal")
		TlTui.ContentBox.ShowPanel("newEntryModal")
	} else {
		TlTui.DisplayFormModal = false
		TlTui.ContentBox.SendToFront("timeline")
		TlTui.ContentBox.SendToBack("newEntryModal")
		TlTui.ContentBox.HidePanel("newEntryModal")
	}
}

func createNewEntryStub(dateFormat string) string {
	stub := "## " + time.Now().Format(dateFormat) + "\n\n\n"
	return stub
}

func createResponseStub(tlfi *core.TlFeedItem, dateFormat string) string {
	a := tlfi.Author
	if strings.Contains(tlfi.Author, " ") {
		a = strings.Split(tlfi.Author, " ")[1]
	}

	stub := "## " + time.Now().Format(dateFormat) + "\nRE: " + a + " " + tlfi.Published.Format(dateFormat) + "\n"
	for _, l := range strings.Split(tlfi.Content, "\n") {
		// Ignore separator line:
		if strings.Contains(l, "           ---------------------------") {
			continue
		}
		stub += "> " + l + "\n"
	}
	stub += "\n"
	return stub
}

func copyToClipboard(content string) error {
	return clipboard.WriteAll(content)
}

func getNewTableCell(title string, config *core.TlConfig) *cview.TableCell {
	c := cview.NewTableCell(title)
	// Text color, default to white:
	c.SetTextColor(tcell.ColorWhite.TrueColor())
	// If in config:
	if config.Tui_color_text != "" {
		h, e := strconv.ParseInt(config.Tui_color_text, 16, 32)
		if e != nil {
			log.Println("Focus color isn't valid (tui_color_text)")
		} else {
			c.SetTextColor(tcell.NewHexColor(int32(h)))
		}
	}

	return c
}

func hightlightContent(content string, highlights []string, hightlightColor string) (string, bool) {
	f := false
	for _, h := range highlights {
		h = strings.TrimSpace(h)
		if strings.Contains(content, h) {
			content = strings.Replace(content, h, "[:#"+hightlightColor+":]"+h+"[:-:]", -1)
			f = true
		}
	}

	return content, f
}
