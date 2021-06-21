package ui

import (
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cbind"
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"

	"git.bacardi55.io/bacardi55/gtl/core"
)

type TlShortcut struct {
	Name        string
	Command     string
	Description string
}

type TlTuiSubs struct {
	Items []core.TlFeed
}

type TlTUI struct {
	App              *cview.Application
	Layout           *cview.Flex
	MainFlex         *cview.Flex
	SideBarBox       *cview.Panels
	ContentBox       *cview.Panels
	HelpBox          *cview.Panels
	RefreshBox       *cview.Panels
	ListTl           *cview.List
	FocusManager     *cview.FocusManager
	Footer           *cview.Panels
	Filter           string
	FilterHighlights bool
	RefreshStream    func(bool)
	LastRefresh      time.Time
	Help             bool
	DisplaySidebar   bool
	Emoji            bool
	Muted            []string
}

func (TlTui *TlTUI) InitApp(useEmoji bool) {
	TlTui.App = cview.NewApplication()
	TlTui.App.EnableMouse(true)

	TlTui.Help = false
	TlTui.LastRefresh = time.Now()
	TlTui.Filter = ""
	TlTui.FilterHighlights = false
	// Todo: make it configurable.
	TlTui.DisplaySidebar = true

	TlTui.Emoji = false
	if useEmoji == true {
		TlTui.Emoji = true
	}
}

func (TlTui *TlTUI) SetAppUI(data *core.TlData) {
	TlTui.Layout = cview.NewFlex()
	TlTui.Layout.SetTitle("Gemini Tiny Logs")
	TlTui.Layout.SetDirection(cview.FlexRow)

	TlTui.MainFlex = cview.NewFlex()
	TlTui.SideBarBox = sideBarBox(data.Feeds)
	TlTui.MainFlex.AddItem(TlTui.SideBarBox, 0, 1, false)
	TlTui.ContentBox = contentBox(data)
	TlTui.MainFlex.AddItem(TlTui.ContentBox, 0, 3, true)
	TlTui.LastRefresh = time.Now()

	TlTui.Footer = createFooter(time.Now(), data.Config.Date_format)
	TlTui.Layout.SetDirection(cview.FlexRow)
	TlTui.Layout.AddItem(TlTui.MainFlex, 0, 1, true)
	//TlTui.Layout.AddItem(TlTui.Footer, 1, 0, false)

	focusManager := cview.NewFocusManager(TlTui.App.SetFocus)
	focusManager.SetWrapAround(true)
	focusManager.Add(TlTui.SideBarBox)
	focusManager.Add(TlTui.ContentBox)
	TlTui.FocusManager = focusManager
	// TODO: Investigate
	// This fix an issue where the first time user hits TAB, it doesn't change focus.
	TlTui.FocusManager.FocusNext()

	TlTui.HelpBox = createHelpBox()
	TlTui.RefreshBox = createRefreshBox()
}

func (TlTui *TlTUI) SetShortcuts() {
	c := cbind.NewConfiguration()

	handleRefresh := func(ev *tcell.EventKey) *tcell.EventKey {
		var refreshStart = func() {
			TlTui.App.SetRoot(TlTui.RefreshBox, true)
		}

		var refreshEnd = func() {
			TlTui.App.SetRoot(TlTui.Layout, true)
		}

		var refresh = func() {
			TlTui.RefreshStream(true)
		}

		TlTui.App.QueueUpdateDraw(refreshStart)
		TlTui.App.QueueUpdate(refresh)
		TlTui.App.QueueUpdateDraw(refreshEnd)
		return nil
	}

	handleHighlights := func(ev *tcell.EventKey) *tcell.EventKey {
		TlTui.FilterHighlights = !TlTui.FilterHighlights
		TlTui.RefreshStream(false)
		return nil
	}

	handleTimeline := func(ev *tcell.EventKey) *tcell.EventKey {
		// Remove Highlights filter:
		TlTui.FilterHighlights = false
		// Remove tinylog filter:
		TlTui.Filter = ""
		// Select "All" in subscription sidebar:
		TlTui.ListTl.SetCurrentItem(0)
		TlTui.RefreshStream(false)
		return nil
	}

	handleTab := func(ev *tcell.EventKey) *tcell.EventKey {
		// If help of if sidebar is hidden, nothing to switch focus to.
		if TlTui.Help == false && TlTui.DisplaySidebar == true {
			TlTui.FocusManager.FocusNext()
		}
		return nil
	}

	handleToggleSidebar := func(ev *tcell.EventKey) *tcell.EventKey {
		if TlTui.DisplaySidebar == true {
			TlTui.MainFlex.RemoveItem(TlTui.SideBarBox)
			TlTui.DisplaySidebar = false
			TlTui.FocusManager.Focus(TlTui.ContentBox)
		} else {
			TlTui.MainFlex.AddItemAtIndex(0, TlTui.SideBarBox, 0, 1, true)
			TlTui.DisplaySidebar = true
			TlTui.FocusManager.Focus(TlTui.SideBarBox)
		}
		return nil
	}

	handleHelp := func(ev *tcell.EventKey) *tcell.EventKey {
		if TlTui.Help == false {
			TlTui.Help = true
			TlTui.App.SetRoot(TlTui.HelpBox, true)
		} else {
			TlTui.Help = false
			TlTui.App.SetRoot(TlTui.Layout, true)
		}
		return nil
	}
	handleQuit := func(ev *tcell.EventKey) *tcell.EventKey {
		// Don't quit if within help.
		if TlTui.Help == true {
			TlTui.Help = false
			TlTui.App.SetRoot(TlTui.Layout, true)
		} else {
			TlTui.App.Stop()
		}
		return nil
	}

	c.SetRune(tcell.ModNone, 'r', handleRefresh)
	c.SetRune(tcell.ModNone, 'h', handleHighlights)
	c.SetRune(tcell.ModNone, 't', handleTimeline)
	c.SetRune(tcell.ModNone, 's', handleToggleSidebar)
	c.SetKey(tcell.ModNone, tcell.KeyTAB, handleTab)
	c.SetRune(tcell.ModNone, '?', handleHelp)
	c.SetRune(tcell.ModNone, 'q', handleQuit)
	TlTui.App.SetInputCapture(c.Capture)
}

// Implement sort.Interface Len.
func (Subs *TlTuiSubs) Len() int {
	return len(Subs.Items)
}

// Implement Interface sort.Interface Less.
func (Subs *TlTuiSubs) Less(i, j int) bool {
	if Subs.Items[i].Status == Subs.Items[j].Status {
		return (strings.Compare(Subs.Items[i].Title, Subs.Items[j].Title) < 0)
	} else {
		return Subs.Items[i].Status < Subs.Items[j].Status
	}
}

// Implement Interface sort.Interface Swap.
func (Subs *TlTuiSubs) Swap(i, j int) {
	Subs.Items[i], Subs.Items[j] = Subs.Items[j], Subs.Items[i]
}
