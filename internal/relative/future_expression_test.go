package relative_test

import (
	"github.com/sho0pi/naturaltime/internal/relative"
	"testing"
	"time"
)

func TestParseFutureExpression_Valid(t *testing.T) {
	baseTime := time.Date(2025, 3, 11, 15, 0, 0, 0, time.UTC)

	testCases := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{"in 5 minutes", "I'll call you in 5 minutes", baseTime.Add(5 * time.Minute)},
		{"in 2 hours", "meeting in 2 hours", baseTime.Add(2 * time.Hour)},
		{"in 1 day", "trip in 1 day", baseTime.AddDate(0, 0, 1)},
		{"in 3 weeks", "vacation in 3 weeks", baseTime.AddDate(0, 0, 21)},
		{"in 7 months", "wedding in 7 months", baseTime.AddDate(0, 7, 0)},
		{"in 10 years", "project in 10 years", baseTime.AddDate(10, 0, 0)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shiftFunc, ok := relative.ParseFutureExpression(tc.input)
			if !ok {
				t.Fatalf("Expected valid parsing, got false for input: %q", tc.input)
			}
			result := shiftFunc(baseTime)
			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestParseFutureExpression_Invalid(t *testing.T) {
	invalidInputs := []string{
		"in five minutes", // does not support spelled-out numbers
		"in 2 centuries",  // unsupported time unit
		"in",              // missing number and unit
		"running fast",    // no time reference
	}

	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			_, ok := relative.ParseFutureExpression(input)
			if ok {
				t.Errorf("Expected parsing to fail for input: %q", input)
			}
		})
	}
}
