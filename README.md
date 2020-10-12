Cronfab
-------

Cronfab is a crontab time-and-date specification parser and processor with a configurable calendar.

All the unix standard features are supported:
- units may be specified by number or name
- lists and ranges are suppored
- step values are supported

Cronfab does not support shell command execution, or specification nicknames (such as `@reboot`, `@annually`, `@yearly`, `@monthly`, `@weekly`, `@daily` or `@hourly`).

Parsers for classic 5-field (year, month, day of month, day of week, hour minute) and extended, 6 field, (year, month, day of month, day of week, hour minute, second) are provided.  Other calendars and/or periods may be added.

