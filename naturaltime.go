package naturaltime

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/dop251/goja"
)

var naturaltimeJavaScript string

// Parser provides natural language parsing capabilities for time expressions.
type Parser struct {
	runtime        *goja.Runtime
	parseRangeFunc goja.Callable
	parseDateFunc  goja.Callable
	thisContext    goja.Value
	nativeParser   *NativeParser // Add native parser
}

// New creates a new natural time expression parser.
func New() (*Parser, error) {
	// Initialize the native parser
	nativeParser := NewNativeParser()

	// Initialize the JavaScript runtime
	runtime := goja.New()

	// Compile and run the embedded JavaScript
	program, err := goja.Compile("naturaltime.js", naturaltimeJavaScript, false)
	if err != nil {
		return nil, fmt.Errorf("failed to compile naturaltime JavaScript: %w", err)
	}

	_, err = runtime.RunProgram(program)
	if err != nil {
		return nil, fmt.Errorf("failed to run naturaltime JavaScript: %w", err)
	}

	// Extract the JavaScript object and its methods
	jsObject := runtime.Get("naturaltime").ToObject(runtime)

	parseRangeFunc, ok := goja.AssertFunction(jsObject.Get("parseRange"))
	if !ok {
		return nil, fmt.Errorf("failed to get 'parseRange' function from JavaScript")
	}

	parseDateFunc, ok := goja.AssertFunction(jsObject.Get("parseDate"))
	if !ok {
		return nil, fmt.Errorf("failed to get 'parseDate' function from JavaScript")
	}

	return &Parser{
		thisContext:    runtime.ToValue(map[string]interface{}{}),
		runtime:        runtime,
		parseRangeFunc: parseRangeFunc,
		parseDateFunc:  parseDateFunc,
		nativeParser:   nativeParser, // Set the native parser
	}, nil
}

// ParseDate parses a natural language date expression and returns the corresponding time.
func (p *Parser) ParseDate(expr string, base time.Time) (*time.Time, error) {
	// Try the native parser first
	if date, err := p.nativeParser.ParseDate(expr); err == nil {
		return date, nil
	}

	// Fall back to the JavaScript parser
	result, err := p.parseDateFunc(p.thisContext, p.runtime.ToValue(expr), p.runtime.ToValue(base.Format(time.RFC3339)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse date expression %q: %w", expr, err)
	}

	switch parsedValue := result.Export().(type) {
	case time.Time:
		return &parsedValue, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected result type when parsing date expression %q", expr)
	}
}
