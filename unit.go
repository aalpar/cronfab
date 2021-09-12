package cronfab

import (
	"sort"
	"time"
)

// Unit units of time
type Unit interface {
	// Less return true if this Unit is a finer grain than the supplied Unit, x
	Less(x Unit) bool
	// String return a simple name for the unit.
	String() string
	// Add returns time, t, with n units added
	Add(t time.Time, n int) time.Time
	// Trunc returns a time with all the lower components zeroed.
	Trunc(t time.Time) time.Time
}

func SortUnits(us []Unit) {
	sort.Slice(us, func(i, j int) bool {
		return us[i].Less(us[j])
	})
}
