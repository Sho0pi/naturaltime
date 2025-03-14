package naturaltime

import (
	"fmt"
	"regexp"
	"time"
)

// NativeParser provides natural language parsing capabilities for time expressions in pure Go.
type NativeParser struct {
	now time.Time
}

// NewNativeParser creates a new NativeParser with the current time as the reference.
func NewNativeParser() *NativeParser {
	return &NativeParser{
		now: time.Now(),
	}
}

// ParseDate parses a natural language date expression and returns the corresponding time.
func (p *NativeParser) ParseDate(expr string) (*time.Time, error) {
	// Try to parse exact dates first
	if date, err := p.parseExactDate(expr); err == nil {
		return date, nil
	}

	// Try to parse relative dates
	if date, err := p.parseRelativeDate(expr); err == nil {
		return date, nil
	}

	return nil, fmt.Errorf("unable to parse date expression: %s", expr)
}

// parseExactDate parses exact date formats.
func (p *NativeParser) parseExactDate(expr string) (*time.Time, error) {
	// Define regex patterns for common date formats
	patterns := []struct {
		pattern string
		layout  string
	}{
		{`(\d{4})-(\d{2})-(\d{2})`, "2006-01-02"},                 // YYYY-MM-DD
		{`(\d{2})/(\d{2})/(\d{4})`, "01/02/2006"},                 // MM/DD/YYYY
		{`(\d{2})/(\d{2})/(\d{2})`, "01/02/06"},                   // MM/DD/YY
		{`(\d{1,2})\s+([A-Za-z]+)\s+(\d{4})`, "2 January 2006"},   // 15 January 2023
		{`([A-Za-z]+)\s+(\d{1,2}),\s+(\d{4})`, "January 2, 2006"}, // January 15, 2023
		{`(\d{2})/(\d{2})/(\d{4})`, "02/01/2006"},                 // DD/MM/YYYY
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.pattern)
		if matches := re.FindStringSubmatch(expr); matches != nil {
			// Parse the date using the corresponding layout
			date, err := time.Parse(pattern.layout, expr)
			if err == nil {
				return &date, nil
			}
		}
	}

	return nil, fmt.Errorf("no exact date format matched")
}

// parseRelativeDate parses relative date expressions.
func (p *NativeParser) parseRelativeDate(expr string) (*time.Time, error) {
	// Truncate the current time to seconds to avoid nanosecond differences
	now := p.now.Truncate(time.Second)

	switch expr {
	case "today":
		today := now
		return &today, nil
	case "tomorrow":
		tomorrow := now.AddDate(0, 0, 1)
		return &tomorrow, nil
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		return &yesterday, nil
	case "next week":
		nextWeek := now.AddDate(0, 0, 7)
		return &nextWeek, nil
	default:
		return nil, fmt.Errorf("unrecognized relative date expression")
	}
}
