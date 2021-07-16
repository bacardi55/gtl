package ui

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cbind"
	"github.com/gdamore/tcell/v2"

	"git.bacardi55.io/bacardi55/gtl/core"
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

	c.SetRune(tcell.ModNone, 'O', linksHandler)

	c.SetRune(tcell.ModNone, 'T', threadHandler)

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

		} else if TlTui.ContentBox.HasFocus() == true {
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
			tlfi, e := getSelectedEntryText()
			if e != nil {
				updateFormModalContent("Couldn't find a valid entry.", "Ok", "", func() {})
				toggleFormModal()
			}
			text = createResponseStub(tlfi, TlTui.Clipboard.DateFormat)
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

func linksHandler(ev *tcell.EventKey) *tcell.EventKey {
	if TlTui.SelectedEntry < 0 {
		updateFormModalContent("No selected entry.", "Ok", "", func() {})
		toggleFormModal()
	}

	tlfi, e := getSelectedEntryText()
	if e != nil {
		log.Println(e)
		updateFormModalContent(e.Error(), "Ok", "", func() {})
		toggleFormModal()
	}

	links := extractLinks(tlfi)
	var err error

	if len(links) < 1 {
		updateFormModalContent("No link to open in this entry.", "Ok", "", func() {})
		toggleFormModal()
		err = nil

	} else if len(links) == 1 {
		// Only 1 link found, open it directly.
		err = openLinkInBrowser(strings.Split(links[0], " ")[1])

	} else if len(links) > 1 {
		// More than 1 link found, ask for confirmation.
		m := TlTui.FormModal
		f := m.GetForm()
		f.Clear(true)

		message := "Multiple links detected, open them all?"
		for i, l := range links {
			message += "\n(" + strconv.Itoa(i+1) + ") " + l + "\n"
		}

		m.SetText(message)
		f.AddButton("Yes", func() {
			for _, l := range links {
				openLinkInBrowser(strings.Split(l, " ")[1])
				time.Sleep(100 * time.Millisecond)
			}
			toggleFormModal()
		})

		f.AddButton("No", func() {
			toggleFormModal()
		})
		toggleFormModal()
	}

	if err != nil {
		log.Println(err)
		updateFormModalContent(err.Error(), "Ok", "", func() {})
		toggleFormModal()
	}

	TlTui.FocusManager.Focus(TlTui.ContentBox)
	return nil
}

func threadHandler(ev *tcell.EventKey) *tcell.EventKey {
	TlTui.FocusManager.Focus(TlTui.ContentBox)
	if TlTui.SelectedEntry < 0 {
		updateFormModalContent("No selected entry.", "Ok", "", func() {})
		toggleFormModal()
		return ev
	}

	//TODO: Use getSelectedEntryText() when cview issue is fixed:
	// https://code.rocketnine.space/tslocum/cview/issues/69
	entry := `3 days ago - Tue 13 Jul 2021 18:33 CEST
🤔 @bacardi55
💬 @frrobert 2021-07-13 15:59 UTC
 > lace with the strict option does that
Interesting :). I was thinking about doing something like this but with a more text/gemini friendly output so that gemini browser can actually parse it (almost as a big tinylog)!
I think I'll poc that quickly after I release the next gtl version to see how it can look like :).
`

	if isReponseToEntry(entry) == true {
		index := findOriginalEntry(entry)
		if index == -1 {
			// Nothing found.
			updateFormModalContent("No original entry found.", "Ok", "", nil)
			toggleFormModal()
		} else {
			// Display a thread.
			// TODO: Improve UI if this is fixed:
			// https://code.rocketnine.space/tslocum/cview/issues/71
			// Otherwise, might need to move to another way of displaying the original post.
			updateFormModalContent("Original entry \n\n"+entry, "Ok", "", nil)
			toggleFormModal()
		}
	} else {
		updateFormModalContent("Not a response format, no original to look for.", "Ok", "", func() {})
		toggleFormModal()
	}

	return nil
}

func isReponseToEntry(entry string) bool {
	lines := strings.Split(entry, "\n")
	if len(lines) < 3 {
		return false
	}

	re := regexp.MustCompile(`(?im)^(↳|\x{1F4AC})`)
	return re.MatchString(lines[2])
}

// Find the index in the stream of the original entry.
// Return -1 if original entry isn't found.
func findOriginalEntry(entry string) int {
	lines := strings.Split(entry, "\n")
	r := strings.Split(string(lines[2]), " ")

	if len(r) < 3 {
		return -1
	}

	author := r[1]
	date := core.ParseTlDate(strings.Join(r[2:], " "))

	for i, s := range TlTui.TlStream.Items {
		tmp := strings.Split(s.Author, " ")
		a := ""
		// Removing avatar if any.
		if len(tmp) > 1 {
			a = strings.Join(tmp[1:], " ")
		} else {
			a = tmp[0]
		}

		if strings.Contains(a, author) && s.Published == date {
			return i
		}
	}

	return -1
}

func getSelectedEntryText() (*core.TlFeedItem, error) {
	// TODO: GetRegionText().
	entry := `22 hours ago - Sun 11 Jul 2021 22:04 CEST
🤔 @bacardi55
Just opened my first issue on cview tracker:
→ https://code.rocketnine.space/tslocum/cview/issues/69
→ https://code.rocketnine.space/tslocum/cview/issues/69
When this get resolve, the v0.6.0 could start again with multiple new things coming… :)
With a second link to test:
→ gemini://gmi.bacardi55.io
→ ftp://test.com`

	lines := strings.Split(entry, "\n")

	if len(lines) < 3 {
		return nil, fmt.Errorf("Couldn't parse selected entry - nb of line issue.")
	}

	date := strings.Split(lines[0], "-")
	if len(date) < 2 {
		return nil, fmt.Errorf("Couldn't parse selected entry - date issue.")
	}
	d := core.ParseTlDate(strings.TrimSpace(date[1]))

	tlfi := &core.TlFeedItem{
		Author:    lines[1],
		Content:   strings.Join(lines[2:], "\n"),
		Published: d,
	}

	return tlfi, nil
}

func extractLinks(tlfi *core.TlFeedItem) []string {
	log.Println(tlfi.Content)

	re := regexp.MustCompile("(?im)→ (gemini|gopher|https{0,1})://(.*)$")
	return re.FindAllString(tlfi.Content, -1)
}
