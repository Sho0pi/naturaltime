package fixed

import (
	"fmt"
	"github.com/sho0pi/naturaltime/internal/relative"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// monthAliases maps common month names and abbreviations (case-insensitive)
// to their corresponding month number (1 = January, 2 = February, etc.).
var monthAliases = map[string]int{
	"january": 1, "jan": 1,
	"february": 2, "feb": 2,
	"march": 3, "mar": 3,
	"april": 4, "apr": 4,
	"may":  5,
	"june": 6, "jun": 6,
	"july": 7, "jul": 7,
	"august": 8, "aug": 8,
	"september": 9, "sep": 9, "sept": 9,
	"october": 10, "oct": 10,
	"november": 11, "nov": 11,
	"december": 12, "dec": 12,
}

// numericDateRegex matches numeric date formats such as:
// "17.2", "7/7", "1/1/2024", or "09.01.98".
// It expects the format: day[separator]month[separator][year].
// The separator can be a dot, slash, or dash.
var numericDateRegex = regexp.MustCompile(`\b(?P<day>\d{1,2})[./-](?P<month>\d{1,2})(?:[./-](?P<year>\d{2,4}))?\b`)

// textualDateRegex is built dynamically in init() to match textual date formats.
// It supports three alternatives:
//  1. "MonthName Day[st|nd|rd|th][, Year]": e.g., "Feb 3rd", "March 28, 2030"
//  2. "Day MonthName[, Year]": e.g., "28 March 2030"
//  3. "MonthName" alone: e.g., "January" (defaults to day 1)
//
// The regex uses a dynamic pattern for month names based on monthAliases.
var textualDateRegex *regexp.Regexp

func init() {
	// Build a regex alternation group for valid month names from monthAliases keys.
	var monthNames []string
	for k := range monthAliases {
		monthNames = append(monthNames, k)
	}
	monthPattern := "(?i:" + strings.Join(monthNames, "|") + ")"
	// Build the complete regex pattern with three alternatives:
	// Option 1: MonthName Day (with optional ordinal suffix) (optional year)
	// Option 2: Day MonthName (optional year)
	// Option 3: Only MonthName
	pattern := fmt.Sprintf(`\b(?:(?P<monthName>%s)\s+(?P<day>\d{1,2})(?:st|nd|rd|th)?(?:,?\s+(?P<year>\d{2,4}))?|\b(?P<day2>\d{1,2})\s+(?P<monthName2>%s)(?:,?\s+(?P<year2>\d{2,4}))?|\b(?P<onlyMonth>%s))\b`, monthPattern, monthPattern, monthPattern)
	textualDateRegex = regexp.MustCompile(pattern)
}

// ParseFixedDateExpression parses fixed date expressions from an input string.
// Supported formats include numeric dates (e.g., "17.2", "7/7", "1/1/2024", "09.01.98")
// and textual dates (e.g., "Feb 3rd", "28 March 2030", "January").
// It returns a TimeShiftFunc—a function that takes a reference time (of type time.Time)
// and returns a new time.Time with the parsed date (year, month, day) and the time set to midnight.
// If the year is omitted in the input, the function uses the reference year.
// For two-digit years, a heuristic is applied: if the year is less than 50, it is assumed to be 2000+year;
// otherwise, 1900+year.
func ParseFixedDateExpression(input string) (relative.TimeShiftFunc, bool) {
	// First, try to match the numeric date format.
	if matches := numericDateRegex.FindStringSubmatch(input); matches != nil {
		day, _ := strconv.Atoi(matches[1])
		month, _ := strconv.Atoi(matches[2])
		year := 0
		if matches[3] != "" {
			year, _ = strconv.Atoi(matches[3])
			// Two-digit year heuristic
			if year < 100 {
				if year < 50 {
					year += 2000
				} else {
					year += 1900
				}
			}
		}
		return func(ref time.Time) time.Time {
			if year == 0 {
				year = ref.Year()
			}
			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, ref.Location())
		}, true
	}
	// If numeric format did not match, try the textual date format.
	if matches := textualDateRegex.FindStringSubmatch(input); matches != nil {
		var day, year, month int
		// Option 1: Format "MonthName Day [Year]" (e.g., "Feb 3rd", "March 28, 2030")
		if matches[1] != "" {
			month = monthAliases[strings.ToLower(matches[1])]
			day, _ = strconv.Atoi(matches[2])
			if matches[3] != "" {
				year, _ = strconv.Atoi(matches[3])
				if year < 100 {
					if year < 50 {
						year += 2000
					} else {
						year += 1900
					}
				}
			}
		} else if matches[4] != "" { // Option 2: Format "Day MonthName [Year]" (e.g., "28 March 2030")
			day, _ = strconv.Atoi(matches[4])
			month = monthAliases[strings.ToLower(matches[5])]
			if matches[6] != "" {
				year, _ = strconv.Atoi(matches[6])
				if year < 100 {
					if year < 50 {
						year += 2000
					} else {
						year += 1900
					}
				}
			}
		} else if matches[7] != "" { // Option 3: Only MonthName (e.g., "January" defaults to the 1st)
			month = monthAliases[strings.ToLower(matches[7])]
			day = 1
		}
		return func(ref time.Time) time.Time {
			if year == 0 {
				year = ref.Year()
			}
			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, ref.Location())
		}, true
	}
	return nil, false
}
