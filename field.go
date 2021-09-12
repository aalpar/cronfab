package cronfab

// CrontabField a crontab field which is an array of constraints
type CrontabField [][3]int

func (cf CrontabField) String() string {
	q := ""
	for i := 0; i < len(cf); i++ {
		if i > 0 {
			q += ","
		}
		q += cf.GetConstraint(i).String()
	}
	return q
}

// Len return the number of constraints
func (cf CrontabField) Len() int {
	return len(cf)
}

// GetConstraint return constraint at i
func (cf CrontabField) GetConstraint(i int) CrontabConstraint {
	return (cf)[i]
}

// SetConstraint set the constraint at i
func (cf *CrontabField) SetConstraint(i int, c CrontabConstraint) {
	(*cf)[i] = c
}

func (cf CrontabField) Ceil(x int) (int, bool) {
	q := 0
	roll := false
	for i := range cf {
		q, roll = cf.GetConstraint(i).Ceil(x)
		if !roll {
			return q, false
		}
	}
	return cf.GetConstraint(0).GetMin(), true
}
