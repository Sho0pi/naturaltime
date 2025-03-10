package naturaltime_test

import (
	"testing"
	"time"

	"github.com/sho0pi/naturaltime"
)

func TestParseDate(t *testing.T) {
	parser, err := naturaltime.New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	// Use a fixed date but with local timezone
	localLoc := time.Local
	now := time.Date(2023, 1, 15, 12, 0, 0, 0, localLoc)

	tests := []struct {
		expression  string
		expected    time.Time
		shouldBeNil bool
	}{
		{"today", time.Date(2023, 1, 15, 12, 0, 0, 0, localLoc), false},
		{"tomorrow", time.Date(2023, 1, 16, 12, 0, 0, 0, localLoc), false},
		{"yesterday", time.Date(2023, 1, 14, 12, 0, 0, 0, localLoc), false},
		{"next Monday", time.Date(2023, 1, 16, 12, 0, 0, 0, localLoc), false}, // Jan 15, 2023 is a Sunday, so next Monday is Jan 16
		{"3pm", time.Date(2023, 1, 15, 15, 0, 0, 0, localLoc), false},
		{"January 20", time.Date(2023, 1, 20, 12, 0, 0, 0, localLoc), false},
		{"invalid date expression", time.Time{}, true},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			result, err := parser.ParseDate(test.expression, now)

			if err != nil {
				t.Fatalf("ParseDate(%q, %v) returned error: %v", test.expression, now, err)
			}

			if test.shouldBeNil {
				if result != nil {
					t.Errorf("ParseDate(%q, %v) = %v, want nil", test.expression, now, result)
				}
				return
			}

			if result == nil {
				t.Fatalf("ParseDate(%q, %v) returned nil", test.expression, now)
			}

			// Only compare the parts we care about (truncate to minute precision for simplicity)
			if !result.Truncate(time.Minute).Equal(test.expected.Truncate(time.Minute)) {
				t.Errorf("ParseDate(%q, %v) = %v, want %v", test.expression, now, result, test.expected)
			}
		})
	}
}

func TestParseRange(t *testing.T) {
	parser, err := naturaltime.New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	// Use a fixed date but with local timezone
	localLoc := time.Local
	now := time.Date(2023, 1, 15, 12, 0, 0, 0, localLoc)

	tests := []struct {
		expression string
		startTime  time.Time
		endTime    time.Time
		duration   string // Expected duration as a string (e.g., "2h")
	}{
		{
			"today from 2pm to 4pm",
			time.Date(2023, 1, 15, 14, 0, 0, 0, localLoc),
			time.Date(2023, 1, 15, 16, 0, 0, 0, localLoc),
			"2h0m0s",
		},
		{
			"tomorrow 9am-5pm",
			time.Date(2023, 1, 16, 9, 0, 0, 0, localLoc),
			time.Date(2023, 1, 16, 17, 0, 0, 0, localLoc),
			"8h0m0s",
		},
		{
			"next Monday 10:00-11:30",
			time.Date(2023, 1, 16, 10, 0, 0, 0, localLoc),
			time.Date(2023, 1, 16, 11, 30, 0, 0, localLoc),
			"1h30m0s",
		},
		{
			"from 3pm to 5pm",
			time.Date(2023, 1, 15, 15, 0, 0, 0, localLoc),
			time.Date(2023, 1, 15, 17, 0, 0, 0, localLoc),
			"2h0m0s",
		},
		// implement this in the future
		//{
		//	"for 2 hours",
		//	time.Date(2023, 1, 15, 12, 0, 0, 0, localLoc),
		//	time.Date(2023, 1, 15, 14, 0, 0, 0, localLoc),
		//	"2h0m0s",
		//},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			result, err := parser.ParseRange(test.expression, now)
			if err != nil {
				t.Fatalf("ParseRange(%q, %v) returned error: %v", test.expression, now, err)
			}

			if result == nil {
				t.Fatalf("ParseRange(%q, %v) returned nil", test.expression, now)
			}

			expectedDuration, _ := time.ParseDuration(test.duration)

			// Only compare the parts we care about (truncate to minute precision for simplicity)
			if !result.Start().Truncate(time.Minute).Equal(test.startTime.Truncate(time.Minute)) {
				t.Errorf("ParseRange(%q, %v).Start() = %v, want %v", test.expression, now, result.Start(), test.startTime)
			}

			if !result.End().Truncate(time.Minute).Equal(test.endTime.Truncate(time.Minute)) {
				t.Errorf("ParseRange(%q, %v).End() = %v, want %v", test.expression, now, result.End(), test.endTime)
			}

			// For duration, we can compare directly
			if result.Duration != expectedDuration {
				t.Errorf("ParseRange(%q, %v).Duration = %v, want %v", test.expression, now, result.Duration, expectedDuration)
			}
		})
	}
}
