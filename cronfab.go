package cronfab

import (
	"time"
)

type SecondUnit struct{}

func (SecondUnit) String() string {
	return "second"
}

func (SecondUnit) Trunc(t time.Time) time.Time {
	return t.Round(time.Second)
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

type MinuteUnit struct{}

func (MinuteUnit) String() string {
	return "minute"
}

func (MinuteUnit) Trunc(t time.Time) time.Time {
	return t.Truncate(time.Minute)
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

type HourUnit struct{}

func (HourUnit) String() string {
	return "hour"
}

func (HourUnit) Trunc(t time.Time) time.Time {
	return t.Truncate(time.Hour)
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

type DayUnit struct{}

func (DayUnit) String() string {
	return "day"
}

func (DayUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
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

type MonthUnit struct{}

func (MonthUnit) String() string {
	return "month"
}

func (MonthUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, t.Location())
}

func (MonthUnit) Add(t time.Time, n int) time.Time {
	return t.AddDate(0, n, 0)
}

func (MonthUnit) Less(u Unit) bool {
	switch u.(type) {
	case SecondUnit, MinuteUnit, HourUnit, DayUnit, MonthUnit:
		return false
	}
	return true
}

type YearUnit struct{}

func (YearUnit) String() string {
	return "year"
}

func (YearUnit) Add(t time.Time, n int) time.Time {
	return t.AddDate(n, 0, 0)
}

func (YearUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), 0, 0, 0, 0, 0, 0, t.Location())
}

func (YearUnit) Less(u Unit) bool {
	switch u.(type) {
	case SecondUnit, MinuteUnit, HourUnit, DayUnit, YearUnit:
		return false
	}
	return true
}

var DefaultContabConfig = NewCrontabConfig([]FieldConfig{
	{
		unit: MinuteUnit{},
		name: "minute",
		min:  0,
		max:  59,
		dateIndexFn: func(t time.Time) int {
			return t.Minute()
		},
	},
	{
		unit: HourUnit{},
		name: "hour",
		min:  0,
		max:  23,
		dateIndexFn: func(t time.Time) int {
			return t.Hour()
		},
	},
	{
		unit: DayUnit{},
		name: "day of month",
		min:  1,
		max:  31,
		dateIndexFn: func(t time.Time) int {
			return t.Day()
		},
	},
	{
		unit:       MonthUnit{},
		name:       "month",
		rangeNames: []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"},
		min:        1,
		max:        12,
		dateIndexFn: func(t time.Time) int {
			return int(t.Month())
		},
	},
	{
		unit:       DayUnit{},
		name:       "day of week",
		rangeNames: []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"},
		min:        0,
		max:        6,
		dateIndexFn: func(t time.Time) int {
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
		dateIndexFn: func(t time.Time) int {
			return t.Second()
		},
	},
	{
		unit: MinuteUnit{},
		name: "minute",
		min:  0,
		max:  59,
		dateIndexFn: func(t time.Time) int {
			return t.Minute()
		},
	},
	{
		unit: HourUnit{},
		name: "hour",
		min:  0,
		max:  23,
		dateIndexFn: func(t time.Time) int {
			return t.Hour()
		},
	},
	{
		unit: DayUnit{},
		name: "day of month",
		min:  1,
		max:  31,
		dateIndexFn: func(t time.Time) int {
			return t.Day()
		},
	},
	{
		unit:       MonthUnit{},
		name:       "month",
		rangeNames: []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"},
		min:        1,
		max:        12,
		dateIndexFn: func(t time.Time) int {
			return int(t.Month())
		},
	},
	{
		unit:       DayUnit{},
		name:       "day of week",
		rangeNames: []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"},
		min:        0,
		max:        6,
		dateIndexFn: func(t time.Time) int {
			return int(t.Weekday())
		},
	},
})
