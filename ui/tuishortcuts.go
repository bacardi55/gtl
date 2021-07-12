package ui

import (
	"log"
	"strconv"

	"code.rocketnine.space/tslocum/cbind"
	"github.com/gdamore/tcell/v2"
)

type TlShortcut struct {
	Name        string
	Command     string
	Description string
}

func (TlTui *TlTUI) SetShortcuts() {
	c := cbind.NewConfiguration()

	c.SetRune(tcell.ModNone, 'N', openEditorHandler)
	c.SetRune(tcell.ModNone, 'R', openEditorHandler)

	c.SetRune(tcell.ModNone, 'r', refreshHandler)

	c.SetRune(tcell.ModNone, 'h', mainDisplayHandler)
	c.SetRune(tcell.ModNone, 't', mainDisplayHandler)

	c.SetRune(tcell.ModNone, 'J', tlNavHandler)
	c.SetRune(tcell.ModNone, 'K', tlNavHandler)

	c.SetRune(tcell.ModNone, 's', uiChangeHandler)
	c.SetKey(tcell.ModNone, tcell.KeyTAB, uiChangeHandler)

	c.SetRune(tcell.ModNone, '?', helpToggleDisplay)
	c.SetRune(tcell.ModNone, 'q', uiChangeHandler)

	c.SetKey(tcell.ModNone, tcell.KeyESC, uiChangeHandler)

	TlTui.App.SetInputCapture(c.Capture)
}

// Manage refresh shortcut.
func refreshHandler(ev *tcell.EventKey) *tcell.EventKey {
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

func mainDisplayHandler(ev *tcell.EventKey) *tcell.EventKey {
	if TlTui.DisplayFormModal == true {
		return ev
	}

	if ev.Rune() == 'h' {
		// Toggle Highlights filter:
		TlTui.FilterHighlights = !TlTui.FilterHighlights
		TlTui.RefreshStream(false)
		TlTui.FocusManager.Focus(TlTui.ContentBox)

		return nil
	} else if ev.Rune() == 't' {
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

	return ev
}

// Handle keys that changes UI:
// ESC: hide help/modal/selected entry.
// TAB: switch focus.
// s: hide/Show sidebar.
// q: quit.
func uiChangeHandler(ev *tcell.EventKey) *tcell.EventKey {
	if TlTui.DisplayFormModal == true {
		return ev
	}

	if ev.Key() == tcell.KeyTAB {
		// If help of if sidebar is hidden, nothing to switch focus to.
		if TlTui.Help == false && TlTui.DisplaySidebar == true {
			TlTui.FocusManager.FocusNext()
		}
		return nil

	} else if ev.Key() == tcell.KeyESC {
		if TlTui.Help == true {
			// Hide help if displayed.
			TlTui.Help = false
			TlTui.App.SetRoot(TlTui.Layout, true)
			return nil

		} else if TlTui.DisplayFormModal == true {
			// Hide modal if displayed.
			toggleFormModal()
			return nil

		} else if TlTui.ContentBox.HasFocus() == true && TlTui.SelectedEntry != -1 {
			// Unselect entries if any.
			TlTui.SelectedEntry = -1
			TlTui.TimelineTV.Highlight("")
		}
	} else if ev.Rune() == 's' {
		sidebarToggleDisplay()
		return nil
	} else if ev.Rune() == 'q' {
		// Don't quit if within help.
		if TlTui.Help == true {
			TlTui.Help = false
			TlTui.App.SetRoot(TlTui.Layout, true)
			return nil
		} else {
			TlTui.App.Stop()
		}
	}

	return ev
}

// Hide / Show help panel.
func helpToggleDisplay(ev *tcell.EventKey) *tcell.EventKey {
	if TlTui.Help == false {
		TlTui.Help = true
		TlTui.App.SetRoot(TlTui.HelpBox, true)
	} else {
		TlTui.Help = false
		TlTui.App.SetRoot(TlTui.Layout, true)
	}
	return nil
}

// Hide / Show sidebar.
func sidebarToggleDisplay() {
	if TlTui.DisplaySidebar == true {
		TlTui.MainFlex.RemoveItem(TlTui.SideBarBox)
		TlTui.DisplaySidebar = false
		TlTui.FocusManager.Focus(TlTui.ContentBox)
	} else {
		TlTui.MainFlex.AddItemAtIndex(0, TlTui.SideBarBox, 0, 1, true)
		TlTui.DisplaySidebar = true
		TlTui.FocusManager.Focus(TlTui.SideBarBox)
	}
}

// Manage entry selection.
func tlNavHandler(ev *tcell.EventKey) *tcell.EventKey {
	// Only usable on ContentBox.
	if TlTui.ContentBox.HasFocus() == false {
		return ev
	}

	if ev.Rune() == 'J' {
		// Highlight next item.
		if TlTui.SelectedEntry < TlTui.NbEntries {
			TlTui.SelectedEntry += 1
		}
	} else if ev.Rune() == 'K' {
		if TlTui.SelectedEntry > 0 {
			TlTui.SelectedEntry -= 1
		}
	}
	TlTui.TimelineTV.Highlight("entry-" + strconv.Itoa(TlTui.SelectedEntry))
	TlTui.TimelineTV.ScrollToHighlight()

	return nil
}

// Manage editor related feature.
func openEditorHandler(ev *tcell.EventKey) *tcell.EventKey {
	mainButtonName, buttonName, message, execFunc := "Cancel", "", "", func() {}

	if TlTui.Clipboard.Enabled == true {
		var text string
		if ev.Rune() == 'R' {
			if TlTui.SelectedEntry == -1 {
				updateFormModalContent("Select an entry first to be able respond to it.", "Ok", "", func() {})
				toggleFormModal()
				return ev
			}
			// TODO: Blocked by cview issue.
			text = createResponseStub(TlTui.Clipboard.DateFormat)
		} else if ev.Rune() == 'N' {
			text = createNewEntryStub(TlTui.Clipboard.DateFormat)
		}
		if copyToClipboard(text) != nil {
			log.Println("Couldn't copy Stub to clipboard")
		}
	}

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
						refreshHandler(ev)
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
	}
	TlTui.FocusManager.Focus(TlTui.ContentBox)

	return nil
}
