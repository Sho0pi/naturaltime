package relative

import (
	"regexp"
	"strconv"
	"time"
)

// futureTimeRegex is a regular expression used to detect time expressions of the format
// "in [amount] [unit]" within a given string.
// It captures two named groups:
// - "amount": The numerical value representing the amount of time.
// - "unit": The time unit (e.g., seconds, minutes, hours, days, etc.).
// This regex ensures that the phrase starts with "in" and is followed by a valid time unit.
var futureTimeRegex = regexp.MustCompile(`\bin\s+(?P<amount>\d+)\s+(?P<unit>s|sec|second|seconds|m|min|mins|minute|minutes|h|hr|hrs|hour|hours|d|day|days|w|week|weeks|mo|mon|mos|month|months|y|yr|year|years)\b`)

// ParseFutureExpression parses time expressions of the format "in [amount] [unit]" and returns
// a TimeShiftFunc, which can be applied to a time.Time value to get the corresponding future time.
//
// This function:
// 1. Uses futureTimeRegex to extract the numerical amount and time unit from the input string.
// 2. Converts the extracted amount to an integer.
// 3. Maps the unit alias to its standardized form.
// 4. Retrieves the corresponding time modification function from timeShiftFunctions.
// 5. Returns a function that applies the time shift by the given amount.
//
// If the input does not match the expected format or the conversion fails,
// the function returns (nil, false).
func ParseFutureExpression(input string) (TimeShiftFunc, bool) {
	matches := futureTimeRegex.FindStringSubmatch(input)
	if matches == nil {
		return nil, false
	}

	// Extract numerical amount and convert to int
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, false
	}

	// Lookup standard unit from alias
	standardUnit := aliasToStandardUnit[matches[2]]
	shiftFunc, exists := timeShiftFunctions[standardUnit]
	if !exists {
		return nil, false
	}

	// Return a function that applies the time shift
	return func(t time.Time) time.Time {
		return shiftFunc(t, amount)
	}, true
}
