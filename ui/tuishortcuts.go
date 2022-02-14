package ui

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
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

func (TlTui *TlTUI) SetShortcuts() {
	c := cbind.NewConfiguration()

	c.SetRune(tcell.ModNone, 'N', openEditorHandler)
	c.SetRune(tcell.ModNone, 'R', openEditorHandler)

	c.SetRune(tcell.ModNone, 'O', entryHandler)
	c.SetRune(tcell.ModNone, 'T', entryHandler)
	c.SetKey(tcell.ModAlt, tcell.KeyEnter, entryHandler)

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
		TlTui.FocusManager.Focus(TlTui.RefreshBox)
	}

	var refreshEnd = func() {
		TlTui.App.SetRoot(TlTui.Layout, true)
		TlTui.FocusManager.Focus(TlTui.ContentBox)
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
	if ev.Key() == tcell.KeyTAB {
		if TlTui.DisplayFormModal == true {
			return ev
		}
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
			TlTui.TimelineTV.Highlight("")
		}
	} else if ev.Rune() == 's' {
		if TlTui.DisplayFormModal == true {
			return ev
		}
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
		TlTui.FocusManager.Focus(TlTui.ContentBox)
	}

	selectedEntry := getSelectedEntryNumber()
	if ev.Rune() == 'J' {
		max := TlTui.NbEntries
		if TlTui.TlConfig.Tui_max_entries > 0 && TlTui.TlConfig.Tui_max_entries < TlTui.NbEntries {
			max = TlTui.TlConfig.Tui_max_entries
		}

		// Highlight next item.
		if selectedEntry < max-1 {
			selectedEntry += 1
		} else {
			selectedEntry = max - 1
		}
	} else if ev.Rune() == 'K' {
		if selectedEntry > 0 {
			selectedEntry -= 1
		} else {
			selectedEntry = 0
		}
	}
	TlTui.TimelineTV.Highlight("entry-" + strconv.Itoa(selectedEntry))
	TlTui.TimelineTV.ScrollToHighlight()

	return nil
}

// Manage shortcuts related to entries.
func entryHandler(ev *tcell.EventKey) *tcell.EventKey {
	if ev.Rune() == 'O' {
		return linksHandler(ev)
	} else if ev.Rune() == 'T' {
		return threadHandler(ev)
	} else if ev.Modifiers() == tcell.ModAlt && ev.Key() == tcell.KeyEnter {
		return shiftEnterHandler(ev)
	}
	return ev
}

func shiftEnterHandler(ev *tcell.EventKey) *tcell.EventKey {
	if !TlTui.ContentBox.HasFocus() || TlTui.DisplayFormModal {
		return ev
	}

	if getSelectedEntryNumber() < 0 {
		updateFormModalContent("Select an entry first.", "Ok", "", func() {})
		toggleFormModal()
		return nil
	}

	tlfi, e := getSelectedEntryText()
	if e != nil {
		log.Println(e)
		updateFormModalContent(e.Error(), "Ok", "", func() {})
		toggleFormModal()
	}

	TlTui.FormModal.GetForm().Clear(true)
	TlTui.FormModal.SetTextAlign(cview.AlignLeft)

	fe := gemtextFormatModal(tlfi)
	TlTui.FormModal.SetText(fe)

	TlTui.FormModal.GetForm().AddButton("Open Links", func() {
		toggleFormModal()
		linksHandler(nil)
	})
	TlTui.FormModal.GetForm().AddButton("Open Thread", func() {
		toggleFormModal()
		threadHandler(nil)
	})
	TlTui.FormModal.GetForm().AddButton("Reply", func() {
		toggleFormModal()
		e := tcell.NewEventKey(0, 'R', tcell.ModNone)
		openEditorHandler(e)
	})
	TlTui.FormModal.GetForm().AddButton("Cancel", func() {
		toggleFormModal()
	})

	toggleFormModal()

	return ev
}

// Manage editor related feature.
func openEditorHandler(ev *tcell.EventKey) *tcell.EventKey {
	if TlTui.Clipboard.Enabled || TlTui.TlConfig.Tui_show_stub {
		var text string
		if ev.Rune() == 'R' {
			if getSelectedEntryNumber() < 0 {
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
		if TlTui.Clipboard.Enabled {
			if copyToClipboard(text) != nil {
				log.Println("Couldn't copy Stub to clipboard")
			}
		}

		if TlTui.TlConfig.Tui_show_stub {
			TlTui.FormModal.GetForm().Clear(true)
			// The format is ugly because an issue in cview:
			// https://code.rocketnine.space/tslocum/cview/issues/72#issuecomment-3968
			TlTui.FormModal.SetText(strings.Replace(text, "\n", "\n\n", -1))
			TlTui.FormModal.SetTextAlign(cview.AlignLeft)
			TlTui.FormModal.GetForm().AddButton("Reply", func() {
				toggleFormModal()
				launchEditor()
			})
			toggleFormModal()
			return nil
		}
	}

	launchEditor()
	return nil
}

func launchEditor() {
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
						refreshHandler(nil)
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
}

func linksHandler(ev *tcell.EventKey) *tcell.EventKey {
	if getSelectedEntryNumber() < 0 {
		updateFormModalContent("No selected entry.", "Ok", "", func() {})
		toggleFormModal()
		return nil
	}

	tlfi, e := getSelectedEntryText()
	if e != nil {
		log.Println(e)
		updateFormModalContent(e.Error(), "Ok", "", func() {})
		toggleFormModal()
		return nil
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

		message := "Multiple links detected, open them all?\n"
		for i, l := range links {
			message += "(" + strconv.Itoa(i+1) + ") " + strings.TrimSpace(l) + "\n\n"
			f.AddButton(strconv.Itoa(i+1), func() {
				openLinkInBrowser(strings.Split(l, " ")[1])
				time.Sleep(100 * time.Millisecond)
			})
		}

		m.SetText(message)
		f.AddButton("All", func() {
			for _, l := range links {
				openLinkInBrowser(strings.Split(l, " ")[1])
				time.Sleep(100 * time.Millisecond)
			}
			toggleFormModal()
		})

		f.AddButton("Cancel", func() {
			toggleFormModal()
		})
		TlTui.FormModal.SetTextAlign(cview.AlignLeft)
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
	if getSelectedEntryNumber() < 0 {
		updateFormModalContent("No selected entry.", "Ok", "", func() {})
		toggleFormModal()
		return ev
	}

	tlfi, e := getSelectedEntryText()
	if e != nil {
		updateFormModalContent("Couldn't find a valid entry.", "Ok", "", func() {})
		toggleFormModal()
	}
	entry := tlfi.Content

	if isReponseToEntry(entry) == true {
		tlfi := findOriginalEntry(entry)
		if tlfi == nil {
			// Nothing found.
			updateFormModalContent("No original entry found.", "Ok", "", nil)
			toggleFormModal()
		} else {
			fe := gemtextFormatModal(tlfi)
			fe = "Original entry:\n\n" + fe
			updateFormModalContent(fe, "Ok", "", nil)
			TlTui.FormModal.SetTextAlign(cview.AlignLeft)
			toggleFormModal()
		}
	} else {
		updateFormModalContent("Not a response format, no original to look for.", "Ok", "", func() {})
		toggleFormModal()
	}

	return nil
}

func gemtextFormatModal(tlfi *core.TlFeedItem) string {
	// TODO: Bug in cview when modal text contains [-:-:-] or other.
	// https://code.rocketnine.space/tslocum/cview/issues/72
	t := time.Now()
	d := formatElapsedTime(t.Sub(tlfi.Published))
	a := tlfi.Author
	//c := strings.Replace(gemtextFormat(tlfi.Content, false, TlTui.TlConfig.Tui_status_emoji), "\n", "\n\n", -1)
	c := tlfi.Content
	fe := d + " - " + tlfi.Published.Format(TlTui.TlConfig.Date_format) + "\n\n" + a + "\n\n" + c + "\n"

	return fe
}

func isReponseToEntry(entry string) bool {
	lines := strings.Split(entry, "\n")

	re := regexp.MustCompile(`(?im)^(↳|\x{1F4AC})`)
	return re.MatchString(lines[0])
}

// Find the index in the stream of the original entry.
// Return -1 if original entry isn't found.
func findOriginalEntry(entry string) *core.TlFeedItem {
	lines := strings.Split(entry, "\n")

	line := lines[0]
	index := 0
	author := ""
	if i := strings.Index(line, "@"); i > -1 {
		author = strings.Split(line[i:], " ")[0]
		index = i + len(author)
	}

	index2 := strings.Index(line, "→")
	if index2 == -1 {
		index2 = len(line) - 1
	}
	stringDate := strings.TrimSpace(line[index:index2])

	date := core.ParseTlDate(stringDate)

	for _, s := range TlTui.TlStream.Items {
		a := ""

		tmp := strings.Split(s.Author, " ")
		// Removing avatar if any.
		if len(tmp) > 1 {
			a = strings.Join(tmp[1:], " ")
		} else {
			a = tmp[0]
		}

		if (strings.Contains(a, author) || strings.Contains(author, a)) && s.Published.Truncate(time.Minute) == date.Truncate(time.Minute) {
			return s
		}
	}

	return nil
}

func getSelectedEntryText() (*core.TlFeedItem, error) {
	entry := TlTui.TimelineTV.GetRegionText("entry-" + strconv.Itoa(getSelectedEntryNumber()))
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
		Content:   strings.TrimSpace(strings.Join(lines[2:], "\n")),
		Published: d,
	}

	return tlfi, nil
}

func extractLinks(tlfi *core.TlFeedItem) []string {
	re := regexp.MustCompile("(?im)→ (gemini|gopher|https{0,1})://(.*)$")
	return re.FindAllString(tlfi.Content, -1)
}

func getSelectedEntryNumber() int {
	h := TlTui.TimelineTV.GetHighlights()

	if len(h) == 1 {
		i, e := strconv.Atoi(strings.Replace(h[0], "entry-", "", -1))
		if e == nil {
			return i
		}
	}

	return -1
}
