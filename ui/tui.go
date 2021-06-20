package ui

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"

	"git.bacardi55.io/bacardi55/gtl/core"
)

var shortcuts = []TlShortcut{
	{"Refresh", "r", "Refresh timeline but keep active filters."},
	{"Timeline", "t", "Display timeline, remove active filters."},
	{"Highlights", "h", "Display only entries containing highlights, keep tinylog filters active."},
	{"Focus", "TAB", "Switch focus between the timeline and the subsciption list."},
	{"Sidebar toggle", "s", "Hide/Show Subscription sidebar."},
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

	TlTui.RefreshStream = func(refresh bool) {
		if TlTui.Filter == "All Subscriptions" {
			TlTui.Filter = ""
		}

		if refresh == true {
			e := data.RefreshFeeds()
			if e != nil {
				log.Fatalln("Couldn't refresh feeds")
			}
			TlTui.ListTl = createListTl(data.Feeds)
			TlTui.SideBarBox.AddPanel("subscriptions", TlTui.ListTl, true, true)
		}
		tv := getContentTextView(data)
		TlTui.ContentBox.SetTitle(createTimelineTitle(TlTui.LastRefresh, TlTui.FilterHighlights))
		TlTui.ContentBox.AddPanel("timeline", tv, true, true)

		// Needs to happen after the getContentTextView function for displaying
		// a seperator between new and old entries.
		if refresh == true {
			TlTui.LastRefresh = time.Now()
		}

		//tv = createFooterTextView(TlTui.LastRefresh, data.Config.Date_format)
		//TlTui.Footer.AddPanel("footer", tv, true, true)
	}

	TlTui.App.SetRoot(TlTui.Layout, true)
	if err := TlTui.App.Run(); err != nil {
		panic(err)
	}

	return nil
}

func sideBarBox(tl map[string]core.TlFeed) *cview.Panels {
	p := cview.NewPanels()
	p.SetTitle("Subscriptions:")
	p.SetBorder(true)
	p.SetBorderColorFocused(tcell.ColorGreen)

	TlTui.ListTl = createListTl(tl)
	p.AddPanel("subscriptions", TlTui.ListTl, true, true)
	return p
}

func contentBox(data *core.TlData) *cview.Panels {
	p := cview.NewPanels()
	p.SetBorder(true)
	p.SetBorderColorFocused(tcell.ColorGreen)
	p.SetTitle(createTimelineTitle(TlTui.LastRefresh, false))

	tv := getContentTextView(data)

	p.AddPanel("timeline", tv, true, true)
	return p
}

func getContentTextView(data *core.TlData) *cview.TextView {
	var content string
	t := time.Now()
	separator := false
	nbEntries := 0
	for _, i := range data.Stream.Items {
		if TlTui.Filter != "" && TlTui.Filter != i.Author {
			continue
		}
		f := false
		if len(data.Config.Highlights) > 0 {
			if highlights := strings.Split(data.Config.Highlights, ","); len(highlights) > 0 {
				for _, h := range highlights {
					h = strings.TrimSpace(h)
					if strings.Contains(i.Content, h) {
						i.Content = strings.Replace(i.Content, h, "[:red:]"+h+"[:black:]", -1)
						f = true
						break
					}
				}
			}
		}

		var c string
		ignoreEntry := false
		if TlTui.FilterHighlights == true && f == true {
			// No bold because all would be bold.
			c = gemtextFormat(i.Content, false)
		} else if TlTui.FilterHighlights == false {
			c = gemtextFormat(i.Content, f)
			if f == true {
				c = "[::b]" + c + "[::-]"
			}
		} else {
			ignoreEntry = true
		}

		if ignoreEntry != true {
			a := fmt.Sprintf("[red]" + i.Author + "[white::]")
			d := "[skyblue::]" + formatElapsedTime(t.Sub(i.Published)) + "[white::]"
			if isTlEntryNew(i, TlTui.LastRefresh) != true && separator != true {
				if nbEntries > 0 {
					content = content + "            --------------------------- \n"
				}
				separator = true
			}
			content = content + fmt.Sprintf("\n%v - %v\n%v\n%v\n", d, i.Published.Format(data.Config.Date_format), a, c)
			nbEntries++
		}
	}

	tv := cview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetText(content)

	return tv
}

func isTlEntryNew(tlfi *core.TlFeedItem, lastRefresh time.Time) bool {
	return tlfi.Published.After(lastRefresh)
}

func createListTl(tl map[string]core.TlFeed) *cview.List {
	list := createList("", false)
	list.ShowSecondaryText(true)
	list.SetMainTextColor(tcell.Color196)
	list.SetSecondaryTextColor(tcell.ColorSkyblue)

	i := createListItem("All Subscriptions", "> Press '?' for help")
	i.SetSelectedFunc(func() {
		TlTui.Filter = TlTui.ListTl.GetCurrentItem().GetMainText()
		TlTui.RefreshStream(false)
	})
	list.AddItem(i)
	orderedTl := getOrderedSubscriptions(tl)
	sort.Sort(&orderedTl)

	for _, f := range orderedTl.Items {
		t := getStatusIcon(f.Status) + " - " + f.DisplayName
		it := createListItem(t, "=> "+f.Link)
		list.AddItem(it)
		it.SetSelectedFunc(func() {
			TlTui.Filter = strings.TrimSpace(strings.Split(TlTui.ListTl.GetCurrentItem().GetMainText(), "-")[1])
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

func createFooter(latestRefresh time.Time, format string) *cview.Panels {
	p := cview.NewPanels()

	tv := createFooterTextView(latestRefresh, format)

	p.AddPanel("footer", tv, true, true)
	return p
}

func createFooterTextView(latestRefresh time.Time, format string) *cview.TextView {
	tv := cview.NewTextView()
	tv.SetMaxLines(1)
	tv.SetTextAlign(cview.AlignCenter)
	content := "Last Refresh:\t" + latestRefresh.Format(format)
	tv.SetText(content)
	return tv
}

func createTimelineTitle(t time.Time, highlights bool) string {
	if highlights == true {
		return fmt.Sprintf("  Highlights - Refreshed at %v  ", t.Format("15:04 MST"))
	} else {
		return fmt.Sprintf("  Timeline - Refreshed at %v  ", t.Format("15:04 MST"))
	}
}

func createHelpBox() *cview.Panels {
	p := cview.NewPanels()
	p.SetBorder(true)
	p.SetTitle("Help:")
	p.SetPadding(2, 0, 5, 0)

	helpTable := cview.NewTable()
	helpTable.SetBorders(false)
	helpTable.SetFixed(1, 0)
	helpTable.SetSelectable(false, false)
	helpTable.SetSortClicked(false)

	// 3 rows: Name, Command, Description.
	c := cview.NewTableCell("Name")
	c.SetAttributes(tcell.AttrBold | tcell.AttrUnderline)
	helpTable.SetCell(0, 0, c)

	c = cview.NewTableCell("Command")
	c.SetAttributes(tcell.AttrBold | tcell.AttrUnderline)
	helpTable.SetCell(0, 1, c)

	c = cview.NewTableCell("Description")
	c.SetAttributes(tcell.AttrBold | tcell.AttrUnderline)
	helpTable.SetCell(0, 2, c)

	for i, shortcut := range shortcuts {
		tc := cview.NewTableCell(shortcut.Name)
		helpTable.SetCell(i+1, 0, tc)

		tc = cview.NewTableCell(shortcut.Command)
		tc.SetAttributes(tcell.AttrBold)
		helpTable.SetCell(i+1, 1, tc)

		tc = cview.NewTableCell(shortcut.Description)
		helpTable.SetCell(i+1, 2, tc)
	}

	helpTable.InsertRow(1)

	p.AddPanel("help", helpTable, true, true)
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

func getStatusIcon(status int) string {
	if status == core.FeedValid {
		return "✔"
	} else if status == core.FeedUnreachable {
		return "💀"
	} else if status == core.FeedWrongFormat {
		return "❌"
	} else if status == core.FeedSSLError {
		return "🔓"
	}
	return ""
}

func gemtextFormat(s string, isHighlighted bool) string {
	closeFormat := "[white::-]"
	if isHighlighted == true {
		closeFormat = "[white::b]"
	}

	// Format quotes:
	re := regexp.MustCompile("(\n(> .*\n))")
	s = re.ReplaceAllString(s, "\n[grey::i] $2"+closeFormat)

	// Format links:
	re = regexp.MustCompile("([\n]*=> [^\n]*[\n]*)")
	s = re.ReplaceAllString(s, "[skyblue::b]$1"+closeFormat)

	// Format lists:
	re = regexp.MustCompile("([\n]{0,1})([*]{1} [^\n]*)")
	s = re.ReplaceAllString(s, "$1 $2")

	return s
}
