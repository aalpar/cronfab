package cronfab

import (
	"log"
	"time"
)

// FieldConfig models the forms of time component ranges
type FieldConfig struct {
	unit        Unit
	name        string
	rangeNames  []string
	min         int
	max         int
	dateIndexFn func(time.Time) int
}

func (configField *FieldConfig) NormalizeMax(tabConstraint CrontabConstraint, t time.Time) time.Time {
	// number is after the range of ctd... set cdt to beginning of next period and loop around
	return configField.unit.Add(t,
		(configField.max-configField.dateIndexFn(t))+(1-configField.min)+tabConstraint.GetMin())
}

func (configField *FieldConfig) NormalizeMin(tabConstraint CrontabConstraint, t time.Time) time.Time {
	// number is after the range of ctd... set cdt to beginning of next period and loop around
	return configField.unit.Add(t, tabConstraint.GetMin()-configField.dateIndexFn(t))
}

func (configField *FieldConfig) NormalizeInRange(tabConstraint CrontabConstraint, t time.Time) time.Time {
	// we have a value within range, now see if it matches with
	// the current step value
	y0 := (configField.dateIndexFn(t) - tabConstraint.GetMin()) % tabConstraint.GetStep()
	if y0 != 0 {
		y1 := tabConstraint.GetStep() - y0
		return configField.unit.Add(t, y1)
	}
	return t
}

func (configField FieldConfig) IsSatisfactory(tabField CrontabField, t time.Time) bool {
	// go through the constrains of this field in the specification
	satisfactory := false
	for j := 0; j < tabField.Len() && !satisfactory; j++ {
		constraint := tabField.GetConstraint(j)
		if configField.dateIndexFn(t) > constraint.GetMax() {
			log.Printf("beyond max")
			continue
		}
		if configField.dateIndexFn(t) < constraint.GetMin() {
			continue
		}
		if (configField.dateIndexFn(t)-constraint.GetMin())%constraint.GetStep() != 0 {
			continue
		}
		satisfactory = true
	}
	return satisfactory
}

// Ceil performs a ceiling function for a calendar unit within the constrains of the crontab
func (configField *FieldConfig) Ceil(tabField CrontabField, t time.Time) time.Time {
	var t1 = t
	// go through the constrains of this field in the specification
	for j := 0; j < tabField.Len(); j++ {

		var t2 time.Time
		t3 := t
		for t3 != t2 {
			t2 = t3
			if configField.dateIndexFn(t2) > tabField.GetConstraint(j).GetMax() {
				// number is after the range o ctd... set cdt to beginning of next period and loop around
				t3 = configField.NormalizeMax(tabField.GetConstraint(j), t2)
			} else if configField.dateIndexFn(t) < tabField.GetConstraint(j).GetMin() {
				// number is before the range o ctd... set cdt to beginning of this period and loop around
				t3 = configField.NormalizeMin(tabField.GetConstraint(j), t2)
			} else {
				// we have a value within range, now see if it matches with
				// the current step value
				t3 = configField.NormalizeInRange(tabField.GetConstraint(j), t2)
			}
		}

		// first constraint is a special case.  pick the satisfied constraint closest
		// to the current time
		if j == 0 || t3.Before(t1) {
			t1 = t3
		}

	}
	return t1
}
