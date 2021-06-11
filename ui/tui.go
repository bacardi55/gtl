package ui

import (
	"fmt"

  "code.rocketnine.space/tslocum/cview"

	"git.bacardi55.io/bacardi55/gtl/core"
)

func DisplayStreamTui(data *core.TlData) error {
  /*
	e := data.RefreshFeeds()
	if e != nil {
		return fmt.Errorf("Couldn't refresh feeds")
	}
  */

	app := cview.NewApplication()
	app.EnableMouse(true)

	flex := cview.NewFlex()
  flex.SetTitle("Gemini Tiny Logs")
	flex.AddItem(sideBarBox(), 0, 1, false)
	flex.AddItem(contentBox(), 0, 3, true)

	app.SetRoot(flex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}

  return nil
}

func sideBarBox() *cview.Flex {
	mflex := cview.NewFlex()
  mflex.SetDirection(cview.FlexRow)
	mflex.AddItem(createMenu(), 0, 1, false)
	mflex.AddItem(createListTl(), 0, 1, false)

  return mflex
}

func contentBox() *cview.Panels {
  p := cview.NewPanels()
  p.SetBorder(true)
  p.SetTitle("Timeline:")
  rows := 25
  var content string
  for r := 0; r < rows; r++ {
    content = content + fmt.Sprintf("%v 🤔 @author@capsule.tld - X hours ago - Fri 11 Jun 2021 20:00 UTC\nA tinylog entry :)\nCan be multiline of course.\n\n", r)
  }

  t := cview.NewTextView()
  t.SetText(content)
  p.AddPanel("timeline", t, true, true)
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

func createListTl() *cview.List {
  list := createList("Subscriptions:")
  list.ShowSecondaryText(true)

  i := createListItem("@author1", "=> gemini://capsule.tld")
  list.AddItem(i)

  i = createListItem("🤔 @author2", "=> gemini://capsule2.tld")
  list.AddItem(i)

  i = createListItem("@author3@capsule3.tld", "=> gemini://capsule3.tld")
  list.AddItem(i)

  i = createListItem("🤔 @author4@capsule4.tld", "=> gemini://capsule4.tld")
  list.AddItem(i)

  i = createListItem("@author5", "=> gemini://capsule5.tld")
  list.AddItem(i)

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
