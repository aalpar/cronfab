package cronfab

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrMaxit = errors.New("maximum number of iterations met")
)

// CrontabConfig models the possible time specifications for a crontab entry
type CrontabConfig struct {
	Fields     []FieldConfig
	FieldUnits map[string][]int
	Units      []Unit
}

// NewCrontabConfig return a new crontab config for the supplied filed configs
func NewCrontabConfig(fields []FieldConfig) *CrontabConfig {
	q := &CrontabConfig{
		Fields:     fields,
		FieldUnits: map[string][]int{},
	}
	for i := 0; i < len(q.Fields); i++ {
		f := fields[i]
		u := f.unit
		q.Units = append(q.Units, u)
		SortUnits(q.Units)
		unm := u.String()
		q.FieldUnits[unm] = append(q.FieldUnits[unm], i)
	}
	// fields are sorted in unit order
	return q
}

const (
	MAXIT = 20000
)

// Next return the next time after n as specified in the CrontabLine
func (cc *CrontabConfig) Next(ctl CrontabLine, n time.Time) (time.Time, error) {
	unitsRank := cc.Units
	u := unitsRank[0]
	n = u.Add(n, 1)
	roll := false
	newn := time.Time{}
	k := 0
	j := 0
	for k < len(unitsRank) {
		for k = 0; k < len(unitsRank); k++ {
			u = unitsRank[k]
			fieldsForUnit := cc.FieldUnits[u.String()]
			for _, i := range fieldsForUnit {
				newn, roll = cc.Fields[i].Ceil(ctl[i], n)
				if !newn.Equal(n) || roll {
					break
				}
			}
			if roll {
				newn = unitsRank[k].Add(newn, 1)
				newn = unitsRank[k].Trunc(newn)
			}
			if !newn.Equal(n) || roll {
				break
			}
		}
		if j > MAXIT {
			return n, ErrMaxit
		}
		j++
		n = newn
	}
	return n, nil
}

// NameToNumber convert a constraint mnemonic to an index
func (cc *CrontabConfig) NameToNumber(i int, s string) int {
	return lookupNameIndex(cc.Fields[i].rangeNames, s)
}

// SetGroupName convert the group member name to an index for field fieldi
func (cc *CrontabConfig) SetGroupName(state State, i, fieldi int, a *CrontabConstraint, s string) error {
	k := cc.NameToNumber(fieldi, s) + cc.Fields[fieldi].min
	if k < cc.Fields[fieldi].min || k > cc.Fields[fieldi].max {
		return &ErrorBadIndex{FieldName: cc.Fields[fieldi].name, Value: k}
	}
	if state == StateInName {
		*a = [3]int{k, k, 1}
	} else if state == StateInEndRangeName {
		(*a)[1] = k
	} else {
		return &ErrorParse{Index: i, State: state}
	}
	return nil
}

// SetGroupNumber convert the group member number to an index for field fieldi
func (cc *CrontabConfig) SetGroupNumber(state State, i, fieldi int, a *CrontabConstraint, s string) error {
	k, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if k < cc.Fields[fieldi].min || k > cc.Fields[fieldi].max {
		return &ErrorBadIndex{FieldName: cc.Fields[fieldi].name, Value: k}
	}
	if state == StateInNumber {
		*a = [3]int{k, k, 1}
	} else if state == StateInEndRangeNumber {
		(*a)[1] = k
	} else if state == StateInStepNumber {
		if k < 1 {
			return &ErrorBadIndex{FieldName: cc.Fields[fieldi].name, Value: k}
		}
		(*a)[2] = k
	} else {
		return &ErrorParse{Index: i, State: state}
	}
	return nil
}

// Len crontab config is sortable
func (cc CrontabConfig) Len() int {
	return len(cc.Fields)
}

// lookupNameIndex lookup index of s in ss
func lookupNameIndex(ss []string, s string) int {
	q := -1
	// apples to apples
	s = strings.ToLower(s)
	for i := range ss {
		// go through all the names looking for a prefix match
		if s == ss[i] {
			// exact matches always result in an index
			return i
		} else if strings.HasPrefix(ss[i], s) {
			// unambiguous prefix matches result in an index
			if q >= 0 {
				return -1
			}
			q = i
		}
	}
	return q
}
