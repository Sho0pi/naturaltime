package naturaltime

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dop251/goja"
)

//go:embed dist/naturaltime.out.js
var naturaltimeJavaScript string

// TimeRangeResult represents the JSON structure returned from the JavaScript parser.
type timeRangeResult struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end,omitzero"`
}

// Parser provides natural language parsing capabilities for time expressions.
// It uses a JavaScript implementation embedded in the Go binary.
type Parser struct {
	runtime        *goja.Runtime
	parseRangeFunc goja.Callable
	parseDateFunc  goja.Callable
	thisContext    goja.Value
}

// New creates a new natural time expression parser.
// It initializes the JavaScript runtime and prepares the parsing functions.
//
// Returns:
//   - An initialized Parser
//   - An error if initialization fails
func New() (*Parser, error) {
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
	}, nil
}

// ParseDate parses a natural language date expression and returns the corresponding time.
//
// Parameters:
//   - expr: The natural language expression to parse
//   - base: The reference time for relative expressions
//
// Returns:
//   - A pointer to the parsed time.Time, or nil if the expression could not be parsed
//   - An error if parsing fails
func (p *Parser) ParseDate(expr string, base time.Time) (*time.Time, error) {
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

// ParseRange parses a natural language time expression and returns a single time range.
//
// Parameters:
//   - expr: The natural language expression to parse
//   - base: The reference time for relative expressions
//
// Returns:
//   - A pointer to the parsed Range
//   - An error if parsing fails or if multiple ranges are returned
func (p *Parser) ParseRange(expr string, base time.Time) (*Range, error) {
	ranges, err := p.ParseMulti(expr, base)
	if err != nil {
		return nil, err
	}

	if len(ranges) != 1 {
		return nil, fmt.Errorf("expected a single range from expression %q, got %d", expr, len(ranges))
	}

	return &ranges[0], nil
}

// ParseMulti parses a natural language time expression that may represent multiple time ranges.
//
// Parameters:
//   - expr: The natural language expression to parse
//   - base: The reference time for relative expressions
//
// Returns:
//   - A slice of Range objects
//   - An error if parsing fails
func (p *Parser) ParseMulti(expr string, base time.Time) ([]Range, error) {
	result, err := p.parseRangeFunc(p.thisContext, p.runtime.ToValue(expr), p.runtime.ToValue(base.Format(time.RFC3339)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse range expression %q: %w", expr, err)
	}

	// Convert JavaScript result to Go
	jsonBytes, err := json.Marshal(result.Export())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal results for expression %q: %w", expr, err)
	}

	var timeRangeResults []timeRangeResult
	err = json.Unmarshal(jsonBytes, &timeRangeResults)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal results for expression %q: %w", expr, err)
	}

	// Convert time range results to Range objects
	ranges := make([]Range, 0, len(timeRangeResults))
	for _, rangeResult := range timeRangeResults {
		ranges = append(ranges, RangeFromTimes(rangeResult.Start, rangeResult.End))
	}

	return ranges, nil
}
