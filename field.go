package cronfab

// CrontabField a crontab field which is an array of constraints
type CrontabField [][3]int

func (cc CrontabField) String() string {
	q := ""
	for i := 0; i < len(cc); i++ {
		if i > 0 {
			q += ","
		}
		q += cc.GetConstraint(i).String()
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
