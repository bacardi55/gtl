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
		}
		TlTui.ListTl = createListTl(data.Feeds)
		TlTui.SideBarBox.AddPanel("subscriptions", TlTui.ListTl, true, true)

		TlTui.TimelineTV = getContentTextView(data)
		TlTui.ContentBox.AddPanel("timeline", TlTui.TimelineTV, true, true)

		// Needs to happen after the getContentTextView function for displaying
		// a seperator between new and old entries.
		if refresh == true {
			TlTui.LastRefresh = time.Now()
		}
		// Needs to happen after TlTui.LastRefresh is updated.
		TlTui.ContentBox.SetTitle(createTimelineTitle(TlTui.LastRefresh, TlTui.FilterHighlights, TlTui.Filter))
	}

	TlTui.App.SetRoot(TlTui.Layout, true)
	if err := TlTui.App.Run(); err != nil {
		panic(err)
	}

	return nil
}

func sideBarBox(tl map[string]core.TlFeed) *cview.Panels {
	p := cview.NewPanels()
	p.SetTitle(" [::u]Subscribed Authors[::-]: ")
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
	p.SetTitle(createTimelineTitle(TlTui.LastRefresh, false, ""))
	p.SetPadding(0, 0, 1, 0)

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
		if data.Config.Tui_max_entries > 0 && nbEntries >= data.Config.Tui_max_entries {
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

		f := false
		if len(data.Config.Highlights) > 0 {
			if highlights := strings.Split(data.Config.Highlights, ","); len(highlights) > 0 {
				for _, h := range highlights {
					h = strings.TrimSpace(h)
					if strings.Contains(i.Content, h) {
						i.Content = strings.Replace(i.Content, h, "[:red:]"+h+"[:-:]", -1)
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
			c = gemtextFormat(i.Content, false, TlTui.Emoji)
		} else if TlTui.FilterHighlights == false {
			c = gemtextFormat(i.Content, f, TlTui.Emoji)
			if f == true {
				c = "[:-:b]" + c + "[:-:-]"
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
			content = content + fmt.Sprintf("\n[\"entry-"+strconv.Itoa(nbEntries)+"\"]%v - %v\n%v\n%v[\"\"]\n", d, i.Published.Format(data.Config.Date_format), a, c)
			nbEntries++
		}
	}

	tv := cview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetRegions(true)
	tv.SetToggleHighlights(false)
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

func createTimelineTitle(t time.Time, highlights bool, filter string) string {
	start := ""

	if highlights == true {
		start = start + "[::bu]Highlights[::-]"
	} else {
		start = start + "[::bu]Timeline[::-]"
	}

	if filter != "" {
		start = start + " from [red::]" + filter + "[white::]"
		return fmt.Sprintf("  %v  ", start)
	} else {
		return fmt.Sprintf("  %v - [::i]Refreshed at %v[::-]  ", start, t.Format("15:04 MST"))
	}
}

func createHelpBox() *cview.Panels {
	p := cview.NewPanels()
	p.SetBorder(true)
	p.SetTitle(" [::u]Help[::-]: ")
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
	p.SetTitle(" [::bu]Refreshing stream[::-]: ")
	p.SetPadding(2, 0, 5, 0)

	tv := cview.NewTextView()
	tv.SetTextAlign(cview.AlignCenter)
	tv.SetText("TinyLogs are being refreshed, please wait…")

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

func gemtextFormat(s string, isHighlighted bool, emoji bool) string {
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

	// Format responses:
	// Must be after link format.
	re = regexp.MustCompile(`(?im)^(\[skyblue::b\]→ [^ ]* ){0,1}(re: )(.*)$`)
	startFormat := ""
	if emoji == true {
		startFormat = "💬 "
	} else {
		startFormat = "↳ "
	}
	startFormat = startFormat + "$3 $1"
	s = re.ReplaceAllString(s, startFormat+closeFormat)

	// Format lists:
	re = regexp.MustCompile("(?im)^([*] [^\n]*)")
	if isHighlighted == true {
		s = re.ReplaceAllString(s, "  [::b]$1"+closeFormat)
	} else {
		s = re.ReplaceAllString(s, "  [-:-:-]$1"+closeFormat)
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
	return m
}

func updateFormModalContent(message string, mainButtonName string, buttonName string, execFunc func()) {
	m := TlTui.FormModal
	f := m.GetForm()
	f.ClearButtons()

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

func createResponseStub(dateFormat string) string {
	// Stubs is blocked by cview bug.
	// Temporary empty stub:
	stub := "## " + time.Now().Format(dateFormat) + "\nRE:\n@author:\n\n"
	return stub
}

func copyToClipboard(content string) error {
	return clipboard.WriteAll(content)
}
