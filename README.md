
# NaturalTime 🕒 (Native Go Implementation)

[![Go Reference](https://pkg.go.dev/badge/github.com/sho0pi/naturaltime.svg)](https://pkg.go.dev/github.com/sho0pi/naturaltime)
[![Go Report Card](https://goreportcard.com/badge/github.com/sho0pi/naturaltime)](https://goreportcard.com/report/github.com/sho0pi/naturaltime)

**A pure Go library for parsing natural language time expressions!** This native Go implementation eliminates JavaScript dependencies while maintaining advanced time parsing capabilities.

## ✨ Features

- 100% Go implementation - no JavaScript runtime!
- Parse natural language dates into `time.Time` objects
- Extract specific dates or **date ranges**
- Support for multiple time ranges in single expressions
- Idiomatic Go API with zero external dependencies
- Active development - contributions welcome!

## 📦 Installation

```bash
go get github.com/sho0pi/naturaltime
```

## 🚀 Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/sho0pi/naturaltime"
)

func main() {
	// Create a new parser
	parser, err := naturaltime.New()
	if err != nil {
		panic(err)
	}
	
	now := time.Now()
	
	// Example 1: Parse a simple date expression
	date, err := parser.ParseDate("tomorrow at 3pm", now)
	if err != nil {
		panic(err)
	}
	if date != nil {
		fmt.Println(date)
	}
	
	// Example 2: Parse a time range expression
	timeRange, err := parser.ParseRange("from 2pm to 4pm tomorrow", now)
	if err != nil {
		panic(err)
	}
	fmt.Printf("  Start: %s, End: %s", timeRange.Start(), timeRange.End())
	fmt.Printf("  Duration: %s\n\n", timeRange.Duration)
	
	// Example 3: Parse multiple time ranges
	ranges, err := parser.ParseMulti("Monday and Tuesday from 9am to 5pm", now)
	if err != nil {
		panic(err)
	}
	for i, r := range ranges {
		fmt.Printf("  Range %d: %s to %s (%s)\n", 
			i+1, 
			r.Start(), 
			r.End(),
			r.Duration)
	}
	
}
```

## Supported Expressions

The library can parse a wide variety of natural language time expressions, including:

- Relative dates: "today", "tomorrow", "next week"
- Specific dates: "January 15, 2023", "15/01/2023"
- Time expressions: "3pm", "15:00"
- Durations: "from 2pm to 4pm", "for 2 hours"
- Combined expressions: "tomorrow from 9am to 5pm"
- Recurring times: "every Monday", "every weekday"

## 🛠 Roadmap & Contribution

We're actively working to achieve feature parity with the JS-based version! Help us:
- Add new time expression patterns
- Improve range detection
- Optimize parsing performance
- Expand test coverage

Check our [issues](https://github.com/sho0pi/naturaltime/issues) for good first contributions!

## License

MIT
