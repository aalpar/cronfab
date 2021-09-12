package cronfab

import (
	"time"
)

// SecondUnit units in seconds
type SecondUnit struct{}

func (SecondUnit) String() string {
	return "second"
}

func (SecondUnit) Add(t time.Time, n int) time.Time {
	return t.Add(time.Duration(n) * time.Second)
}

func (SecondUnit) Less(u Unit) bool {
	switch u.(type) {
	case SecondUnit:
		return false
	}
	return true
}

func (SecondUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

// MinuteUnit units in minutes
type MinuteUnit struct{}

func (MinuteUnit) String() string {
	return "minute"
}

func (MinuteUnit) Add(t time.Time, n int) time.Time {
	return t.Add(time.Duration(n) * time.Minute)
}

func (MinuteUnit) Less(u Unit) bool {
	switch u.(type) {
	case SecondUnit, MinuteUnit:
		return false
	}
	return true
}

func (MinuteUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

// HourUnit units in hours
type HourUnit struct{}

func (HourUnit) String() string {
	return "hour"
}

func (x HourUnit) Less(u Unit) bool {
	switch u.(type) {
	case SecondUnit, MinuteUnit, HourUnit:
		return false
	}
	return true
}

func (HourUnit) Add(t time.Time, n int) time.Time {
	return t.Add(time.Duration(n) * time.Hour)
}

func (HourUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

// DayUnit units in days
type DayUnit struct{}

func (DayUnit) String() string {
	return "day"
}

func (DayUnit) Add(t time.Time, n int) time.Time {
	return t.AddDate(0, 0, n)
}

func (DayUnit) Less(u Unit) bool {
	switch u.(type) {
	case SecondUnit, MinuteUnit, HourUnit, DayUnit:
		return false
	}
	return true
}

func (DayUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// WeekUnit units in days
type WeekOfMonth struct{}

func (WeekOfMonth) String() string {
	return "week"
}

func (WeekOfMonth) Add(t time.Time, n int) time.Time {
	return t.AddDate(0, 0, n*7)
}

func (WeekOfMonth) Less(u Unit) bool {
	switch u.(type) {
	case SecondUnit, MinuteUnit, HourUnit, DayUnit, WeekOfMonth:
		return false
	}
	return true
}

func (WeekOfMonth) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, t.Location())
}

// MonthUnit units in months
type MonthUnit struct{}

func (MonthUnit) String() string {
	return "month"
}

func (MonthUnit) Add(t time.Time, n int) time.Time {
	return t.AddDate(0, n, 0)
}

func (MonthUnit) Less(u Unit) bool {
	switch u.(type) {
	case SecondUnit, MinuteUnit, HourUnit, DayUnit, WeekOfMonth, MonthUnit:
		return false
	}
	return true
}

func (MonthUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, t.Location())
}

// YearUnit units in years
type YearUnit struct{}

func (YearUnit) String() string {
	return "year"
}

func (YearUnit) Add(t time.Time, n int) time.Time {
	return t.AddDate(n, 0, 0)
}

func (YearUnit) Less(u Unit) bool {
	switch u.(type) {
	case SecondUnit, MinuteUnit, HourUnit, DayUnit, WeekOfMonth, MonthUnit, YearUnit:
		return false
	}
	return true
}

func (YearUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), 0, 0, 0, 0, 0, 0, t.Location())
}

var DefaultContabConfig = NewCrontabConfig([]FieldConfig{
	{
		unit: MinuteUnit{},
		name: "minute",
		min:  0,
		max:  59,
		getIndex: func(t time.Time) int {
			return t.Minute()
		},
	},
	{
		unit: HourUnit{},
		name: "hour",
		min:  0,
		max:  23,
		getIndex: func(t time.Time) int {
			return t.Hour()
		},
	},
	{
		unit: DayUnit{},
		name: "day of month",
		min:  1,
		max:  31,
		getIndex: func(t time.Time) int {
			return t.Day()
		},
	},
	{
		unit:       MonthUnit{},
		name:       "month",
		rangeNames: []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"},
		min:        1,
		max:        12,
		getIndex: func(t time.Time) int {
			return int(t.Month())
		},
	},
	{
		unit:       DayUnit{},
		name:       "day of week",
		rangeNames: []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"},
		min:        0,
		max:        6,
		getIndex: func(t time.Time) int {
			return int(t.Weekday())
		},
	},
})

var SecondContabConfig = NewCrontabConfig([]FieldConfig{
	{
		unit: SecondUnit{},
		name: "second",
		min:  0,
		max:  59,
		getIndex: func(t time.Time) int {
			return t.Second()
		},
	},
	{
		unit: MinuteUnit{},
		name: "minute",
		min:  0,
		max:  59,
		getIndex: func(t time.Time) int {
			return t.Minute()
		},
	},
	{
		unit: HourUnit{},
		name: "hour",
		min:  0,
		max:  23,
		getIndex: func(t time.Time) int {
			return t.Hour()
		},
	},
	{
		unit: DayUnit{},
		name: "day of month",
		min:  1,
		max:  31,
		getIndex: func(t time.Time) int {
			return t.Day()
		},
	},
	{
		unit: WeekOfMonth{},
		name: "week of month",
		min:  1,
		max:  5,
		getIndex: func(t time.Time) int {
			q := (t.Day() + (6 - int(t.Weekday())))
			return (q / 7) + 1
		},
	},
	{
		unit:       MonthUnit{},
		name:       "month",
		rangeNames: []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"},
		min:        1,
		max:        12,
		getIndex: func(t time.Time) int {
			return int(t.Month())
		},
	},
	{
		unit:       DayUnit{},
		name:       "day of week",
		rangeNames: []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"},
		min:        0,
		max:        6,
		getIndex: func(t time.Time) int {
			return int(t.Weekday())
		},
	},
})
