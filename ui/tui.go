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
	{"Fitler tinylog", "Enter/Left click", "Only display entries from this tinylog"},
	{"(Un)Mute tinylog", "Alt-Enter/Right click", "Hide entries from this tinylog"},
	{"Help", "?", "Toggle displaying this help."},
	{"Quit", "q / Ctrl-c", "Quit GTL."},
}

var TlTui TlTUI

func displayStreamTui(data *core.TlData) error {
	e := data.RefreshFeeds()
	if e != nil {
		return fmt.Errorf("Couldn't refresh feeds")
	}

	TlTui.InitApp(data.Config.Tui_status_emoji)
	TlTui.SetAppUI(data)
	TlTui.SetShortcuts()

  if data.Config.Allow_edit == true && data.Config.Tinylog_path != "" {
    if err := TlTui.InitTlEditor(data.Config.Tinylog_path, data.Config.Post_edit_script); err != nil {
      log.Println("Error while enabling tinylog edition:\n", err)
    } else {
      log.Println("Tinylog edition enabled.")
    }
  } else {
    log.Println("Tinylog edition not enabled.")
  }

	TlTui.RefreshStream = func(refresh bool) {
		if TlTui.Filter == "All Subscriptions" {
			TlTui.Filter = ""
			TlTui.Muted = []string{}
		}

		if refresh == true {
			e := data.RefreshFeeds()
			if e != nil {
				log.Fatalln("Couldn't refresh feeds")
			}
		}
		TlTui.ListTl = createListTl(data.Feeds)
		TlTui.SideBarBox.AddPanel("subscriptions", TlTui.ListTl, true, true)

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
	p.SetTitle(" Subscriptions: ")
	p.SetBorder(true)
	p.SetBorderColorFocused(tcell.ColorGreen)
	p.SetPadding(1, 1, 0, 0)

	TlTui.ListTl = createListTl(tl)
	p.AddPanel("subscriptions", TlTui.ListTl, true, true)
	return p
}

func contentBox(data *core.TlData) *cview.Panels {
	p := cview.NewPanels()
	p.SetBorder(true)
	p.SetBorderColorFocused(tcell.ColorGreen)
	p.SetTitle(createTimelineTitle(TlTui.LastRefresh, false))
	p.SetPadding(0, 0, 1, 0)

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
	list.SetSelectedAlwaysCentered(false)

	orderedTl := getOrderedSubscriptions(tl)
	sort.Sort(&orderedTl)

	list.AddContextItem("(Un)Mute tinylog", 'm', func(index int) {
		// index == 0 means All Subscriptions.
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

	i := createListItem("All Subscriptions", "> Press '?' for help")
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
	p.SetTitle(" Help: ")
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

func createRefreshBox() *cview.Panels {
	p := cview.NewPanels()
	p.SetBorder(true)
	p.SetTitle(" Refreshing stream: ")
	p.SetPadding(2, 0, 5, 0)

	tv := cview.NewTextView()
	tv.SetTextAlign(cview.AlignCenter)
	tv.SetText("Feeds are being refreshed, please wait…")

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
	if TlTui.Emoji == true {
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

func gemtextFormat(s string, isHighlighted bool) string {
	closeFormat := "[white::-]"
	if isHighlighted == true {
		closeFormat = "[white::b]"
	}

	// Format quotes:
	re := regexp.MustCompile("(?im)^(> .*[^\n])([\n]*)")
	if isHighlighted == true {
		s = re.ReplaceAllString(s, "[grey::bi] $1"+closeFormat+"$2")
	} else {
		s = re.ReplaceAllString(s, "[grey::i] $1"+closeFormat+"$2")
	}

	// Format links:
	re = regexp.MustCompile("(?im)^(=>)( [^\n]*[\n]*)")
	s = re.ReplaceAllString(s, "[skyblue::b]→$2"+closeFormat)

	// Format lists:
	re = regexp.MustCompile("(?im)^([*] [^\n]*)")
	if isHighlighted == true {
		s = re.ReplaceAllString(s, " [::bi]$1"+closeFormat)
	} else {
		s = re.ReplaceAllString(s, " [::i]$1"+closeFormat)
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

func createFormModal() *cview.Modal {
  m := cview.NewModal()

  // TODO: create accept function to validate entry format.
  f := m.GetForm()
  f.AddInputField("newEntryDate", time.Now().Format("2006-01-02 15:04 -0700"), 0, nil, nil)
  f.AddInputField("newEntryContent", "", 0, nil, nil)
  //input := f.GetFormItemByLabel("newEntryDate")
  //input.SetFieldNote("The content of the tinylog entry.")

  f.AddButton("Add", func() {
    log.Println("In ADD")
  })

  f.AddButton("Cancel", func() {
    log.Println("In Cancel")
    // TODO: clean form.
    toggleFormModal()
  })

  return m
}
