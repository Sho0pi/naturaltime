package fixed

import (
	"github.com/sho0pi/naturaltime/internal/relative"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// timeRegex matches fixed time expressions like "5pm", "23:10", "6am", "00:00", etc.
// It captures groups for hour, minute (if provided), and AM/PM notation (if applicable).
var fixedTimeRegex = regexp.MustCompile(`\b(?P<hour>\d{1,2})(?::(?P<minute>\d{2}))?\s*(?P<ampm>am|pm)?\b`)

// ParseTimeExpression extracts fixed time expressions and returns a function that modifies a given time to match the extracted time.
// Supported formats:
// - "5" (interpreted as 5:00 AM in 24-hour format)
// - "5am" / "7pm" (12-hour format with AM/PM)
// - "5:10" / "05:30" (24-hour format)
// - "23:10pm" (invalid, but handled by the function)
// - "00:00" (midnight)
func ParseTimeExpression(input string) (relative.TimeShiftFunc, bool) {
	matches := fixedTimeRegex.FindStringSubmatch(input)
	if matches == nil {
		return nil, false
	}

	hour, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, false
	}

	minute := 0
	if matches[2] != "" {
		minute, err = strconv.Atoi(matches[2])
		if err != nil {
			return nil, false
		}
	}

	ampm := strings.ToLower(matches[3])
	if ampm == "pm" && hour < 12 {
		hour += 12
	} else if ampm == "am" && hour == 12 {
		hour = 0
	}

	return func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), hour, minute, 0, 0, t.Location())
	}, true
}
