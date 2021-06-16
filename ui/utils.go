package ui

import (
	"fmt"
	"strconv"
	"time"

	"code.rocketnine.space/tslocum/cview"

	"git.bacardi55.io/bacardi55/gtl/core"
)

type TlUI struct {
	Mode string
}

func (Ui *TlUI) Run(data *core.TlData, limit int) error {
	if Ui.Mode == "cli" {
		return displayStreamCli(data, limit)
	} else if Ui.Mode == "tui" {
		return displayStreamTui(data)
	} else {
		return fmt.Errorf("Unknown mode.")
	}
}

type TlTUI struct {
	App              *cview.Application
	Layout           *cview.Flex
	MainFlex         *cview.Flex
	SideBarBox       *cview.Panels
	ContentBox       *cview.Panels
	ListTl           *cview.List
	FocusManager     *cview.FocusManager
	Footer           *cview.Panels
	Filter           string
	FilterHighlights bool
	RefreshStream    func(bool)
	LastRefresh      time.Time
	Help             bool
}

func formatElapsedTime(elapsed time.Duration) string {
	ret := elapsed.Round(time.Second).String()

	if d := int(elapsed.Hours() / 24); d > 0 {
		ret = strconv.Itoa(d) + " day"
		if d > 1 {
			ret = ret + "s"
		}
	} else if h := int(elapsed.Hours()); h > 0 {
		ret = strconv.Itoa(h) + " hour"
		if h > 1 {
			ret = ret + "s"
		}
	} else if m := int(elapsed.Minutes()); m > 0 {
		ret = strconv.Itoa(m) + " minute"
		if m > 1 {
			ret = ret + "s"
		}
	} else {
		ret = strconv.Itoa(int(elapsed.Round(time.Second).Seconds())) + " seconds"
	}

	return ret + " ago"
}
