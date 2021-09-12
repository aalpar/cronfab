Cronfab
=======

Cronfab is a crontab time-and-date specification parser and processor with a configurable calendar.

All the unix standard features are supported:
- units may be specified by number or name
- lists and ranges are suppored
- step values are supported

Cronfab does not support shell command execution, or specification nicknames (such as `@reboot`, `@annually`, `@yearly`, `@monthly`, `@weekly`, `@daily` or `@hourly`).

Parsers for classic 6-field (year, month, day of month, day of week, hour of day, minute or hour) and extended, 8 field, (6-field version extended to second of minute and week of month) are provided.  Other calendars and/or periods may be added.

Examples and Tests are the best source of documentation.

Example
-------

The example below outputs the parsed crontab entry and then runs for 20s, producing output every 5s.  The example demonstrates generating timeseries for a parsed crontab.

```
// crontab event at every 5 second interval
func ExampleEveryFiveSeconds() {
	markers, err := cronfab.SecondContabConfig.ParseCronTab("*/5 * * * * * *")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("%v\n", markers)

	// run for 4 intervals
	for i := 0; i < 4; i++ {
		t1, err := cronfab.SecondContabConfig.Next(markers, time.Now())
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		dt := t1.Sub(time.Now())
		time.Sleep(dt)
		fmt.Fprintf(os.Stderr, "time: %v\n", time.Now())
	}

	// Output: 0-59/5 0-59/1 0-23/1 1-31/1 1-5/1 1-12/1 0-6/1
}
```
