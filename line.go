package cronfab

// CrontabLine represents the parsed user supplied contab specification
type CrontabLine [][][3]int

func (cl CrontabLine) String() string {
	q := ""
	for i := 0; i < len(cl); i++ {
		if i > 0 {
			q += " "
		}
		q += cl.GetField(i).String()
	}
	return q
}

// GetField return the crontab field at index i
func (cl CrontabLine) GetField(i int) CrontabField {
	return cl[i]
}

// SetField set the crontab field at index i
func (cl *CrontabLine) SetField(i int, f CrontabField) {
	(*cl)[i] = f
}

// SetConstraint set the field constraint for crontab field i, constraint j
func (cl *CrontabLine) SetConstraint(i int, j int, c CrontabConstraint) {
	(*cl)[i][j] = c
}
