package cronfab

import (
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
			in:  "*/2 11,10 4-8 *",
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

		t.Run(tcase.in, func(t *testing.T) {
			markers, err := DefaultCrontabConfig.ParseCronTab(tcase.in)
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

func TestNameIndex(t *testing.T) {
	tcases := []struct {
		in  string
		out int
		err error
	}{
		{
			in:  "*",
			out: -1,
		},
		{
			in:  "o",
			out: 1,
		},
		{
			in:  "z",
			out: 0,
		},
		{
			in:  "four",
			out: 4,
		},
		{
			in:  "fourteen",
			out: 9,
		},
		{
			in:  "fourth",
			out: -1,
		},
	}
	stringSet := []string{"zero", "one", "two", "three", "four", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen"}
	for i := range tcases {

		tcase := tcases[i]

		t.Run(tcase.in, func(t *testing.T) {
			out := lookupNameIndex(stringSet, tcase.in)
			if out != tcase.out {
				t.Errorf("unexpected value: %v != %v", out, tcase.out)
			}
		})

	}

}

func TestNext(t *testing.T) {

	tcases := []struct {
		in     string
		start  string
		expect string
		err    error
	}{
		{
			in:     "*/5 * * * *",
			start:  "0001-01-01T00:01:00Z",
			expect: "0001-01-01T00:05:00Z",
		},
		{
			in:     "*/5 * * * *",
			start:  "0001-01-01T00:05:00Z",
			expect: "0001-01-01T00:10:00Z",
		},
		{
			in:     "5,20,25,40 * * * *",
			start:  "0001-01-01T00:01:00Z",
			expect: "0001-01-01T00:05:00Z",
		},
		{
			in:     "5,20,25,40 * * * *",
			start:  "0001-01-01T00:07:00Z",
			expect: "0001-01-01T00:20:00Z",
		},
		{
			in:     "5,20,25,40 * * * *",
			start:  "0001-01-02T00:01:00Z",
			expect: "0001-01-02T00:05:00Z",
		},
		{
			in:     "5,20,25,40 * * * *",
			start:  "0001-01-02T00:27:00Z",
			expect: "0001-01-02T00:40:00Z",
		},
		{
			in:     "5,20,25,40 2-10 * * *",
			start:  "0001-01-01T00:01:00Z",
			expect: "0001-01-01T02:05:00Z",
		},
		{
			in:     "40,25,20,5 2-10 * * *",
			start:  "0001-01-01T00:01:00Z",
			expect: "0001-01-01T02:05:00Z",
		},
		{
			in:     "* * * * thur",
			start:  "0001-01-01T00:01:00Z",
			expect: "0001-01-04T00:02:00Z",
		},
		{
			in:     "* * * * sun",
			start:  "0001-01-01T00:01:00Z",
			expect: "0001-01-07T00:00:00Z",
		},
		{
			in:     "5 0 * * sun",
			start:  "0001-01-02T00:07:00Z",
			expect: "0001-01-07T00:05:00Z",
		},
		{
			in:     "5 0 27 * *",
			start:  "0001-01-29T00:00:00Z",
			expect: "0001-02-27T00:05:00Z",
		},
		{
			in:     "5 0 27 * wed",
			start:  "0001-01-29T00:00:00Z",
			expect: "0001-06-27T00:05:00Z",
		},
	}

	for _, tcase := range tcases {
		ok := t.Run(tcase.in, func(t *testing.T) {
			cf, err := DefaultCrontabConfig.ParseCronTab(tcase.in)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			// start with the first
			t0, err := time.Parse(time.RFC3339, tcase.start)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			t.Logf("t0 = %q %s", t0.Format(time.RFC3339), t0.Weekday())

			t1, err := DefaultCrontabConfig.Next(cf, t0)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			t.Logf("t0 next = %q %s", t1.Format(time.RFC3339), t1.Weekday())

			t2, err := time.Parse(time.RFC3339, tcase.expect)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			t.Logf("expected t2 = %q %s", t2.Format(time.RFC3339), t2.Weekday())

			if !t1.Equal(t2) {
				t.Fatalf("unexpected value: %q %s", t1.Format(time.RFC3339), t1.Weekday())
			}
		})
		if !ok {
			break
		}
	}

}

func TestNextes(t *testing.T) {

	tcases := []struct {
		in    string
		start string
		outs  map[int]string
		err   error
	}{
		{
			in:    "*/5 * * * *",
			start: "0001-01-01T00:01:00Z",
			outs: map[int]string{
				0: "0001-01-01T00:01:00Z",
				1: "0001-01-01T00:05:00Z",
				2: "0001-01-01T00:10:00Z",
			},
		},
		{
			in:    "*/5 * * * *",
			start: "0001-01-01T00:07:00Z",
			outs: map[int]string{
				0: "0001-01-01T00:07:00Z",
				1: "0001-01-01T00:10:00Z",
				2: "0001-01-01T00:15:00Z",
			},
		},
		{
			in:    "1-3/5 * * * *",
			start: "0001-01-01T00:07:00Z",
			outs: map[int]string{
				0: "0001-01-01T00:07:00Z",
				1: "0001-01-01T01:01:00Z",
				2: "0001-01-01T02:01:00Z",
				3: "0001-01-01T03:01:00Z",
			},
		},
		{
			in:    "5,20,25,40 * * * *",
			start: "0001-01-01T00:01:00Z",
			outs: map[int]string{
				0: "0001-01-01T00:01:00Z",
				1: "0001-01-01T00:05:00Z",
				2: "0001-01-01T00:20:00Z",
				3: "0001-01-01T00:25:00Z",
				4: "0001-01-01T00:40:00Z",
				5: "0001-01-01T01:05:00Z",
			},
		},
		{
			in:    "40,25,20,5 * * * *",
			start: "0001-01-01T00:07:00Z",
			outs: map[int]string{
				0: "0001-01-01T00:07:00Z",
				1: "0001-01-01T00:20:00Z",
				2: "0001-01-01T00:25:00Z",
			},
		},
		{
			in:    "5,20,25,40 * * * *",
			start: "0001-01-01T00:07:00Z",
			outs: map[int]string{
				0: "0001-01-01T00:07:00Z",
				1: "0001-01-01T00:20:00Z",
				2: "0001-01-01T00:25:00Z",
			},
		},
		{
			in:    "5,20,25,40 * * * *",
			start: "0001-01-02T00:01:00Z",
			outs: map[int]string{
				0: "0001-01-02T00:01:00Z",
				1: "0001-01-02T00:05:00Z",
				2: "0001-01-02T00:20:00Z",
			},
		},
		{
			in:    "5,20,25,40 * * * *",
			start: "0001-01-02T00:27:00Z",
			outs: map[int]string{
				0: "0001-01-02T00:27:00Z",
				1: "0001-01-02T00:40:00Z",
				2: "0001-01-02T01:05:00Z",
			},
		},
		{
			in:    "5,20,25,40 2-10 * * *",
			start: "0001-01-01T00:01:00Z",
			outs: map[int]string{
				0:  "0001-01-01T00:01:00Z",
				1:  "0001-01-01T02:05:00Z",
				2:  "0001-01-01T02:20:00Z",
				3:  "0001-01-01T02:25:00Z",
				4:  "0001-01-01T02:40:00Z",
				5:  "0001-01-01T03:05:00Z",
				6:  "0001-01-01T03:20:00Z",
				7:  "0001-01-01T03:25:00Z",
				8:  "0001-01-01T03:40:00Z",
				9:  "0001-01-01T04:05:00Z",
				10: "0001-01-01T04:20:00Z",
				11: "0001-01-01T04:25:00Z",
				12: "0001-01-01T04:40:00Z",
				13: "0001-01-01T05:05:00Z",
				14: "0001-01-01T05:20:00Z",
				15: "0001-01-01T05:25:00Z",
				16: "0001-01-01T05:40:00Z",
				17: "0001-01-01T06:05:00Z",
				18: "0001-01-01T06:20:00Z",
				19: "0001-01-01T06:25:00Z",
				20: "0001-01-01T06:40:00Z",
				21: "0001-01-01T07:05:00Z",
				22: "0001-01-01T07:20:00Z",
				23: "0001-01-01T07:25:00Z",
				24: "0001-01-01T07:40:00Z",
				25: "0001-01-01T08:05:00Z",
				26: "0001-01-01T08:20:00Z",
				27: "0001-01-01T08:25:00Z",
				28: "0001-01-01T08:40:00Z",
				29: "0001-01-01T09:05:00Z",
				30: "0001-01-01T09:20:00Z",
				31: "0001-01-01T09:25:00Z",
				32: "0001-01-01T09:40:00Z",
				33: "0001-01-01T10:05:00Z",
				34: "0001-01-01T10:20:00Z",
				35: "0001-01-01T10:25:00Z",
				36: "0001-01-01T10:40:00Z",
				37: "0001-01-02T02:05:00Z",
				38: "0001-01-02T02:20:00Z",
				39: "0001-01-02T02:25:00Z",
				40: "0001-01-02T02:40:00Z",
			},
		},
		{
			in:    "* * * * thur",
			start: "0001-01-01T00:01:00Z",
			outs: map[int]string{
				0: "0001-01-01T00:01:00Z",
				1: "0001-01-04T00:02:00Z",
				2: "0001-01-04T00:03:00Z",
			},
		},
		{
			in:    "* * * * sun",
			start: "0001-01-01T00:01:00Z",
			outs: map[int]string{
				0: "0001-01-01T00:01:00Z",
				1: "0001-01-07T00:00:00Z",
				2: "0001-01-07T00:01:00Z",
			},
		},
		{
			in:    "5 0 * * sun",
			start: "0001-01-02T00:07:00Z",
			outs: map[int]string{
				0: "0001-01-02T00:07:00Z",
				1: "0001-01-07T00:05:00Z",
				2: "0001-01-14T00:05:00Z",
			},
		},
		{
			in:    "5 0 27 * *",
			start: "0001-01-28T00:00:00Z",
			outs: map[int]string{
				0: "0001-01-28T00:00:00Z",
				1: "0001-02-27T00:05:00Z",
				2: "0001-03-27T00:05:00Z",
			},
		},
		{
			in:    "5 0 27 * wed",
			start: "0001-01-28T00:00:00Z",
			outs: map[int]string{
				0: "0001-01-28T00:00:00Z",
				1: "0001-06-27T00:05:00Z",
				2: "0002-02-27T00:05:00Z",
			},
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.in, func(t *testing.T) {
			cf, err := DefaultCrontabConfig.ParseCronTab(tcase.in)
			if err != nil {
				t.Fatal("err: ", err)
			}

			// start with the first
			t0, err := time.Parse(time.RFC3339, tcase.start)
			if err != nil {
				t.Fatal("err: ", err)
			}

			i := 0
			for {

				datetime, ok := tcase.outs[i]

				if ok {
					t.Logf("t0 = %q %s", t0.Format(time.RFC3339), t0.Weekday())
					t.Logf("datetime = %q", datetime)

					t2, err := time.Parse(time.RFC3339, datetime)
					if err != nil {
						t.Fatalf("err: %v", err)
					}

					if !t0.Equal(t2) {
						t.Logf("got t0 = %q %s", t0.Format(time.RFC3339), t0.Weekday())
						t.Logf("expecting t2 = %q %s", t2.Format(time.RFC3339), t2.Weekday())
						t.Errorf("unexpected value: [%d] y%d m%d d%d H%d M%d", i, t0.Year(), int(t0.Month()), t0.Day(), t0.Hour(), t0.Minute())
						break
					}

					delete(tcase.outs, i)
					if len(tcase.outs) == 0 {
						break
					}
				}

				t0, err = DefaultCrontabConfig.Next(cf, t0)
				if err != nil {
					t.Fatal("err: ", err)
				}
				i++
			}
		})
	}

}
