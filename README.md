# NaturalTime 🕒

[![Go Reference](https://pkg.go.dev/badge/github.com/sho0pi/naturaltime.svg)](https://pkg.go.dev/github.com/sho0pi/naturaltime)
[![Go Report Card](https://goreportcard.com/badge/github.com/sho0pi/naturaltime)](https://goreportcard.com/report/github.com/sho0pi/naturaltime)


A powerful Go library for parsing natural language time expressions with exceptional support for time ranges! NaturalTime is a wrapper around the excellent [chrono-node](https://github.com/wanasit/chrono) JavaScript library, providing Go developers with advanced natural language time parsing capabilities.

## ✨ Features

- Parse natural language date expressions into `time.Time` objects
- Extract specific dates or **date ranges**
- Support for multiple time ranges in a single expression
- Integration with Go's `time` package
- Clean, idiomatic Go API

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

## How It Works

The `naturaltime` library embeds the JavaScript code from chrono-node and executes it within a Go application using the [goja](https://github.com/dop251/goja) JavaScript runtime. This approach provides the rich natural language parsing capabilities of chrono-node while maintaining a pure Go API.

## License

MIT

## Acknowledgments

- [chrono-node](https://github.com/wanasit/chrono) - The underlying JavaScript natural language date parser
- [goja](https://github.com/dop251/goja) - The JavaScript runtime in Go