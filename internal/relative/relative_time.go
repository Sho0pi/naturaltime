package relative

import (
	"regexp"
	"time"
)

// relativeTimeModifiers maps modifier keywords to their integer representation
var relativeTimeModifiers = map[string]int{
	"last":       -1,
	"past":       -1,
	"next":       1,
	"after this": 1,
}

// relativeTimeRegex is a regular expression designed to identify time expressions such as
// "next week", "last month", "past year", and "after this day" within a text.
// It captures two named groups:
// - "modifier": A keyword that defines the direction of time (e.g., "next", "last", "past", "after this").
// - "unit": A valid time unit (e.g., seconds, minutes, hours, days, weeks, months, years).
// This regex ensures flexibility in recognizing various relative time phrases.
var relativeTimeRegex = regexp.MustCompile(`\b(?P<modifier>next|last|past|after this)\s+(?P<unit>s|sec|second|seconds|m|min|mins|minute|minutes|h|hr|hrs|hour|hours|d|day|days|w|week|weeks|mo|mon|mos|month|months|y|yr|year|years)\b`)

// ParseRelativeExpression extracts and interprets relative time expressions such as
// "next month", "last year", "past week", and "after this day".
// It returns a TimeShiftFunc that modifies a given time accordingly.
//
// The function operates as follows:
// 1. It scans the input string for a match using relativeTimeRegex.
// 2. It extracts the modifier (e.g., "next", "last") and maps it to a numerical shift value.
// 3. It standardizes the time unit by looking it up in aliasToStandardUnit.
// 4. It retrieves the appropriate time shift function from timeShiftFunctions.
// 5. It returns a function that applies the time shift in the direction specified by the modifier.
//
// If parsing fails due to a missing or invalid modifier/unit, the function returns (nil, false).
func ParseRelativeExpression(input string) (TimeShiftFunc, bool) {
	matches := relativeTimeRegex.FindStringSubmatch(input)
	if matches == nil {
		return nil, false
	}

	// Extract modifier and unit
	modifier := relativeTimeModifiers[matches[1]]
	standardUnit := aliasToStandardUnit[matches[2]]
	timeShiftFunc, exists := timeShiftFunctions[standardUnit]
	if !exists {
		return nil, false
	}

	// Return a function that applies the time shift
	return func(t time.Time) time.Time {
		return timeShiftFunc(t, modifier)
	}, true
}
