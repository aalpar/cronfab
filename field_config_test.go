package cronfab

import (
	"fmt"
	"testing"
	"time"
)

func TestCeil(t *testing.T) {

	fc1 := &FieldConfig{
		Unit: MinuteUnit{},
		Name: "minute",
		Min:  0,
		Max:  59,
		GetIndex: func(t time.Time) int {
			return t.Minute()
		},
	}

	tcases := []struct {
		in         string
		expect     string
		roll       bool
		constraint [][3]int
	}{
		{
			in:         "2020-10-15T17:00:00Z",
			constraint: [][3]int{{3, 32, 5}},
			expect:     "2020-10-15T17:03:00Z",
			roll:       false,
		},
		{
			in:         "2020-10-15T17:07:00Z",
			constraint: [][3]int{{7, 32, 5}},
			expect:     "2020-10-15T17:07:00Z",
			roll:       false,
		},
		{
			in:         "2020-10-15T17:07:01Z",
			constraint: [][3]int{{7, 32, 1}},
			expect:     "2020-10-15T17:07:01Z",
			roll:       false,
		},
		{
			in:         "2020-10-15T17:32:00Z",
			constraint: [][3]int{{0, 30, 5}},
			expect:     "2020-10-15T17:32:00Z",
			roll:       true,
		},
		{
			in:         "2020-10-15T17:52:01Z",
			constraint: [][3]int{{5, 30, 5}},
			expect:     "2020-10-15T17:52:01Z",
			roll:       true,
		},
	}
	for i, tc := range tcases {
		ok := t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t0, err := time.Parse(time.RFC3339, tc.in)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			expect, err := time.Parse(time.RFC3339, tc.expect)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			t1, roll := fc1.Ceil(tc.constraint, t0)
			if !t1.Equal(expect) {
				t.Fatalf("missmatch: %q %q", t1.Format(time.RFC3339), expect.Format(time.RFC3339))
			}
			if tc.roll != roll {
				t.Fatalf("missmatch: %t %t", roll, tc.roll)
			}
		})
		if !ok {
			break
		}
	}

}
