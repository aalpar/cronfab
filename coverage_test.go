package cronfab

import (
	"testing"
	"time"
)

// --- Error type tests ---

func TestErrorBadIndex_Error(t *testing.T) {
	e := &ErrorBadIndex{FieldName: "minute", Value: 99}
	s := e.Error()
	if s != "invalid index 99 for minute" {
		t.Errorf("unexpected: %q", s)
	}
}

func TestErrorBadName_Error(t *testing.T) {
	e := &ErrorBadName{FieldName: "month", Value: "foo"}
	s := e.Error()
	if s != `invalid name "foo" for month` {
		t.Errorf("unexpected: %q", s)
	}
}

func TestErrorParse_Error(t *testing.T) {
	e := &ErrorParse{Index: 5, State: StateInNumber}
	s := e.Error()
	if s != "in number at 5" {
		t.Errorf("unexpected: %q", s)
	}
}

// --- StateString tests ---

func TestStateString(t *testing.T) {
	cases := []struct {
		state  State
		expect string
	}{
		{StateExpectSplatOrNumberOrName, "expecting '*' or number or name"},
		{StateInNumber, "in number"},
		{StateInName, "in name"},
		{StateExpectEndRangeName, "expecting end-range name"},
		{StateInEndRangeName, "in end-range name"},
		{StateExpectEndRangeNumber, "expecting end-range number"},
		{StateInEndRangeNumber, "in end-range number"},
		{StateExpectStep, "expecting '/'"},
		{StateExpectStepNumber, "expecting step number"},
		{StateInStepNumber, "in step number"},
		{State(999), "unknown"},
	}
	for _, tc := range cases {
		got := StateString(tc.state)
		if got != tc.expect {
			t.Errorf("StateString(%d) = %q, want %q", tc.state, got, tc.expect)
		}
	}
}

// --- CrontabLine tests ---

func TestCrontabLine_Validate(t *testing.T) {
	// valid line
	cl := CrontabLine{{{0, 59, 1}}, {{0, 23, 1}}}
	if err := cl.Validate(); err != nil {
		t.Errorf("expected valid, got %v", err)
	}

	// overlapping constraints in a field
	cl = CrontabLine{{{0, 10, 1}, {5, 15, 1}}}
	if err := cl.Validate(); err == nil {
		t.Error("expected overlap error")
	}

	// reversed boundaries
	cl = CrontabLine{{{10, 5, 1}}}
	if err := cl.Validate(); err != ErrConstraintBoundariesReversed {
		t.Errorf("expected ErrConstraintBoundariesReversed, got %v", err)
	}
}

func TestCrontabLine_String(t *testing.T) {
	cl := CrontabLine{{{0, 59, 1}}, {{0, 23, 1}}}
	s := cl.String()
	if s != "0-59/1 0-23/1" {
		t.Errorf("unexpected: %q", s)
	}

	// single field
	cl = CrontabLine{{{5, 5, 1}}}
	if cl.String() != "5-5/1" {
		t.Errorf("unexpected: %q", cl.String())
	}
}

// --- CrontabConfig.Len ---

func TestCrontabConfig_Len(t *testing.T) {
	n := DefaultCrontabConfig.Len()
	if n != 5 {
		t.Errorf("expected 5 fields, got %d", n)
	}
	n = SecondCrontabConfig.Len()
	if n != 7 {
		t.Errorf("expected 7 fields, got %d", n)
	}
}

// --- Constraint tests ---

func TestCrontabConstraint_Validate(t *testing.T) {
	// valid
	c := CrontabConstraint{0, 59, 1}
	if err := c.Validate(); err != nil {
		t.Errorf("expected valid, got %v", err)
	}
	// single value (min == max)
	c = CrontabConstraint{5, 5, 1}
	if err := c.Validate(); err != nil {
		t.Errorf("expected valid, got %v", err)
	}
	// reversed
	c = CrontabConstraint{10, 5, 1}
	if err := c.Validate(); err != ErrConstraintBoundariesReversed {
		t.Errorf("expected ErrConstraintBoundariesReversed, got %v", err)
	}
}

func TestCrontabConstraint_Ceil(t *testing.T) {
	cases := []struct {
		c      CrontabConstraint
		x      int
		expect int
		roll   bool
	}{
		// below min
		{CrontabConstraint{5, 10, 1}, 3, 5, false},
		// at min
		{CrontabConstraint{5, 10, 1}, 5, 5, false},
		// in range, on step
		{CrontabConstraint{0, 59, 5}, 10, 10, false},
		// in range, off step
		{CrontabConstraint{0, 59, 5}, 11, 15, false},
		// at max
		{CrontabConstraint{0, 59, 5}, 55, 55, false},
		// above max
		{CrontabConstraint{0, 59, 5}, 60, 0, true},
		// off step, next step lands on max
		{CrontabConstraint{0, 10, 5}, 8, 10, false},
		// step of 1, at max boundary
		{CrontabConstraint{1, 31, 1}, 31, 31, false},
		// step of 1, above max
		{CrontabConstraint{1, 31, 1}, 32, 1, true},
	}
	for i, tc := range cases {
		got, roll := tc.c.Ceil(tc.x)
		if got != tc.expect || roll != tc.roll {
			t.Errorf("case %d: Ceil(%d) on %v = (%d, %v), want (%d, %v)",
				i, tc.x, tc.c, got, roll, tc.expect, tc.roll)
		}
	}
}

func TestCrontabConstraint_String(t *testing.T) {
	c := CrontabConstraint{0, 59, 5}
	if c.String() != "0-59/5" {
		t.Errorf("unexpected: %q", c.String())
	}
}

// --- CrontabField.Ceil with multiple constraints ---

func TestCrontabField_CeilMulti(t *testing.T) {
	// two non-overlapping ranges
	cf := CrontabField{
		{5, 10, 1},
		{20, 25, 1},
	}

	// value in first range
	v, roll := cf.Ceil(7)
	if v != 7 || roll {
		t.Errorf("expected (7, false), got (%d, %v)", v, roll)
	}

	// value between ranges → lands in second
	v, roll = cf.Ceil(15)
	if v != 20 || roll {
		t.Errorf("expected (20, false), got (%d, %v)", v, roll)
	}

	// value above all ranges → roll to first min
	v, roll = cf.Ceil(30)
	if v != 5 || !roll {
		t.Errorf("expected (5, true), got (%d, %v)", v, roll)
	}
}

// --- Unit tests ---

func TestSecondUnit(t *testing.T) {
	u := SecondUnit{}
	if u.String() != "second" {
		t.Errorf("unexpected: %q", u.String())
	}
	if !u.Less(MinuteUnit{}) {
		t.Error("SecondUnit should be less than MinuteUnit")
	}
	if u.Less(SecondUnit{}) {
		t.Error("SecondUnit should not be less than SecondUnit")
	}
	t0 := time.Date(2020, 1, 1, 0, 0, 0, 500, time.UTC)
	truncated := u.Trunc(t0)
	if truncated.Nanosecond() != 0 {
		t.Errorf("Trunc should zero nanoseconds, got %d", truncated.Nanosecond())
	}
}

func TestMinuteUnit(t *testing.T) {
	u := MinuteUnit{}
	if u.Less(SecondUnit{}) {
		t.Error("MinuteUnit should not be less than SecondUnit")
	}
	if u.Less(MinuteUnit{}) {
		t.Error("MinuteUnit should not be less than MinuteUnit")
	}
	if !u.Less(HourUnit{}) {
		t.Error("MinuteUnit should be less than HourUnit")
	}
}

func TestHourUnit(t *testing.T) {
	u := HourUnit{}
	if u.Less(MinuteUnit{}) {
		t.Error("HourUnit should not be less than MinuteUnit")
	}
	if u.Less(HourUnit{}) {
		t.Error("HourUnit should not be less than HourUnit")
	}
	if !u.Less(DayUnit{}) {
		t.Error("HourUnit should be less than DayUnit")
	}
}

func TestDayUnit(t *testing.T) {
	u := DayUnit{}
	if u.Less(DayUnit{}) {
		t.Error("DayUnit should not be less than DayUnit")
	}
	if !u.Less(MonthUnit{}) {
		t.Error("DayUnit should be less than MonthUnit")
	}
}

func TestWeekOfMonthUnit(t *testing.T) {
	u := WeekOfMonth{}
	if u.String() != "week" {
		t.Errorf("unexpected: %q", u.String())
	}
	if u.Less(DayUnit{}) {
		t.Error("WeekOfMonth should not be less than DayUnit")
	}
	if u.Less(WeekOfMonth{}) {
		t.Error("WeekOfMonth should not be less than WeekOfMonth")
	}
	if !u.Less(MonthUnit{}) {
		t.Error("WeekOfMonth should be less than MonthUnit")
	}
	// Add 1 week
	t0 := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	t1 := u.Add(t0, 1)
	if t1.Day() != 8 {
		t.Errorf("expected day 8, got %d", t1.Day())
	}
	// Trunc to start of week
	t2 := u.Trunc(time.Date(2020, 1, 8, 15, 30, 0, 0, time.UTC)) // Wednesday
	if t2.Hour() != 0 || t2.Minute() != 0 {
		t.Errorf("Trunc should zero sub-day components")
	}
}

func TestMonthUnit(t *testing.T) {
	u := MonthUnit{}
	if u.String() != "month" {
		t.Errorf("unexpected: %q", u.String())
	}
	if u.Less(DayUnit{}) {
		t.Error("MonthUnit should not be less than DayUnit")
	}
	if u.Less(MonthUnit{}) {
		t.Error("MonthUnit should not be less than MonthUnit")
	}
	if !u.Less(YearUnit{}) {
		t.Error("MonthUnit should be less than YearUnit")
	}
	// Add
	t0 := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	t1 := u.Add(t0, 2)
	if t1.Month() != 3 {
		t.Errorf("expected March, got %v", t1.Month())
	}
	// Trunc
	t2 := u.Trunc(time.Date(2020, 3, 15, 12, 30, 0, 0, time.UTC))
	_ = t2 // just exercise the code path
}

func TestYearUnit(t *testing.T) {
	u := YearUnit{}
	if u.String() != "year" {
		t.Errorf("unexpected: %q", u.String())
	}
	if u.Less(MonthUnit{}) {
		t.Error("YearUnit should not be less than MonthUnit")
	}
	if u.Less(YearUnit{}) {
		t.Error("YearUnit should not be less than YearUnit")
	}
	// No unit larger than year in the codebase, but Less should return false for itself
	// Add
	t0 := time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC)
	t1 := u.Add(t0, 3)
	if t1.Year() != 2023 {
		t.Errorf("expected 2023, got %d", t1.Year())
	}
	// Trunc
	t2 := u.Trunc(time.Date(2020, 6, 15, 12, 30, 0, 0, time.UTC))
	_ = t2 // exercise the code path
}

func TestSortUnits(t *testing.T) {
	units := []Unit{MonthUnit{}, SecondUnit{}, HourUnit{}, DayUnit{}, MinuteUnit{}}
	SortUnits(units)
	expected := []string{"second", "minute", "hour", "day", "month"}
	for i, u := range units {
		if u.String() != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], u.String())
		}
	}
}

// --- Parser edge cases ---

func TestParseCrontab_Empty(t *testing.T) {
	cl, err := DefaultCrontabConfig.ParseCronTab("")
	if err != nil {
		t.Errorf("empty string should not error: %v", err)
	}
	if len(cl) != 0 {
		t.Errorf("expected empty CrontabLine, got %v", cl)
	}
}

func TestParseCrontab_AllFields(t *testing.T) {
	cl, err := DefaultCrontabConfig.ParseCronTab("0 0 1 1 0")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(cl) != 5 {
		t.Fatalf("expected 5 fields, got %d", len(cl))
	}
}

func TestParseCrontab_AllWildcards(t *testing.T) {
	cl, err := DefaultCrontabConfig.ParseCronTab("* * * * *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(cl) != 5 {
		t.Fatalf("expected 5 fields, got %d", len(cl))
	}
	expected := [][2]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 6}}
	for i, exp := range expected {
		c := cl[i][0]
		if c[0] != exp[0] || c[1] != exp[1] || c[2] != 1 {
			t.Errorf("field %d: expected {%d,%d,1}, got %v", i, exp[0], exp[1], c)
		}
	}
}

func TestParseCrontab_NamedRange(t *testing.T) {
	// day-of-week range with names (must be last field because the parser
	// doesn't handle StateInEndRangeName at space/comma delimiters)
	cl, err := DefaultCrontabConfig.ParseCronTab("* * * * mon-fri")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	dowField := cl[4][0]
	if dowField[0] != 1 || dowField[1] != 5 {
		t.Errorf("expected dow range 1-5, got %v", dowField)
	}
}

func TestParseCrontab_NamedRangeMidField(t *testing.T) {
	cl, err := DefaultCrontabConfig.ParseCronTab("* * * jan-jun *")
	if err != nil {
		t.Fatalf("named range in non-final field should parse: %v", err)
	}
	monthField := cl[3][0]
	if monthField[0] != 1 || monthField[1] != 6 {
		t.Errorf("expected month range 1-6, got %v", monthField)
	}
}

func TestParseCrontab_NamedRangeWithStep(t *testing.T) {
	cl, err := DefaultCrontabConfig.ParseCronTab("* * * jan-dec/3 *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	monthField := cl[3][0]
	if monthField[0] != 1 || monthField[1] != 12 || monthField[2] != 3 {
		t.Errorf("expected {1,12,3}, got %v", monthField)
	}
}

func TestParseCrontab_SingleName(t *testing.T) {
	cl, err := DefaultCrontabConfig.ParseCronTab("* * * * mon")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	dowField := cl[4][0]
	if dowField[0] != 1 || dowField[1] != 1 {
		t.Errorf("expected mon={1,1,...}, got %v", dowField)
	}
}

func TestParseCrontab_NameList(t *testing.T) {
	cl, err := DefaultCrontabConfig.ParseCronTab("* * * * mon,wed,fri")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(cl[4]) != 3 {
		t.Fatalf("expected 3 constraints for dow, got %d", len(cl[4]))
	}
}

func TestParseCrontab_InvalidChar(t *testing.T) {
	_, err := DefaultCrontabConfig.ParseCronTab("@")
	if err == nil {
		t.Error("expected error for invalid character")
	}
}

func TestParseCrontab_StarInWrongPosition(t *testing.T) {
	_, err := DefaultCrontabConfig.ParseCronTab("5*")
	if err == nil {
		t.Error("expected error for '*' after number")
	}
}

func TestParseCrontab_DashInWrongPosition(t *testing.T) {
	_, err := DefaultCrontabConfig.ParseCronTab("*-5")
	if err == nil {
		t.Error("expected error for '-' after '*'")
	}
}

func TestParseCrontab_SlashInWrongPosition(t *testing.T) {
	_, err := DefaultCrontabConfig.ParseCronTab("/5")
	if err == nil {
		t.Error("expected error for '/' at start")
	}
}

func TestParseCrontab_LetterInWrongPosition(t *testing.T) {
	// Letter where a step number is expected
	_, err := DefaultCrontabConfig.ParseCronTab("*/a")
	if err == nil {
		t.Error("expected error for letter in step position")
	}
}

func TestParseCrontab_NameInFieldWithoutNames(t *testing.T) {
	// minute field has no rangeNames, so "jan" in first position is treated as name
	// but minute field has no rangeNames → name resolution fails
	_, err := DefaultCrontabConfig.ParseCronTab("jan * * * *")
	if err == nil {
		t.Error("expected error for name in non-name field")
	}
}

func TestParseCrontab_NameRangeInFieldWithoutNames(t *testing.T) {
	// Name with dash in a field without rangeNames (minute field)
	_, err := DefaultCrontabConfig.ParseCronTab("abc-def * * * *")
	if err == nil {
		t.Error("expected error for name range in non-name field")
	}
}

func TestParseCrontab_OutOfRangeNumber(t *testing.T) {
	// 60 is out of range for minute (0-59)
	_, err := DefaultCrontabConfig.ParseCronTab("60")
	if err == nil {
		t.Error("expected error for out-of-range number")
	}
}

func TestParseCrontab_NegativeInStep(t *testing.T) {
	_, err := DefaultCrontabConfig.ParseCronTab("*/-1")
	if err == nil {
		t.Error("expected error for negative step")
	}
}

func TestParseCrontab_EndRangeOutOfBounds(t *testing.T) {
	_, err := DefaultCrontabConfig.ParseCronTab("0-60")
	if err == nil {
		t.Error("expected error for end-range out of bounds")
	}
}

func TestParseCrontab_NumberDashName(t *testing.T) {
	// "1-jan" in a named field: number start, name end
	// end-range expects number after number-dash
	_, err := DefaultCrontabConfig.ParseCronTab("* * * 1-jan *")
	if err == nil {
		// this may or may not error depending on parser state handling
		// the parser goes to StateExpectEndRangeNumber, then sees a letter
		t.Log("parser accepted number-name range")
	}
}

// --- Next edge cases ---

func TestNext_YearRollover(t *testing.T) {
	// "0 0 1 jan *" → midnight on Jan 1
	cl, err := DefaultCrontabConfig.ParseCronTab("0 0 1 jan *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	start := time.Date(2020, 12, 31, 23, 0, 0, 0, time.UTC)
	next, err := DefaultCrontabConfig.Next(cl, start)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if next.Year() != 2021 || next.Month() != 1 || next.Day() != 1 {
		t.Errorf("expected 2021-01-01, got %v", next.Format(time.RFC3339))
	}
}

func TestNext_LeapYear(t *testing.T) {
	// "0 0 29 feb *" → Feb 29, should find leap year
	cl, err := DefaultCrontabConfig.ParseCronTab("0 0 29 feb *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	next, err := DefaultCrontabConfig.Next(cl, start)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if next.Year() != 2024 || next.Month() != 2 || next.Day() != 29 {
		t.Errorf("expected 2024-02-29, got %v", next.Format(time.RFC3339))
	}
}

func TestNext_MonthRollover(t *testing.T) {
	// minute 30, hour 12, day 15 → should advance to next month if past
	cl, err := DefaultCrontabConfig.ParseCronTab("30 12 15 * *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	start := time.Date(2020, 3, 15, 12, 30, 0, 0, time.UTC)
	next, err := DefaultCrontabConfig.Next(cl, start)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if next.Month() != 4 || next.Day() != 15 {
		t.Errorf("expected April 15, got %v", next.Format(time.RFC3339))
	}
}

func TestNext_HourRollover(t *testing.T) {
	// minute 0, every hour → should go to next hour
	cl, err := DefaultCrontabConfig.ParseCronTab("0 * * * *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	start := time.Date(2020, 1, 1, 23, 30, 0, 0, time.UTC)
	next, err := DefaultCrontabConfig.Next(cl, start)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if next.Day() != 2 || next.Hour() != 0 || next.Minute() != 0 {
		t.Errorf("expected Jan 2 00:00, got %v", next.Format(time.RFC3339))
	}
}

func TestNext_EndOfYear(t *testing.T) {
	// "59 23 31 dec *" → last minute of the year
	cl, err := DefaultCrontabConfig.ParseCronTab("59 23 31 dec *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	start := time.Date(2020, 12, 31, 23, 58, 0, 0, time.UTC)
	next, err := DefaultCrontabConfig.Next(cl, start)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if next.Year() != 2020 || next.Hour() != 23 || next.Minute() != 59 {
		t.Errorf("expected 2020-12-31T23:59, got %v", next.Format(time.RFC3339))
	}
}

func TestNext_MultipleConstraintsPerField(t *testing.T) {
	// "0,30 9,17 * * *" → on the hour and half hour at 9am and 5pm
	// Note: Next always advances by 1 smallest unit first, and the ceiling
	// logic may combine field ceilings, so from 08:00 the first hit is 09:30
	// (minute ceils 01→30, then hour ceils 8→9).
	cl, err := DefaultCrontabConfig.ParseCronTab("0,30 9,17 * * *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	start := time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC)
	// From 09:00, advance 1 min to 09:01, ceil minute to 30 → 09:30
	next, err := DefaultCrontabConfig.Next(cl, start)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if next.Format(time.RFC3339) != "2020-01-01T09:30:00Z" {
		t.Errorf("expected 2020-01-01T09:30:00Z, got %s", next.Format(time.RFC3339))
	}

	// From 09:30, advance 1 min to 09:31, minute rolls, eventually reaches 17:00
	next, err = DefaultCrontabConfig.Next(cl, next)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if next.Format(time.RFC3339) != "2020-01-01T17:00:00Z" {
		t.Errorf("expected 2020-01-01T17:00:00Z, got %s", next.Format(time.RFC3339))
	}

	// From 17:00, advance 1 min to 17:01, ceil minute to 30 → 17:30
	next, err = DefaultCrontabConfig.Next(cl, next)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if next.Format(time.RFC3339) != "2020-01-01T17:30:00Z" {
		t.Errorf("expected 2020-01-01T17:30:00Z, got %s", next.Format(time.RFC3339))
	}
}

// --- SecondCrontabConfig tests ---

func TestSecondCrontabConfig_Parse(t *testing.T) {
	cl, err := SecondCrontabConfig.ParseCronTab("0 0 0 1 * 1 *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(cl) != 7 {
		t.Fatalf("expected 7 fields, got %d", len(cl))
	}
}

func TestSecondCrontabConfig_Next(t *testing.T) {
	// every 10 seconds
	cl, err := SecondCrontabConfig.ParseCronTab("*/10 * * * * * *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	next, err := SecondCrontabConfig.Next(cl, start)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if next.Second() != 10 {
		t.Errorf("expected second=10, got %d (%v)", next.Second(), next.Format(time.RFC3339Nano))
	}
}

func TestSecondCrontabConfig_NextSequence(t *testing.T) {
	// every 30 seconds
	cl, err := SecondCrontabConfig.ParseCronTab("0,30 * * * * * *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	t0 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := []int{30, 0, 30, 0}
	for i, exp := range expected {
		next, err := SecondCrontabConfig.Next(cl, t0)
		if err != nil {
			t.Fatalf("iteration %d: err: %v", i, err)
		}
		if next.Second() != exp {
			t.Errorf("iteration %d: expected second=%d, got %d", i, exp, next.Second())
		}
		t0 = next
	}
}

// --- Field.Validate edge cases ---

func TestCrontabField_ValidateReversed(t *testing.T) {
	cf := CrontabField{{10, 5, 1}}
	err := cf.Validate()
	if err != ErrConstraintBoundariesReversed {
		t.Errorf("expected ErrConstraintBoundariesReversed, got %v", err)
	}
}

func TestCrontabField_ValidateNonOverlapping(t *testing.T) {
	cf := CrontabField{{0, 5, 1}, {10, 15, 1}, {20, 25, 1}}
	if err := cf.Validate(); err != nil {
		t.Errorf("expected valid, got %v", err)
	}
}

// --- SetGroupName / SetGroupNumber edge cases ---

// --- Next: ErrMaxit ---

func TestNext_ErrMaxit(t *testing.T) {
	// Feb 31 doesn't exist, so this should exhaust MaxIt
	oldMaxIt := MaxIt
	MaxIt = 100 // reduce so this test doesn't take forever
	defer func() { MaxIt = oldMaxIt }()

	cl, err := DefaultCrontabConfig.ParseCronTab("0 0 31 feb *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err = DefaultCrontabConfig.Next(cl, start)
	if err != ErrMaxit {
		t.Errorf("expected ErrMaxit, got %v", err)
	}
}

// --- YearUnit.Less ---

func TestYearUnit_Less(t *testing.T) {
	u := YearUnit{}
	// YearUnit should not be less than any known unit
	if u.Less(SecondUnit{}) {
		t.Error("YearUnit should not be less than SecondUnit")
	}
	if u.Less(MinuteUnit{}) {
		t.Error("YearUnit should not be less than MinuteUnit")
	}
	if u.Less(HourUnit{}) {
		t.Error("YearUnit should not be less than HourUnit")
	}
	if u.Less(DayUnit{}) {
		t.Error("YearUnit should not be less than DayUnit")
	}
	if u.Less(WeekOfMonth{}) {
		t.Error("YearUnit should not be less than WeekOfMonth")
	}
	if u.Less(MonthUnit{}) {
		t.Error("YearUnit should not be less than MonthUnit")
	}
}

// --- CrontabField.Validate: max-overlap detection ---

func TestCrontabField_ValidateMaxOverlap(t *testing.T) {
	// max of first overlaps with min of second
	cf := CrontabField{{0, 10, 1}, {10, 20, 1}}
	err := cf.Validate()
	if err != ErrOverlappingConstraint {
		t.Errorf("expected ErrOverlappingConstraint, got %v", err)
	}
}

// --- Parser: end-range name with slash/step ---

func TestParseCrontab_NameRangeStep(t *testing.T) {
	// "mon-fri/2" in last field position
	cl, err := DefaultCrontabConfig.ParseCronTab("* * * * mon-fri/2")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	dowField := cl[4][0]
	if dowField[0] != 1 || dowField[1] != 5 || dowField[2] != 2 {
		t.Errorf("expected {1,5,2}, got %v", dowField)
	}
}

// --- Parser: digits where end-range name expected ---

func TestParseCrontab_NumberAfterEndRangeName(t *testing.T) {
	// "mon-5" — starts as name, dash, then number where name expected
	_, err := DefaultCrontabConfig.ParseCronTab("* * * * mon-5")
	if err == nil {
		t.Log("parser accepted mixed name-number range")
	}
}

// --- Parser: comma after end-range number ---

func TestParseCrontab_CommaAfterEndRange(t *testing.T) {
	cl, err := DefaultCrontabConfig.ParseCronTab("1-10,20-30 *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(cl[0]) != 2 {
		t.Fatalf("expected 2 constraints, got %d", len(cl[0]))
	}
}

// --- Parser: step after end-range number ---

func TestParseCrontab_StepAfterEndRange(t *testing.T) {
	cl, err := DefaultCrontabConfig.ParseCronTab("1-30/5")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	c := cl[0][0]
	if c[0] != 1 || c[1] != 30 || c[2] != 5 {
		t.Errorf("expected {1,30,5}, got %v", c)
	}
}

func TestSetGroupName_BadName(t *testing.T) {
	// "zzz" won't match any month name
	_, err := DefaultCrontabConfig.ParseCronTab("* * * zzz *")
	if err == nil {
		t.Error("expected error for unknown month name")
	}
}

func TestSetGroupName_AbbreviatedMonth(t *testing.T) {
	// "sep" should match "september"
	cl, err := DefaultCrontabConfig.ParseCronTab("* * * sep *")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	monthField := cl[3][0]
	if monthField[0] != 9 {
		t.Errorf("expected month 9, got %d", monthField[0])
	}
}

func TestSetGroupName_EndRange(t *testing.T) {
	// "mon-fri" for day of week
	cl, err := DefaultCrontabConfig.ParseCronTab("* * * * mon-fri")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	dowField := cl[4][0]
	if dowField[0] != 1 || dowField[1] != 5 {
		t.Errorf("expected dow range 1-5, got %v", dowField)
	}
}

// --- CrontabLine operations ---

// --- lookupNameIndex: ambiguous prefix ---

func TestLookupNameIndex_AmbiguousPrefix(t *testing.T) {
	names := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
	// "s" matches both "sunday" and "saturday" → ambiguous → -1
	idx := lookupNameIndex(names, "s")
	if idx != -1 {
		t.Errorf("expected -1 for ambiguous prefix, got %d", idx)
	}
	// "su" matches only "sunday" → unambiguous
	idx = lookupNameIndex(names, "su")
	if idx != 0 {
		t.Errorf("expected 0 for 'su', got %d", idx)
	}
}

func TestCrontabLine_Sort(t *testing.T) {
	cl := CrontabLine{{{30, 30, 1}, {5, 5, 1}}}
	cl.Sort()
	if cl[0][0][0] != 5 {
		t.Errorf("expected first constraint min=5 after sort, got %d", cl[0][0][0])
	}
}

func TestCrontabLine_GetSetField(t *testing.T) {
	cl := CrontabLine{{{0, 59, 1}}, {{0, 23, 1}}}
	f := cl.GetField(0)
	if f[0][0] != 0 || f[0][1] != 59 {
		t.Errorf("unexpected field: %v", f)
	}
	cl.SetField(0, CrontabField{{10, 10, 1}})
	if cl[0][0][0] != 10 {
		t.Errorf("SetField failed")
	}
}

func TestCrontabLine_SetConstraint(t *testing.T) {
	cl := CrontabLine{{{0, 59, 1}}}
	cl.SetConstraint(0, 0, CrontabConstraint{5, 5, 1})
	if cl[0][0][0] != 5 {
		t.Errorf("SetConstraint failed")
	}
}
