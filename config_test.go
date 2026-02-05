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
	testCrontabConfig, err := NewCrontabConfig([]FieldConfig{
		{
			Unit: testDayUnit{},
			Name: "day of month",
			Min:  1,
			Max:  31,
			GetIndex: func(t time.Time) int {
				return t.Day()
			},
		},
		{
			Unit:       testMonthUnit{},
			Name:       "month",
			RangeNames: []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"},
			Min:        1,
			Max:        12,
			GetIndex: func(t time.Time) int {
				return int(t.Month())
			},
		},
		{
			Unit:       testDayUnit{},
			Name:       "day of week",
			RangeNames: []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"},
			Min:        0,
			Max:        6,
			GetIndex: func(t time.Time) int {
				return int(t.Weekday())
			},
		},
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(testCrontabConfig.Units) != 2 {
		t.Fatalf("unexpected value.")
	}
	if testCrontabConfig.Units[0] != Unit(testDayUnit{}) {
		t.Fatalf("unexpected value.")
	}
	if testCrontabConfig.Units[1] != Unit(testMonthUnit{}) {
		t.Fatalf("unexpected value.")
	}
}

func TestNewCrontabConfig_Validation(t *testing.T) {
	valid := func(t time.Time) int { return 0 }

	cases := []struct {
		name   string
		fields []FieldConfig
	}{
		{"no fields", nil},
		{"nil Unit", []FieldConfig{{Unit: nil, Name: "x", Min: 0, Max: 1, GetIndex: valid}}},
		{"empty Name", []FieldConfig{{Unit: testDayUnit{}, Name: "", Min: 0, Max: 1, GetIndex: valid}}},
		{"nil GetIndex", []FieldConfig{{Unit: testDayUnit{}, Name: "x", Min: 0, Max: 1, GetIndex: nil}}},
		{"Min > Max", []FieldConfig{{Unit: testDayUnit{}, Name: "x", Min: 10, Max: 5, GetIndex: valid}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCrontabConfig(tc.fields)
			if err == nil {
				t.Errorf("expected error for %s", tc.name)
			}
		})
	}
}
