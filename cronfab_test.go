package cronfab

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestParseCrontab(t *testing.T) {
	tcases := []struct {
		in  string
		out [][][3]int
		err error
	}{
		{
			in:  "*",
			out: [][][3]int{{{0, 59, 1}}},
		},
		{
			in:  "* *",
			out: [][][3]int{{{0, 59, 1}}, {{0, 23, 1}}},
		},
		{
			in:  "* 10",
			out: [][][3]int{{{0, 59, 1}}, {{10, 10, 1}}},
		},
		{
			in:  "10 *",
			out: [][][3]int{{{10, 10, 1}}, {{0, 23, 1}}},
		},
		{
			in:  "*/2 10",
			out: [][][3]int{{{0, 59, 2}}, {{10, 10, 1}}},
		},
		{
			in:  "*/2 10,11 4-8 *",
			out: [][][3]int{{{0, 59, 2}}, {{10, 10, 1}, {11, 11, 1}}, {{4, 8, 1}}, {{1, 12, 1}}},
		},
		{
			in:  "*/1 4-8/2,10-12/2 */4",
			out: [][][3]int{{{0, 59, 1}}, {{4, 8, 2}, {10, 12, 2}}, {{1, 31, 4}}},
		},
		{
			in:  "* * * jan-mar",
			out: [][][3]int{{{0, 59, 1}}, {{0, 23, 1}}, {{1, 31, 1}}, {{1, 3, 1}}},
		},
		{
			in:  "*/100",
			err: &ErrorBadIndex{"minute", 100},
		},
		{
			in:  "*/0",
			err: &ErrorBadIndex{"minute", 0},
		},
		{
			in:  "*/-0",
			err: &ErrorParse{2, StateExpectStepNumber},
		},
	}
	for i := range tcases {

		tcase := tcases[i]

		t.Run(fmt.Sprintf(tcase.in), func(t *testing.T) {
			markers, err := DefaultContabConfig.ParseCronTab(tcase.in)
			if err != nil && tcase.err != nil {
				if !reflect.DeepEqual(err, tcase.err) {
					t.Errorf("unexpected value: %v != %v", err, tcase.err)
				}
			} else if err != nil && tcase.err == nil {
				t.Errorf("err: %v", err)
			} else {
				if !reflect.DeepEqual(markers, CrontabLine(tcase.out)) {
					t.Errorf("unexpected value: %v != %v", markers, tcase.out)
				}
			}
		})

	}

}

func TestNext(t *testing.T) {

	tcases := []struct {
		in  string
		add time.Duration
		out time.Duration
		err error
	}{
		{
			in:  "*/5 * * * *",
			add: 1 * time.Minute,
			out: 5 * time.Minute,
		},
		{
			in:  "*/5 * * * *",
			add: 7 * time.Minute,
			out: 10 * time.Minute,
		},
		{
			in:  "5,20,25,40 * * * *",
			add: 1 * time.Minute,
			out: 5 * time.Minute,
		},
		{
			in:  "5,20,25,40 * * * *",
			add: 7 * time.Minute,
			out: 20 * time.Minute,
		},
		{
			in:  "5,20,25,40 * * * *",
			add: (1 * time.Minute) + (time.Hour * 24),
			out: (5 * time.Minute) + (time.Hour * 24),
		},
		{
			in:  "5,20,25,40 * * * *",
			add: (27 * time.Minute) + (time.Hour * 24),
			out: (40 * time.Minute) + (time.Hour * 24),
		},
		{
			in:  "5,20,25,40 2-10 * * *",
			add: 1 * time.Minute,
			out: 5*time.Minute + (time.Hour * 2),
		},
		{
			in:  "* * * * thur",
			add: 1 * time.Minute,
			out: 2*time.Minute + (4 * (time.Hour * 24)),
		},
		{
			in:  "* * * * sun",
			add: 1 * time.Minute,
			out: 2 * time.Minute,
		},
		{
			in:  "5 0 * * sun",
			add: (24 * time.Hour) + (7 * time.Minute),
			out: ((24 * time.Hour) * 7) + (5 * time.Minute),
		},
		{
			in:  "5 0 27 * *",
			add: 28 * 24 * time.Hour,
			out: ((31 + 27) * 24 * time.Hour) + (5 * time.Minute),
		},
		{
			in:  "5 0 27 * wed",
			add: 28 * 24 * time.Hour,
			out: ((31 + 27) * 24 * time.Hour) + (5 * time.Minute),
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.in, func(t *testing.T) {
			cf, err := DefaultContabConfig.ParseCronTab(tcase.in)
			if err != nil {
				t.Fatal("err: ", err)
			}

			t0 := time.Date(1, 1, 0, 0, 0, 0, 0, time.UTC)
			t0 = t0.Add(tcase.add)

			t.Logf("t0 = %v %s", t0, t0.Weekday())

			t1, err := DefaultContabConfig.Next(cf, t0)
			if err != nil {
				t.Fatal("err: ", err)
			}

			t.Logf("t1 = %v %s", t1, t1.Weekday())

			t2 := time.Date(1, 1, 0, 0, 0, 0, 0, time.UTC)
			t2 = t2.Add(tcase.out)

			t.Logf("t2 = %v %s", t2, t2.Weekday())

			if !t1.Equal(t2) {
				t.Errorf("unexpected value: y%d m%d d%d H%d M%d", t1.Year(), int(t1.Month()), t1.Day(), t1.Hour(), t1.Minute())
			}
		})
	}

}
