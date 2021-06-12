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

func DisplayStreamTui(data *core.TlData) error {
	e := data.RefreshFeeds()
	if e != nil {
		return fmt.Errorf("Couldn't refresh feeds")
	}

	app := cview.NewApplication()
	app.EnableMouse(true)

	layoutFlex := cview.NewFlex()
	layoutFlex.SetTitle("Gemini Tiny Logs")
	layoutFlex.SetDirection(cview.FlexRow)

	mainFlex := cview.NewFlex()
	mainFlex.AddItem(sideBarBox(data.Feeds), 0, 1, false)
	cb := contentBox(data)
	mainFlex.AddItem(cb, 0, 3, true)

	layoutFlex.AddItem(createHeader(), 2, 0, false)
	layoutFlex.AddItem(mainFlex, 0, 1, true)

	// Shortcuts:
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlR {
			log.Println("Ctrl-R - Shortcut used, refreshing content.")

			e := data.RefreshFeeds()
			if e != nil {
				log.Println("Couldn't refresh feeds")
				// TODO: Display a message?
			}
			tv := getContentTextView(data, false)
			cb.AddPanel("timeline", tv, true, true)
			return nil
		} else if event.Key() == tcell.KeyCtrlH || event.Key() == tcell.KeyCtrlT {
			log.Println("Ctrl-H - Shortcut used, refreshing content.")

			// No refresh, only filtering.
			tv := getContentTextView(data, (event.Key() == tcell.KeyCtrlH))
			cb.AddPanel("timeline", tv, true, true)
			return nil
		}
		return event
	})

	app.SetRoot(layoutFlex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}

	return nil
}

func sideBarBox(tl map[string]core.TlFeed) *cview.Flex {
	mflex := cview.NewFlex()
	mflex.SetDirection(cview.FlexRow)
	mflex.AddItem(createListTl(tl), 0, 3, false)

	return mflex
}

func contentBox(data *core.TlData) *cview.Panels {
	p := cview.NewPanels()
	p.SetBorder(true)
	p.SetTitle("Timeline:")

	tv := getContentTextView(data, false)

	p.AddPanel("timeline", tv, true, true)
	return p
}

func getContentTextView(data *core.TlData, highlightsOnly bool) *cview.TextView {
	var content string
	t := time.Now()
	for _, i := range data.Stream.Items {
		f := false
		if len(data.Config.Highlights) > 0 {
			if highlights := strings.Split(data.Config.Highlights, ","); len(highlights) > 0 {
				for _, h := range highlights {
					h = strings.TrimSpace(h)
					if strings.Contains(i.Content, h) {
						i.Content = strings.Replace(i.Content, h, "[:red]"+h+"[:-]", -1)
						f = true
						break
					}
				}
			}
		}

		var c string
		ignoreEntry := false
		if highlightsOnly == true && f == true {
			// No bold because all would be bold.
			c = i.Content
		} else if highlightsOnly == false {
			c = i.Content
			if f == true {
				c = "[::b]" + c + "[-:-:-]"
			}
		} else {
			ignoreEntry = true
		}

		if ignoreEntry != true {
			a := fmt.Sprintf("[red]" + i.Author + "[-]")
			d := "[skyblue]" + formatElapsedTime(t.Sub(i.Published)) + "[white]"
			content = content + fmt.Sprintf("%v - %v\n%v\n%v\n\n", d, i.Published.Format(data.Config.Date_format), a, c)
		}
	}

	tv := cview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetText(content)

	return tv
}

func createListTl(tl map[string]core.TlFeed) *cview.List {
	list := createList("Subscriptions:")
	list.ShowSecondaryText(true)

	for _, f := range tl {
		i := createListItem(f.Title, "=> "+f.Link)
		list.AddItem(i)
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

func createList(title string) *cview.List {
	list := cview.NewList()
	list.SetTitle(title)
	list.SetBorder(true)
	list.SetWrapAround(true)

	return list
}

func createHeader() *cview.TextView {
	tv := cview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetMaxLines(2)
	tv.SetTextAlign(cview.AlignCenter)
	content := "Usage:\t[green]Refresh[-]: [::b]Ctrl-R[-:-:-]\t[green]Timeline[-]: [::b]Ctrl-T[-:-:-]\t[green]Highlights[-]: [::b]Ctrl-H[-:-:-]\t[green]Quit[-]: [::b]Ctrl-C[-:-:-]"
	tv.SetText(content)
	return tv
}
