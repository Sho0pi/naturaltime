package naturaltime

import (
	"fmt"
	"time"
)

// Range is a time range represented as a half-open interval [start, end).
// The start time is inclusive, and the end time is exclusive.
type Range struct {
	start    time.Time
	Duration time.Duration
}

// Start returns when the range begins (inclusive).
func (r Range) Start() time.Time {
	return r.start
}

// End returns when the range ends (exclusive).
func (r Range) End() time.Time {
	return r.start.Add(r.Duration)
}

// Equal returns true if the two ranges are exactly equal,
// having the same start time and duration.
func (r Range) Equal(other Range) bool {
	return r.start.Equal(other.start) && r.Duration == other.Duration
}

// Contains checks if the given time is within this range.
// The range is half-open: [start, end), so start <= t < end.
func (r Range) Contains(t time.Time) bool {
	return !t.Before(r.start) && t.Before(r.End())
}

// ContainsRange checks if another range is fully contained within this range.
func (r Range) ContainsRange(other Range) bool {
	return !other.start.Before(r.start) && !other.End().After(r.End())
}

// Overlaps checks if this range overlaps with another range.
// Two ranges overlap if they share any common time point.
// Note: If 'other' starts exactly at the same time as r.End(), it is not considered an overlap.
// Similarly, if r starts exactly at the same time as other.End(), it is also not considered an overlap.
func (r Range) Overlaps(other Range) bool {
	return !r.End().Before(other.start) && !r.start.After(other.End()) && !r.End().Equal(other.Start()) && !other.End().Equal(r.Start())
}

// Intersection returns the overlapping range between two ranges,
// or a zero-duration range at the latest start time if they don't overlap.
func (r Range) Intersection(other Range) Range {
	if !r.Overlaps(other) {
		// Return zero duration range at the later of the two start times
		if r.start.After(other.start) {
			return Range{r.start, 0}
		}
		return Range{other.start, 0}
	}

	// Find the latest start time
	start := r.start
	if other.start.After(start) {
		start = other.start
	}

	// Find the earliest end time
	end := r.End()
	if other.End().Before(end) {
		end = other.End()
	}

	return Range{start, end.Sub(start)}
}

// Union returns the smallest range that contains both ranges.
// Note: This only makes sense if the ranges overlap or are adjacent.
func (r Range) Union(other Range) Range {
	// Find the earliest start time
	start := r.start
	if other.start.Before(start) {
		start = other.start
	}

	// Find the latest end time
	end := r.End()
	if other.End().After(end) {
		end = other.End()
	}

	return Range{start, end.Sub(start)}
}

// IsAllDay returns true if the range has zero duration.
func (r Range) IsAllDay() bool {
	return r.Duration == 0
}

// String returns a human-readable representation of the range.
func (r Range) String() string {
	return fmt.Sprintf("[%s, %s)", r.start.Format(time.RFC3339), r.End().Format(time.RFC3339))
}

// RangeFromTimes creates a Range from two times.
// The start time should be before or equal to the end time.
func RangeFromTimes(start, end time.Time) Range {
	if end.Before(start) {
		return Range{start, 0} // Return an empty range if end is before start
	}
	return Range{start, end.Sub(start)}
}

// NewRange creates a new Range with the given start time and duration.
func NewRange(start time.Time, duration time.Duration) Range {
	if duration < 0 {
		duration = 0 // Ensure non-negative duration
	}
	return Range{start, duration}
}
