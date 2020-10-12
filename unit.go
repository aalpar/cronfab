package cronfab

import "time"

// Unit units of time
type Unit interface {
	Trunc(t time.Time) time.Time
	Add(t time.Time, n int) time.Time
	Less(u Unit) bool
}
