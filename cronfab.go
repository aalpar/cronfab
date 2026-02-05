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

// WeekOfMonth units in week in month ordinal
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
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
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

var DefaultCrontabConfig = MustCrontabConfig([]FieldConfig{
	{
		Unit: MinuteUnit{},
		Name: "minute",
		Min:  0,
		Max:  59,
		GetIndex: func(t time.Time) int {
			return t.Minute()
		},
	},
	{
		Unit: HourUnit{},
		Name: "hour",
		Min:  0,
		Max:  23,
		GetIndex: func(t time.Time) int {
			return t.Hour()
		},
	},
	{
		Unit: DayUnit{},
		Name: "day of month",
		Min:  1,
		Max:  31,
		GetIndex: func(t time.Time) int {
			return t.Day()
		},
	},
	{
		Unit:       MonthUnit{},
		Name:       "month",
		RangeNames: []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"},
		Min:        1,
		Max:        12,
		GetIndex: func(t time.Time) int {
			return int(t.Month())
		},
	},
	{
		Unit:       DayUnit{},
		Name:       "day of week",
		RangeNames: []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"},
		Min:        0,
		Max:        6,
		GetIndex: func(t time.Time) int {
			return int(t.Weekday())
		},
	},
})

var SecondCrontabConfig = MustCrontabConfig([]FieldConfig{
	{
		Unit: SecondUnit{},
		Name: "second",
		Min:  0,
		Max:  59,
		GetIndex: func(t time.Time) int {
			return t.Second()
		},
	},
	{
		Unit: MinuteUnit{},
		Name: "minute",
		Min:  0,
		Max:  59,
		GetIndex: func(t time.Time) int {
			return t.Minute()
		},
	},
	{
		Unit: HourUnit{},
		Name: "hour",
		Min:  0,
		Max:  23,
		GetIndex: func(t time.Time) int {
			return t.Hour()
		},
	},
	{
		Unit: DayUnit{},
		Name: "day of month",
		Min:  1,
		Max:  31,
		GetIndex: func(t time.Time) int {
			return t.Day()
		},
	},
	{
		Unit: WeekOfMonth{},
		Name: "week of month",
		Min:  1,
		Max:  5,
		GetIndex: func(t time.Time) int {
			q := t.Day() + (6 - int(t.Weekday()))
			return (q / 7) + 1
		},
	},
	{
		Unit:       MonthUnit{},
		Name:       "month",
		RangeNames: []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"},
		Min:        1,
		Max:        12,
		GetIndex: func(t time.Time) int {
			return int(t.Month())
		},
	},
	{
		Unit:       DayUnit{},
		Name:       "day of week",
		RangeNames: []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"},
		Min:        0,
		Max:        6,
		GetIndex: func(t time.Time) int {
			return int(t.Weekday())
		},
	},
})
