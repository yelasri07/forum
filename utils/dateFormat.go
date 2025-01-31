package utils

import (
	"strconv"
	"time"
)

func DateFromat(date time.Time) string {
	diff := time.Since(date)
	switch {
	case diff.Seconds() < 60:
		return strconv.Itoa(int(diff.Seconds())) + "s"
	case diff.Minutes() < 60:
		return strconv.Itoa(int(diff.Minutes())) + "min"
	case diff.Hours() < 24:
		return strconv.Itoa(int(diff.Hours())) + "h"
	case diff.Hours() < 24*7:
		return strconv.Itoa(int(diff.Hours()/24)) + "j"
	default:
		return strconv.Itoa(int(diff.Hours()/(24*7))) + "week"
	}
}
