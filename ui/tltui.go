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
	FormModal        *cview.Modal
	DisplayFormModal bool
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

	TlTui.DisplayFormModal = false

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
	TlTui.FocusManager.Focus(TlTui.ContentBox)

	TlTui.HelpBox = createHelpBox()
	TlTui.RefreshBox = createRefreshBox()
}

func (TlTui *TlTUI) InitTlEditor(tinylogPath string, postScriptPath string, postScriptRefresh bool) error {
	TlTui.FormModal = createFormModal()
	TlTui.ContentBox.AddPanel("newEntryModal", TlTui.FormModal, true, false)

	if err := Tle.Init(tinylogPath, postScriptPath, postScriptRefresh); err != nil {
		return err
	}

	return nil
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
		TlTui.FocusManager.Focus(TlTui.ContentBox)
		return nil
	}

	handleHighlights := func(ev *tcell.EventKey) *tcell.EventKey {
		if TlTui.DisplayFormModal == true {
			return ev
		}
		TlTui.FilterHighlights = !TlTui.FilterHighlights
		TlTui.RefreshStream(false)
		TlTui.FocusManager.Focus(TlTui.ContentBox)
		return nil
	}

	handleTimeline := func(ev *tcell.EventKey) *tcell.EventKey {
		if TlTui.DisplayFormModal == true {
			return ev
		}
		// Remove Highlights filter:
		TlTui.FilterHighlights = false
		// Remove tinylog filter:
		TlTui.Filter = ""
		// Select "All" in subscription sidebar:
		TlTui.ListTl.SetCurrentItem(0)
		TlTui.RefreshStream(false)
		TlTui.FocusManager.Focus(TlTui.ContentBox)
		return nil
	}

	handleTab := func(ev *tcell.EventKey) *tcell.EventKey {
		if TlTui.DisplayFormModal == true {
			return ev
		}
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

	handleNewEntry := func(ev *tcell.EventKey) *tcell.EventKey {
		mainButtonName, buttonName, message, execFunc := "Cancel", "", "", func() {}

		if TlTui.App.Suspend(editTl) == true {
			message = "Tinylog edited successfully"

			if Tle.PostEditionScript != "" {
				message = message + "\nDo you want to run the post edition script?"
				buttonName = "Run script"
				execFunc = func() {
					var m string
					if e := Tle.Push(); e != nil {
						m = "Couldn't run script, please check the logs."
						buttonName = "ok"
						updateFormModalContent(m, "ok", "", func() {})
						TlTui.FocusManager.Focus(TlTui.ContentBox)
					} else {
						toggleFormModal()
						if Tle.PostScriptRefresh == true {
							handleRefresh(ev)
						}
					}
				}
			} else {
				message = ""
				mainButtonName = ""
			}
		} else {
			buttonName = ""
			execFunc = nil
			message = "Tinylog couldn't be edited"
		}

		if message != "" {
			updateFormModalContent(message, mainButtonName, buttonName, execFunc)
			toggleFormModal()
			TlTui.FocusManager.Focus(TlTui.ContentBox)
		} else {
			TlTui.FocusManager.Focus(TlTui.ContentBox)
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

	handleEsc := func(ev *tcell.EventKey) *tcell.EventKey {
		if TlTui.Help == true {
			TlTui.Help = false
			TlTui.App.SetRoot(TlTui.Layout, true)
			return nil
		} else if TlTui.DisplayFormModal == true {
			toggleFormModal()
			return nil
		}
		return ev
	}

	c.SetRune(tcell.ModCtrl, 'n', handleNewEntry)
	c.SetRune(tcell.ModNone, 'r', handleRefresh)
	c.SetRune(tcell.ModNone, 'h', handleHighlights)
	c.SetRune(tcell.ModNone, 't', handleTimeline)
	c.SetRune(tcell.ModNone, 's', handleToggleSidebar)
	c.SetKey(tcell.ModNone, tcell.KeyTAB, handleTab)
	c.SetRune(tcell.ModNone, '?', handleHelp)
	c.SetRune(tcell.ModNone, 'q', handleQuit)
	c.SetKey(tcell.ModNone, tcell.KeyESC, handleEsc)
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
