package ui

import (
	"strings"
	"time"

	"code.rocketnine.space/tslocum/cview"

	"git.bacardi55.io/bacardi55/gtl/core"
)

type TlTuiSubs struct {
	Items []core.TlFeed
}

type TlClipboard struct {
	Enabled    bool
	DateFormat string
}

type TlTUI struct {
	App              *cview.Application
	Layout           *cview.Flex
	MainFlex         *cview.Flex
	SideBarBox       *cview.Panels
	ContentBox       *cview.Panels
	TimelineTV       *cview.TextView
	HelpBox          *cview.Panels
	RefreshBox       *cview.Panels
	ListTl           *cview.List
	FocusManager     *cview.FocusManager
	FormModal        *cview.Modal
	DisplayFormModal bool
	Filter           string
	FilterHighlights bool
	RefreshStream    func(bool)
	LastRefresh      time.Time
	Help             bool
	DisplaySidebar   bool
	Clipboard        TlClipboard
	Muted            []string
	NbEntries        int
	TlConfig         *core.TlConfig
	TlStream         *core.TlStream
}

func (TlTui *TlTUI) InitApp() {
	TlTui.App = cview.NewApplication()
	TlTui.App.EnableMouse(true)

	TlTui.Help = false
	TlTui.LastRefresh = time.Now()
	TlTui.Filter = ""
	TlTui.FilterHighlights = false
	// Todo: make it configurable.
	TlTui.DisplaySidebar = true

	TlTui.DisplayFormModal = false

	TlTui.Clipboard = TlClipboard{
		Enabled: false,
	}
}

func (TlTui *TlTUI) SetAppUI(data *core.TlData) {
	TlTui.TlConfig = data.Config
	TlTui.TlStream = data.Stream

	TlTui.Layout = cview.NewFlex()
	TlTui.Layout.SetTitle("Gemini Tiny Logs")
	TlTui.Layout.SetDirection(cview.FlexRow)

	TlTui.MainFlex = cview.NewFlex()
	TlTui.SideBarBox = sideBarBox(data)
	TlTui.MainFlex.AddItem(TlTui.SideBarBox, 0, 1, false)
	TlTui.ContentBox = contentBox(data)
	TlTui.MainFlex.AddItem(TlTui.ContentBox, 0, 3, true)
	TlTui.LastRefresh = time.Now()

	TlTui.Layout.SetDirection(cview.FlexRow)
	TlTui.Layout.AddItem(TlTui.MainFlex, 0, 1, true)

	focusManager := cview.NewFocusManager(TlTui.App.SetFocus)
	focusManager.SetWrapAround(true)
	focusManager.Add(TlTui.SideBarBox)
	focusManager.Add(TlTui.ContentBox)
	TlTui.FocusManager = focusManager
	TlTui.FocusManager.Focus(TlTui.ContentBox)

	TlTui.HelpBox = createHelpBox(TlTui.TlConfig)
	TlTui.RefreshBox = createRefreshBox(data.Config)

	TlTui.NbEntries = len(data.Stream.Items)

	if data.Config.Tui_copy_stub_clipboard == true {
		TlTui.Clipboard.Enabled = true
		TlTui.Clipboard.DateFormat = data.Config.Date_format
	}

}

func (TlTui *TlTUI) InitTlEditor(tinylogPath string, postScriptPath string, postScriptRefresh bool) error {
	TlTui.FormModal = createFormModal(TlTui.TlConfig)
	TlTui.ContentBox.AddPanel("newEntryModal", TlTui.FormModal, true, false)

	if err := Tle.Init(tinylogPath, postScriptPath, postScriptRefresh); err != nil {
		return err
	}

	return nil
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
