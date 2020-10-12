package cronfab

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

// CrontabConfig models the possible time specifications for a crontab entry
type CrontabConfig struct {
	Fields    []FieldConfig
	fieldRank []int
}

// NewCrontabConfig return a new crontab config for the supplied filed configs
func NewCrontabConfig(fields []FieldConfig) *CrontabConfig {
	q := &CrontabConfig{
		Fields:    fields,
		fieldRank: []int{},
	}
	for i := 0; i < len(q.Fields); i++ {
		q.fieldRank = append(q.fieldRank, i)
	}
	// fields are sorted in unit order
	sort.Sort(q)
	return q
}

// Next return the next time
func (cc *CrontabConfig) Next(ctd CrontabLine, t time.Time) (time.Time, error) {
	if len(cc.Fields) == 0 {
		return t, nil
	}

	// truncate time to the lowest unit
	t = cc.Fields[0].unit.Trunc(t)
	t = cc.Fields[0].unit.Add(t, 1)

	// go through all the fields in the specification (which should match the number of fields in the configuration)
	t1 := t
	for i := 0; i < len(ctd); i++ {
		tabField := cc.GetCrontabFieldAtRank(ctd, i)
		configField := cc.GetFieldAtRank(i)

		// go through the constrains of this field in the specification
		t2 := configField.Ceil(tabField, t1)
		// check that all the fields for the same oder are compatible.
		satisfactory := true
		for j := i - 1; j >= 0 && cc.GetFieldAtRank(j).unit == configField.unit; j-- {
			if !cc.GetFieldAtRank(j).IsSatisfactory(cc.GetCrontabFieldAtRank(ctd, j), t2) {
				satisfactory = false
			}
		}
		if satisfactory {
			t1 = t2
		}
	}
	t = t1

	return t, nil
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

func (cc CrontabConfig) GetFieldAtRank(i int) FieldConfig {
	return cc.Fields[cc.fieldRank[i]]
}

func (cc CrontabConfig) GetCrontabFieldAtRank(ctd CrontabLine, i int) CrontabField {
	return ctd[cc.fieldRank[i]]
}

// Less crontab config is sortable
func (cc CrontabConfig) Less(i, j int) bool {
	return cc.Fields[cc.fieldRank[i]].unit.Less(cc.Fields[cc.fieldRank[j]].unit)
}

// Swap crontab config is sortable
func (cc *CrontabConfig) Swap(i, j int) {
	cc.fieldRank[i], cc.fieldRank[j] = cc.fieldRank[j], cc.fieldRank[i]
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
