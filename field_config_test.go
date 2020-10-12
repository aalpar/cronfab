package cronfab

import (
	"testing"
	"time"
)

func TestCeil(t *testing.T) {
	fc1 := &FieldConfig{
		unit: MinuteUnit{},
		name: "minute",
		min:  0,
		max:  59,
		dateIndexFn: func(t time.Time) int {
			return t.Minute()
		},
	}

	dt := time.Date(2020, time.Month(10), 15, 17, 0, 0, 0, time.UTC)
	t0 := fc1.Ceil([][3]int{{1, 59, 5}}, dt)
	if !t0.Equal(time.Date(2020, time.Month(10), 15, 17, 1, 0, 0, time.UTC)) {
		t.Fatalf("missmatch: %v", t0)
	}

	dt = time.Date(2020, time.Month(10), 15, 17, 0, 0, 0, time.UTC)
	t0 = fc1.Ceil([][3]int{{59, 59, 1}}, dt)
	if !t0.Equal(time.Date(2020, time.Month(10), 15, 17, 59, 0, 0, time.UTC)) {
		t.Fatalf("missmatch: %v", t0)
	}

	dt = time.Date(2020, time.Month(10), 15, 17, 58, 0, 0, time.UTC)
	t0 = fc1.Ceil([][3]int{{0, 59, 5}}, dt)
	if !t0.Equal(time.Date(2020, time.Month(10), 15, 18, 0, 0, 0, time.UTC)) {
		t.Fatalf("missmatch: %v", t0)
	}

}

func TestIsSatisfactory(t *testing.T) {
	fc1 := &FieldConfig{
		unit: MinuteUnit{},
		name: "minute",
		min:  0,
		max:  59,
		dateIndexFn: func(t time.Time) int {
			return t.Minute()
		},
	}

	dt := time.Date(2020, time.Month(10), 15, 17, 4, 0, 0, time.UTC)
	y0 := fc1.IsSatisfactory([][3]int{{4, 59, 5}}, dt)
	if !y0 {
		t.Fatalf("missmatch: %v", y0)
	}

	dt = time.Date(2020, time.Month(10), 15, 17, 9, 0, 0, time.UTC)
	y0 = fc1.IsSatisfactory([][3]int{{4, 59, 1}}, dt)
	if !y0 {
		t.Fatalf("missmatch: %v", y0)
	}

	dt = time.Date(2020, time.Month(10), 15, 17, 5, 0, 0, time.UTC)
	y0 = fc1.IsSatisfactory([][3]int{{4, 59, 5}}, dt)
	if y0 {
		t.Fatalf("missmatch: %v", y0)
	}

	dt = time.Date(2020, time.Month(10), 15, 17, 0, 0, 0, time.UTC)
	y0 = fc1.IsSatisfactory([][3]int{{4, 59, 5}}, dt)
	if y0 {
		t.Fatalf("missmatch: %v", y0)
	}

}
