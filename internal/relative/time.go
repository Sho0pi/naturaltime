package relative

import "time"

// TimeModifierFunc defines a function that modifies a given time based on a modifier (e.g., +1 or -1)
type TimeModifierFunc func(t time.Time, modifier int) time.Time

// TimeShiftFunc is a function that shifts time without requiring a modifier
type TimeShiftFunc func(t time.Time) time.Time

// timeShiftFunctions maps human-readable time units to their corresponding time modification functions
var timeShiftFunctions = map[string]TimeModifierFunc{
	"seconds": func(t time.Time, modifier int) time.Time { return t.Add(time.Second * time.Duration(modifier)) },
	"minutes": func(t time.Time, modifier int) time.Time { return t.Add(time.Minute * time.Duration(modifier)) },
	"hours":   func(t time.Time, modifier int) time.Time { return t.Add(time.Hour * time.Duration(modifier)) },
	"days":    func(t time.Time, modifier int) time.Time { return t.AddDate(0, 0, modifier) },
	"months":  func(t time.Time, modifier int) time.Time { return t.AddDate(0, modifier, 0) },
	"years":   func(t time.Time, modifier int) time.Time { return t.AddDate(modifier, 0, 0) },
}

// timeUnitAliases maps various aliases of time units to a standardized unit name
var timeUnitAliases = map[string][]string{
	"seconds": {"s", "sec", "second", "seconds"},
	"minutes": {"m", "min", "mins", "minute", "minutes"},
	"hours":   {"h", "hr", "hrs", "hour", "hours"},
	"days":    {"d", "day", "days"},
	"weeks":   {"w", "week", "weeks"},
	"months":  {"mo", "mon", "mos", "month", "months"},
	"years":   {"y", "yr", "year", "years"},
}

// aliasToStandardUnit provides quick lookups to get a standardized unit from an alias
var aliasToStandardUnit = make(map[string]string)

func init() {
	// Populate alias-to-standard unit lookup table
	for unit, aliases := range timeUnitAliases {
		for _, alias := range aliases {
			aliasToStandardUnit[alias] = unit
		}
	}
}
