package cronfab

import (
	"time"
)

// FieldConfig models the forms of time component ranges
type FieldConfig struct {
	Unit       Unit
	Name       string
	RangeNames []string
	Min        int
	Max        int
	GetIndex   func(time.Time) int
}

// Ceil performs a ceiling function for a calendar unit within the constraints of the crontab
func (configField FieldConfig) Ceil(tabField CrontabField, t time.Time) (time.Time, bool) {
	x0 := configField.GetIndex(t)
	x1, roll := tabField.Ceil(x0)
	if roll {
		return t, true
	}
	return configField.Unit.Add(t, x1-x0), false
}
