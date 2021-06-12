package ui

import (
	"strconv"
	"time"
)

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
