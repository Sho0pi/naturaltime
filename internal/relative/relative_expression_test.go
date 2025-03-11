package relative_test

import (
	"github.com/sho0pi/naturaltime/internal/relative"
	"testing"
	"time"
)

func TestParseRelativeExpression_Valid(t *testing.T) {
	// Use a fixed base time for testing.
	baseTime := time.Date(2025, 3, 11, 15, 0, 0, 0, time.UTC)

	// Table-driven tests for valid relative expressions.
	testCases := []struct {
		name     string
		input    string
		expected time.Time // expected time after applying the shift to baseTime
	}{
		{
			name:     "next week embedded",
			input:    "play football next week with friends",
			expected: baseTime.AddDate(0, 0, 7), // +1 week = 7 days
		},
		{
			name:     "last month embedded",
			input:    "meeting last month at the cafe",
			expected: baseTime.AddDate(0, -1, 0), // -1 month
		},
		{
			name:     "after this day embedded",
			input:    "party after this day in the evening",
			expected: baseTime.AddDate(0, 0, 1), // +1 day
		},
		{
			name:     "past hour embedded",
			input:    "I was there past hour before dinner",
			expected: baseTime.Add(-time.Hour), // -1 hour
		},
		{
			name:     "multiple words with noise",
			input:    "random text, then next week, and some more text",
			expected: baseTime.AddDate(0, 0, 7),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the relative expression from the input.
			shiftFunc, ok := relative.ParseRelativeExpression(tc.input)
			if !ok {
				t.Fatalf("Expected valid parsing, got false for input: %q", tc.input)
			}
			// Apply the returned TimeShiftFunc to the base time.
			result := shiftFunc(baseTime)
			if !result.Equal(tc.expected) {
				t.Errorf("For input %q, expected %v, got %v", tc.input, tc.expected, result)
			}
		})
	}
}

func TestParseRelativeExpression_Invalid(t *testing.T) {
	// Table-driven tests for invalid relative expressions.
	invalidInputs := []string{
		"play football tomorrow", // "tomorrow" is not supported in relative expressions
		"next",                   // missing time unit
		"yesterday",              // not a valid relative expression for this parser
		"next century",           // unit "century" is not defined
		"run fast",               // no valid time phrase present
	}

	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			_, ok := relative.ParseRelativeExpression(input)
			if ok {
				t.Errorf("Expected parsing to fail for input: %q", input)
			}
		})
	}
}
