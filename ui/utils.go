package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"time"

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
	} else if Ui.Mode == "gemini" {
		return displayStreamGemini(data, limit)
	} else {
		return fmt.Errorf("Unknown mode.")
	}
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

// Open http(s), gemini and gopher links in their dedicated browser.
// Never tested on windows or MacOS yet…
func openLinkInBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return err
	}
	return nil
}
