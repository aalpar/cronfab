package cronfab

import (
	"strconv"
	"strings"
	"time"
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

// Next return the next time after n as specified in the CrontabLine
func (cc *CrontabConfig) Next(ctl CrontabLine, n time.Time) (time.Time, error) {
	unitsRank := cc.Units
	u := unitsRank[0]
	n = u.Add(n, 1)
	var roll bool
	var newn time.Time
	var k int
	for k < len(unitsRank) {
		for k = 0; k < len(unitsRank); k++ {
			u = unitsRank[k]
			fieldsForUnit := cc.FieldUnits[u.String()]
			for _, i := range fieldsForUnit {
				newn, roll = cc.Fields[i].Ceil(ctl[i], n)
				if newn.After(n) || roll {
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

func lookupNameIndex(ss []string, s string) int {
	if len(s) < 3 {
		return -1
	}
	s = strings.ToLower(s)
	for i := range ss {
		found := true
		for j := 0; j < len(ss[i]) && j < len(s); j++ {
			if ss[i][j] != s[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}
