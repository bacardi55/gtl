package ui

import (
	"fmt"
	"log"
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"

	"git.bacardi55.io/bacardi55/gtl/core"
)

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
			TlTui.LastRefresh = time.Now()
			TlTui.ListTl = createListTl(data.Feeds)
			TlTui.SideBarBox.AddPanel("subscriptions", TlTui.ListTl, true, true)
		}
		tv := getContentTextView(data)
		TlTui.ContentBox.SetTitle(createTimelineTitle(TlTui.LastRefresh, TlTui.FilterHighlights))
		TlTui.ContentBox.AddPanel("timeline", tv, true, true)

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
			c = i.Content
		} else if TlTui.FilterHighlights == false {
			c = i.Content
			if f == true {
				c = "[::b]" + c + "[::-]"
			}
		} else {
			ignoreEntry = true
		}

		if ignoreEntry != true {
			a := fmt.Sprintf("[red]" + i.Author + "[-::]")
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

	i := createListItem("All Subscriptions", "")
	i.SetSelectedFunc(func() {
		TlTui.Filter = TlTui.ListTl.GetCurrentItem().GetMainText()
		TlTui.RefreshStream(false)
	})
	list.AddItem(i)
	for _, f := range tl {
		it := createListItem(f.DisplayName, "=> "+f.Link)
		list.AddItem(it)
		it.SetSelectedFunc(func() {
			TlTui.Filter = strings.TrimSpace(strings.Split(TlTui.ListTl.GetCurrentItem().GetMainText(), "-")[0])
			TlTui.RefreshStream(false)
		})
	}

	return list
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

func createHeader() *cview.TextView {
	tv := cview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetMaxLines(1)
	tv.SetTextAlign(cview.AlignCenter)
	content := "[::u]Usage[-::-]:\t[green]Refresh[-]: [::b]Ctrl-R[-::-]\t[green]Timeline[-]: [::b]Ctrl-T[-::-]\t[green]Highlights[-:]: [::b]Ctrl-H[-::-]\t[green]Switch focus[-]: [::b]TAB[-::-]\t[green]Quit[-]: [::b]Ctrl-Q/Ctrl-C[-::-]"
	tv.SetText(content)
	return tv
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
