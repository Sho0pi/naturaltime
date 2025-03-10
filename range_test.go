package naturaltime_test

import (
	"testing"
	"time"

	"github.com/sho0pi/naturaltime"
)

func TestRangeContains(t *testing.T) {
	now := time.Now()
	r := naturaltime.NewRange(now, 2*time.Hour)

	if !r.Contains(now) {
		t.Errorf("Range should contain its start time")
	}

	if !r.Contains(now.Add(1 * time.Hour)) {
		t.Errorf("Range should contain times within its duration")
	}

	if r.Contains(now.Add(2 * time.Hour)) {
		t.Errorf("Range should not contain its end time (exclusive)")
	}

	if r.Contains(now.Add(-1 * time.Minute)) {
		t.Errorf("Range should not contain times before its start")
	}

	if r.Contains(now.Add(3 * time.Hour)) {
		t.Errorf("Range should not contain times after its end")
	}
}

func TestRangeContainsRange(t *testing.T) {
	now := time.Now()
	r1 := naturaltime.NewRange(now, 3*time.Hour)
	r2 := naturaltime.NewRange(now.Add(1*time.Hour), 1*time.Hour)
	r3 := naturaltime.NewRange(now.Add(-1*time.Hour), 1*time.Hour)
	r4 := naturaltime.NewRange(now, 4*time.Hour)

	if !r1.ContainsRange(r2) {
		t.Errorf("Range should contain a range that starts after it and ends before it")
	}

	if r1.ContainsRange(r3) {
		t.Errorf("Range should not contain a range that starts before it")
	}

	if r1.ContainsRange(r4) {
		t.Errorf("Range should not contain a range that ends after it")
	}
}

func TestRangeOverlaps(t *testing.T) {
	now := time.Now()
	r1 := naturaltime.NewRange(now, 2*time.Hour)
	r2 := naturaltime.NewRange(now.Add(1*time.Hour), 2*time.Hour)
	r3 := naturaltime.NewRange(now.Add(3*time.Hour), 1*time.Hour)
	r4 := naturaltime.NewRange(now.Add(-1*time.Hour), 30*time.Minute)
	r5 := naturaltime.NewRange(now.Add(2*time.Hour), 1*time.Hour) // Starts exactly at r1's end

	if !r1.Overlaps(r2) {
		t.Errorf("Ranges should overlap when one starts during the other")
	}

	if r1.Overlaps(r3) {
		t.Errorf("Ranges should not overlap when one starts after the other ends")
	}

	if r1.Overlaps(r4) {
		t.Errorf("Ranges should not overlap when one ends before the other starts")
	}

	if r1.Overlaps(r5) {
		t.Errorf("Ranges should not overlap when one starts exactly at the other's end")
	}
}

func TestRangeIntersection(t *testing.T) {
	now := time.Now()

	t.Run("Overlapping ranges", func(t *testing.T) {
		r1 := naturaltime.NewRange(now, 2*time.Hour)
		r2 := naturaltime.NewRange(now.Add(1*time.Hour), 2*time.Hour)

		intersection := r1.Intersection(r2)
		expected := naturaltime.NewRange(now.Add(1*time.Hour), 1*time.Hour)

		if !intersection.Equal(expected) {
			t.Errorf("Intersection = %v, want %v", intersection, expected)
		}
	})

	t.Run("Non-overlapping ranges", func(t *testing.T) {
		r1 := naturaltime.NewRange(now, 1*time.Hour)
		r2 := naturaltime.NewRange(now.Add(2*time.Hour), 1*time.Hour)

		intersection := r1.Intersection(r2)
		expected := naturaltime.NewRange(now.Add(2*time.Hour), 0)

		if !intersection.Equal(expected) {
			t.Errorf("Intersection of non-overlapping ranges = %v, want zero-duration range at %v",
				intersection, expected)
		}
	})
}

func TestRangeUnion(t *testing.T) {
	now := time.Now()

	t.Run("Adjacent ranges", func(t *testing.T) {
		r1 := naturaltime.NewRange(now, 1*time.Hour)
		r2 := naturaltime.NewRange(now.Add(1*time.Hour), 1*time.Hour)

		union := r1.Union(r2)
		expected := naturaltime.NewRange(now, 2*time.Hour)

		if !union.Equal(expected) {
			t.Errorf("Union of adjacent ranges = %v, want %v", union, expected)
		}
	})

	t.Run("Overlapping ranges", func(t *testing.T) {
		r1 := naturaltime.NewRange(now, 2*time.Hour)
		r2 := naturaltime.NewRange(now.Add(1*time.Hour), 2*time.Hour)

		union := r1.Union(r2)
		expected := naturaltime.NewRange(now, 3*time.Hour)

		if !union.Equal(expected) {
			t.Errorf("Union of overlapping ranges = %v, want %v", union, expected)
		}
	})

	t.Run("One range contains the other", func(t *testing.T) {
		r1 := naturaltime.NewRange(now, 3*time.Hour)
		r2 := naturaltime.NewRange(now.Add(1*time.Hour), 1*time.Hour)

		union := r1.Union(r2)
		expected := naturaltime.NewRange(now, 3*time.Hour)

		if !union.Equal(expected) {
			t.Errorf("Union when one range contains the other = %v, want %v", union, expected)
		}
	})
}

func TestRangeFromTimes(t *testing.T) {
	now := time.Now()
	start := now
	end := now.Add(2 * time.Hour)

	t.Run("Valid time range", func(t *testing.T) {
		r := naturaltime.RangeFromTimes(start, end)

		if !r.Start().Equal(start) {
			t.Errorf("Start time = %v, want %v", r.Start(), start)
		}

		if !r.End().Equal(end) {
			t.Errorf("End time = %v, want %v", r.End(), end)
		}

		if r.Duration != 2*time.Hour {
			t.Errorf("Duration = %v, want %v", r.Duration, 2*time.Hour)
		}
	})

	t.Run("End before start", func(t *testing.T) {
		invalidEnd := now.Add(-1 * time.Hour)
		r := naturaltime.RangeFromTimes(start, invalidEnd)

		if !r.Start().Equal(start) {
			t.Errorf("Start time = %v, want %v", r.Start(), start)
		}

		if r.Duration != 0 {
			t.Errorf("Duration should be 0 when end is before start, got %v", r.Duration)
		}

		if !r.End().Equal(start) {
			t.Errorf("End time should equal start time when end is before start, got %v, want %v",
				r.End(), start)
		}
	})
}

func TestNewRange(t *testing.T) {
	now := time.Now()

	t.Run("Positive duration", func(t *testing.T) {
		r := naturaltime.NewRange(now, 2*time.Hour)

		if !r.Start().Equal(now) {
			t.Errorf("Start time = %v, want %v", r.Start(), now)
		}

		if r.Duration != 2*time.Hour {
			t.Errorf("Duration = %v, want %v", r.Duration, 2*time.Hour)
		}
	})

	t.Run("Negative duration", func(t *testing.T) {
		r := naturaltime.NewRange(now, -1*time.Hour)

		if !r.Start().Equal(now) {
			t.Errorf("Start time = %v, want %v", r.Start(), now)
		}

		if r.Duration != 0 {
			t.Errorf("Duration should be 0 when given negative duration, got %v", r.Duration)
		}
	})
}

func TestRangeIsAllDay(t *testing.T) {
	now := time.Now()

	t.Run("Zero duration", func(t *testing.T) {
		r := naturaltime.NewRange(now, 0)

		if !r.IsAllDay() {
			t.Errorf("Range with zero duration should be all day")
		}
	})

	t.Run("Non-zero duration", func(t *testing.T) {
		r := naturaltime.NewRange(now, 1*time.Hour)

		if r.IsAllDay() {
			t.Errorf("Range with non-zero duration should not be all day")
		}
	})
}

func TestRangeEqual(t *testing.T) {
	now := time.Now()

	t.Run("Equal ranges", func(t *testing.T) {
		r1 := naturaltime.NewRange(now, 2*time.Hour)
		r2 := naturaltime.NewRange(now, 2*time.Hour)

		if !r1.Equal(r2) {
			t.Errorf("Ranges with same start and duration should be equal")
		}
	})

	t.Run("Different start times", func(t *testing.T) {
		r1 := naturaltime.NewRange(now, 2*time.Hour)
		r2 := naturaltime.NewRange(now.Add(1*time.Minute), 2*time.Hour)

		if r1.Equal(r2) {
			t.Errorf("Ranges with different start times should not be equal")
		}
	})

	t.Run("Different durations", func(t *testing.T) {
		r1 := naturaltime.NewRange(now, 2*time.Hour)
		r2 := naturaltime.NewRange(now, 3*time.Hour)

		if r1.Equal(r2) {
			t.Errorf("Ranges with different durations should not be equal")
		}
	})
}

func TestRangeString(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2023-01-15T12:00:00Z")
	r := naturaltime.NewRange(now, 2*time.Hour)

	expected := "[2023-01-15T12:00:00Z, 2023-01-15T14:00:00Z)"
	if r.String() != expected {
		t.Errorf("String() = %q, want %q", r.String(), expected)
	}
}
