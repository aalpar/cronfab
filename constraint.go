package cronfab

import (
	"errors"
	"fmt"
)

var (
	ErrConstraintBoundariesReversed = errors.New("constraint boundaries reversed")
)

// CrontabConstraint a crontab constraint expressed as {min, max, step}.
type CrontabConstraint [3]int

func (cc CrontabConstraint) String() string {
	return fmt.Sprintf("%d-%d/%d", cc.GetMin(), cc.GetMax(), cc.GetStep())
}

// GetMin get the minimum value of the range
func (cc CrontabConstraint) GetMin() int {
	return cc[0]
}

// GetMax get the maximum value of the range
func (cc CrontabConstraint) GetMax() int {
	return cc[1]
}

// GetStep get the step value for the range.  1 if no step was specified.
func (cc CrontabConstraint) GetStep() int {
	return cc[2]
}

// Validate return an error if the constraint boundaries are invalid
func (cc CrontabConstraint) Validate() error {
	if cc[0] > cc[1] {
		return ErrConstraintBoundariesReversed
	}
	return nil
}

// Ceil return the current or next greater value that satisfies the constraint.
func (cc CrontabConstraint) Ceil(x int) (int, bool) {
	if x < cc.GetMin() {
		return cc.GetMin(), false
	}
	if x > cc.GetMax() {
		return cc.GetMin(), true
	}
	rem := (x - cc.GetMin()) % cc.GetStep()
	if x >= cc.GetMin() && x <= cc.GetMax() && rem == 0 {
		return x, false
	}
	x0 := x + (cc.GetStep() - rem)
	if x0 > cc.GetMax() {
		return cc.GetMin(), true
	} else {
		return x0, false
	}
}
