package ui

import (
	"fmt"
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cview"

	"git.bacardi55.io/bacardi55/gtl/core"
)

func DisplayStreamTui(data *core.TlData) error {
	e := data.RefreshFeeds()
	if e != nil {
		return fmt.Errorf("Couldn't refresh feeds")
	}

	app := cview.NewApplication()
	app.EnableMouse(true)

	flex := cview.NewFlex()
	flex.SetTitle("Gemini Tiny Logs")
	flex.AddItem(sideBarBox(data.Feeds), 0, 1, false)
	flex.AddItem(contentBox(data), 0, 3, true)

	app.SetRoot(flex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}

	return nil
}

func sideBarBox(tl map[string]core.TlFeed) *cview.Flex {
	mflex := cview.NewFlex()
	mflex.SetDirection(cview.FlexRow)
	mflex.AddItem(createMenu(), 0, 1, false)
	mflex.AddItem(createListTl(tl), 0, 3, false)

	return mflex
}

func contentBox(data *core.TlData) *cview.Panels {
	p := cview.NewPanels()
	p.SetBorder(true)
	p.SetTitle("Timeline:")

	var content string
	t := time.Now()
	for _, i := range data.Stream.Items {
		f := false
		//entryColor := color.New().SprintFunc()
		if len(data.Config.Highlights) > 0 {
			if highlights := strings.Split(data.Config.Highlights, ","); len(highlights) > 0 {
				for _, h := range highlights {
					h = strings.TrimSpace(h)
					if strings.Contains(i.Content, h) {
						i.Content = strings.Replace(i.Content, h, "[:red]"+h+"[:-]", -1)
						//entryColor = color.New(color.Bold).Add(color.Italic).SprintFunc()
						//bo.Println(stream.Items[i].Content, "\n")
						f = true
						break
					}
				}
			}
		}

		c := i.Content
		if f == true {
			c = "[::b]" + c + "[-:-:-]"
		}

		a := fmt.Sprintf("[red]" + i.Author + "[-]")
		d := "[skyblue]" + formatElapsedTime(t.Sub(i.Published)) + "[white]"
		content = content + fmt.Sprintf("%v - %v\n%v\n%v\n\n", d, i.Published.Format(data.Config.Date_format), a, c)
	}

	tv := cview.NewTextView()
	tv.SetText(content)
	tv.SetDynamicColors(true)

	p.AddPanel("timeline", tv, true, true)
	return p
}

func createMenu() *cview.List {
	list := createList("Menu:")
	list.ShowSecondaryText(false)

	i := createListItem("Timeline", "")
	i.SetShortcut(rune('T'))
	list.AddItem(i)

	i = createListItem("Highlights", "")
	i.SetShortcut(rune('H'))
	list.AddItem(i)

	i = createListItem("Refresh", "")
	i.SetShortcut(rune('R'))
	list.AddItem(i)

	i = createListItem("New Entry", "")
	i.SetShortcut(rune('N'))
	list.AddItem(i)

	return list
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
