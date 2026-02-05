Cronfab
=======

Cronfab is a crontab time-and-date specification parser and processor with a configurable calendar.

Unlike [robfig/cron](https://github.com/robfig/cron), cronfab exposes an extensible field system: you can define custom calendar fields (e.g. week-of-month, moon phase) with named ranges and arbitrary units — no forking required.

All the standard crontab features are supported:
- units may be specified by number or name
- lists and ranges are supported
- step values are supported

Cronfab does not support shell command execution or specification nicknames (`@reboot`, `@daily`, etc.).

Built-in Configs
----------------

- **`DefaultCrontabConfig`** — classic 5-field: minute, hour, day-of-month, month, day-of-week
- **`SecondCrontabConfig`** — 7-field: second, minute, hour, day-of-month, week-of-month, month, day-of-week

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

Custom Calendars
----------------

`FieldConfig` fields are exported, so any package can construct its own calendar by implementing the `Unit` interface and calling `NewCrontabConfig`:

```go
config, err := cronfab.NewCrontabConfig([]cronfab.FieldConfig{
	{
		Unit:     cronfab.HourUnit{},
		Name:     "hour",
		Min:      0,
		Max:      23,
		GetIndex: func(t time.Time) int { return t.Hour() },
	},
	{
		Unit:       MoonPhaseUnit{},
		Name:       "moon phase",
		RangeNames: []string{"new", "waxingcrescent", "firstquarter", "waxinggibbous", "full", "waninggibbous", "thirdquarter", "waningcrescent"},
		Min:        0,
		Max:        7,
		GetIndex:   moonPhaseIndex,
	},
})
```

A full working example is in [examples/lunar](examples/lunar).
