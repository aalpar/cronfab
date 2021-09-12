package cronfab

import (
	"sort"
	"time"
)

// Unit units of time
type Unit interface {
	Less(Unit) bool
	String() string
	Add(t time.Time, n int) time.Time
	Trunc(t time.Time) time.Time
}

func SortUnits(us []Unit) {
	sort.Slice(us, func(i, j int) bool {
		return us[i].Less(us[j])
	})
}
