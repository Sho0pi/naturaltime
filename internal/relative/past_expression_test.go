package relative_test

import (
	"github.com/sho0pi/naturaltime/internal/relative"
	"testing"
	"time"
)

func TestParsePastExpression_Valid(t *testing.T) {
	baseTime := time.Date(2025, 3, 11, 15, 0, 0, 0, time.UTC)

	testCases := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{"5 minutes ago", "I left 5 minutes ago", baseTime.Add(-5 * time.Minute)},
		{"2 hours ago", "meeting was 2 hours ago", baseTime.Add(-2 * time.Hour)},
		{"1 day ago", "trip was 1 day ago", baseTime.AddDate(0, 0, -1)},
		{"3 weeks ago", "vacation was 3 weeks ago", baseTime.AddDate(0, 0, -21)},
		{"7 months ago", "wedding was 7 months ago", baseTime.AddDate(0, -7, 0)},
		{"10 years ago", "project was 10 years ago", baseTime.AddDate(-10, 0, 0)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shiftFunc, ok := relative.ParsePastExpression(tc.input)
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

func TestParsePastExpression_Invalid(t *testing.T) {
	invalidInputs := []string{
		"five minutes ago", // does not support spelled-out numbers
		"2 centuries ago",  // unsupported time unit
		"ago",              // missing number and unit
		"random text",      // no time reference
	}

	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			_, ok := relative.ParsePastExpression(input)
			if ok {
				t.Errorf("Expected parsing to fail for input: %q", input)
			}
		})
	}
}
