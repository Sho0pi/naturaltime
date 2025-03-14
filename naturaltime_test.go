package naturaltime_test

import (
	"testing"
	"time"

	"github.com/sho0pi/naturaltime"
)

func TestNativeParser(t *testing.T) {
	parser := naturaltime.NewNativeParser()

	// Use a fixed reference time for testing
	now := time.Now().Truncate(time.Second)

	tests := []struct {
		expression  string
		expected    time.Time
		shouldError bool
	}{
		{"2023-01-15", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"01/15/2023", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"15/01/2023", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"January 15, 2023", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"15 January 2023", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"today", now, false},
		{"tomorrow", now.AddDate(0, 0, 1), false},
		{"yesterday", now.AddDate(0, 0, -1), false},
		{"next week", now.AddDate(0, 0, 7), false},
		{"invalid date", time.Time{}, true},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			result, err := parser.ParseDate(test.expression)
			if test.shouldError {
				if err == nil {
					t.Errorf("Expected error for expression %q, but got none", test.expression)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseDate(%q) returned error: %v", test.expression, err)
			}

			// Truncate the result and expected time to seconds before comparing
			truncatedResult := result.Truncate(time.Second)
			truncatedExpected := test.expected.Truncate(time.Second)

			if !truncatedResult.Equal(truncatedExpected) {
				t.Errorf("ParseDate(%q) = %v, want %v", test.expression, truncatedResult, truncatedExpected)
			}
		})
	}
}
