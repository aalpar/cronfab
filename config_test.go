package cronfab

import (
	"testing"
	"time"
)

// DayUnit units in days
type testDayUnit struct{}

func (testDayUnit) String() string {
	return "day"
}

func (testDayUnit) Add(t time.Time, n int) time.Time {
	return t.AddDate(0, 0, n)
}

func (testDayUnit) Less(u Unit) bool {
	switch u.(type) {
	case testDayUnit:
		return false
	}
	return true
}

func (testDayUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

type testMonthUnit struct{}

func (testMonthUnit) String() string {
	return "month"
}

func (testMonthUnit) Add(t time.Time, n int) time.Time {
	return t.AddDate(0, n, 0)
}

func (testMonthUnit) Less(u Unit) bool {
	switch u.(type) {
	case testDayUnit, testMonthUnit:
		return false
	}
	return true
}

func (testMonthUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, t.Location())
}

type testYearUnit struct{}

func (testYearUnit) String() string {
	return "year"
}

func (testYearUnit) Add(t time.Time, n int) time.Time {
	return t.AddDate(n, 0, 0)
}

func (testYearUnit) Less(u Unit) bool {
	switch u.(type) {
	case testDayUnit, testMonthUnit, testYearUnit:
		return false
	}
	return true
}

func (testYearUnit) Trunc(t time.Time) time.Time {
	return time.Date(t.Year(), 0, 0, 0, 0, 0, 0, t.Location())
}

func TestNewCrontabConfig(t *testing.T) {
	var testContabConfig = NewCrontabConfig([]FieldConfig{
		{
			unit: testDayUnit{},
			name: "day of month",
			min:  1,
			max:  31,
			getIndex: func(t time.Time) int {
				return t.Day()
			},
		},
		{
			unit:       testMonthUnit{},
			name:       "month",
			rangeNames: []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"},
			min:        1,
			max:        12,
			getIndex: func(t time.Time) int {
				return int(t.Month())
			},
		},
		{
			unit:       testDayUnit{},
			name:       "day of week",
			rangeNames: []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"},
			min:        0,
			max:        6,
			getIndex: func(t time.Time) int {
				return int(t.Weekday())
			},
		},
	})
	if len(testContabConfig.Units) != 2 {
		t.Fatalf("unexpected value.")
	}
	if testContabConfig.Units[0] != Unit(testDayUnit{}) {
		t.Fatalf("unexpected value.")
	}
	if testContabConfig.Units[1] != Unit(testMonthUnit{}) {
		t.Fatalf("unexpected value.")
	}
}
