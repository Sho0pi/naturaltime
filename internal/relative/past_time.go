package relative

import (
	"regexp"
	"strconv"
	"time"
)

// pastTimeRegex is a regular expression used to identify time expressions
// in the format of "[number] [unit] ago". It captures two named groups:
// - "number": The numerical value representing the amount of time.
// - "unit": The unit of time (e.g., seconds, minutes, hours, days, etc.).
// The regex ensures that the phrase "ago" appears at the end
var pastTimeRegex = regexp.MustCompile(`\b(?P<number>\d+)\s+(?P<unit>s|sec|second|seconds|m|min|mins|minute|minutes|h|hr|hrs|hour|hours|d|day|days|w|week|weeks|mo|mon|mos|month|months|y|yr|year|years)\s+ago\b`)

// ParsePastExpression extracts and converts time expressions of the format
// "[number] [unit] ago" into a TimeShiftFunc, which can be applied to a
// time.Time value to get the corresponding pastime.
//
// It first matches the input string using pastTimeRegex. If a match is found:
// 1. The number of time units is extracted and converted to an integer.
// 2. The corresponding standard time unit is looked up from aliasToStandardUnit.
// 3. The appropriate time-shifting function is retrieved from timeShiftFunctions.
// 4. A closure function is returned, applying the negative modifier (-number) to the time.
//
// If parsing fails at any step, the function returns (nil, false).
func ParsePastExpression(input string) (TimeShiftFunc, bool) {
	matches := pastTimeRegex.FindStringSubmatch(input)
	if matches == nil {
		return nil, false
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, false
	}

	unit, exists := aliasToStandardUnit[matches[2]]
	if !exists {
		return nil, false
	}

	unitFunc, exists := timeShiftFunctions[unit]
	if !exists {
		return nil, false
	}

	return func(t time.Time) time.Time {
		return unitFunc(t, -num) // Negative modifier for "ago"
	}, true
}
