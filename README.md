Cronfab
=======

Cronfab is a crontab time-and-date specification parser and processor with a configurable calendar.

Unlike [robfig/cron](https://github.com/robfig/cron), cronfab exposes an extensible field system: you can define custom calendar fields (e.g. week-of-month) with named ranges and arbitrary units.

All the unix standard features are supported:
- units may be specified by number or name
- lists and ranges are supported
- step values are supported

Cronfab does not support shell command execution, or specification nicknames (such as `@reboot`, `@annually`, `@yearly`, `@monthly`, `@weekly`, `@daily` or `@hourly`).

Parsers for classic 6-field (year, month, day of month, day of week, hour of day, minute or hour) and extended, 8 field, (6-field version extended to second of minute and week of month) are provided.  Other calendars and/or periods may be added.

Example
-------

```go
markers, err := cronfab.SecondCrontabConfig.ParseCronTab("*/5 * * * * * *")
if err != nil {
	log.Fatal(err)
}

next, err := cronfab.SecondCrontabConfig.Next(markers, time.Now())
if err != nil {
	log.Fatal(err)
}
fmt.Println(next)
```

Tests are the best source of additional usage examples.
