package cronfab

import (
	"unicode"
	"unicode/utf8"
)

type State int

const (
	StateExpectSplatOrNumberOrName = State(iota)
	StateInNumber                  // next must be '0'-'9' or ',' or ' ' or '-'
	StateInName                    // next must be 'a'-'z' or ',' or ' ' or '-'
	StateExpectEndRangeNumber      // next must be '0'-'9'
	StateInEndRangeNumber          // next must be '0'-'9' or '/' or ',' or ' '
	StateExpectEndRangeName        // next must be '0'-'9'
	StateInEndRangeName            // next must be 'a'-'z' or '/' or ',' or ' '
	StateExpectStep                // next must be '/' or ',' or ' '
	StateExpectStepNumber          // next must be '0'-'9'
	StateInStepNumber              // next must be '0'-'9' or ',' or ' '
)

// StateString returns a human readable string for the parser state.  Used in error messages
func StateString(x State) string {
	switch x {
	case StateExpectSplatOrNumberOrName:
		return "expecting '*' or number or name"
	case StateInNumber:
		return "in number"
	case StateInName:
		return "in name"
	case StateExpectEndRangeName:
		return "expecting end-range name"
	case StateInEndRangeName:
		return "in end-range name"
	case StateExpectEndRangeNumber:
		return "expecting end-range number"
	case StateInEndRangeNumber:
		return "in end-range number"
	case StateExpectStep:
		return "expecting '/'"
	case StateExpectStepNumber:
		return "expecting step number"
	case StateInStepNumber:
		return "in step number"
	}
	return "unknown"
}

// ParseCronTab parses a crontab string using the crontab configuration
func (cc *CrontabConfig) ParseCronTab(s string) (CrontabLine, error) {
	if len(s) == 0 {
		return CrontabLine{}, nil
	}
	i := 0
	j := 0
	fieldi := 0
	listi := 0
	ss := s
	numbers := CrontabConstraint([3]int{cc.Fields[fieldi].Min, cc.Fields[fieldi].Max, 1})
	markers := CrontabLine{{{}}}
	state := StateExpectSplatOrNumberOrName
	r, n := utf8.DecodeRuneInString(ss)
	for r != utf8.RuneError && len(ss) > 0 {
		if r >= '0' && r <= '9' {
			if state == StateInNumber {
			} else if state == StateInStepNumber {
			} else if state == StateInEndRangeNumber {
			} else if state == StateInName {
			} else if state == StateInEndRangeName {
			} else if state == StateExpectSplatOrNumberOrName {
				state = StateInNumber
				j = i
			} else if state == StateExpectEndRangeNumber {
				state = StateInEndRangeNumber
				j = i
			} else if state == StateExpectStepNumber {
				state = StateInStepNumber
				j = i
			} else {
				return CrontabLine{}, &ErrorParse{Index: i, State: state}
			}
		} else if unicode.IsLetter(r) {
			if state == StateInName {
			} else if state == StateInEndRangeName {
			} else if state == StateExpectSplatOrNumberOrName {
				state = StateInName
				j = i
			} else if state == StateExpectEndRangeName {
				state = StateInEndRangeName
				j = i
			} else {
				return CrontabLine{}, &ErrorParse{Index: i, State: state}
			}
		} else if r == '*' {
			if state == StateExpectSplatOrNumberOrName {
				numbers = CrontabConstraint([3]int{cc.Fields[fieldi].Min, cc.Fields[fieldi].Max, 1})
				state = StateExpectStep
			} else {
				return CrontabLine{}, &ErrorParse{Index: i, State: state}
			}
		} else if r == '-' {
			if state == StateInNumber {
				err := cc.SetGroupNumber(state, i, fieldi, &numbers, s[j:i])
				if err != nil {
					return CrontabLine{}, err
				}
				state = StateExpectEndRangeNumber
			} else if state == StateInName && len(cc.Fields[fieldi].RangeNames) > 0 {
				err := cc.SetGroupName(state, i, fieldi, &numbers, s[j:i])
				if err != nil {
					return CrontabLine{}, err
				}
				state = StateExpectEndRangeName
			} else {
				return CrontabLine{}, &ErrorParse{Index: i, State: state}
			}
		} else if r == '/' {
			if state == StateExpectStep {
				state = StateExpectStepNumber
			} else if state == StateInEndRangeNumber {
				err := cc.SetGroupNumber(state, i, fieldi, &numbers, s[j:i])
				if err != nil {
					return CrontabLine{}, err
				}
				state = StateExpectStepNumber
			} else if state == StateInEndRangeName && len(cc.Fields[fieldi].RangeNames) > 0 {
				err := cc.SetGroupName(state, i, fieldi, &numbers, s[j:i])
				if err != nil {
					return CrontabLine{}, err
				}
				state = StateExpectStepNumber
			} else {
				return CrontabLine{}, &ErrorParse{Index: i, State: state}
			}
		} else if r == ',' || r == ' ' {
			if state == StateExpectStep {
			} else if state == StateInEndRangeNumber || state == StateInStepNumber || state == StateInNumber {
				err := cc.SetGroupNumber(state, i, fieldi, &numbers, s[j:i])
				if err != nil {
					return CrontabLine{}, err
				}
			} else if (state == StateInEndRangeName || state == StateInName) && len(cc.Fields[fieldi].RangeNames) > 0 {
				err := cc.SetGroupName(state, i, fieldi, &numbers, s[j:i])
				if err != nil {
					return CrontabLine{}, err
				}
			} else {
				return CrontabLine{}, &ErrorParse{Index: i, State: state}
			}
			if r == ',' {
				markers.SetField(fieldi, append(markers[fieldi], CrontabConstraint([3]int{})))
				markers.SetConstraint(fieldi, listi, numbers)
				listi++
			} else if r == ' ' {
				markers = append(markers, CrontabField([][3]int{{}}))
				markers[fieldi][listi] = numbers
				fieldi++
				listi = 0
			}
			numbers = CrontabConstraint([3]int{cc.Fields[fieldi].Min, cc.Fields[fieldi].Max, 1})
			state = StateExpectSplatOrNumberOrName
		} else {
			return CrontabLine{}, &ErrorParse{Index: i, State: state}
		}
		i += n
		ss = ss[n:]
		r, n = utf8.DecodeRuneInString(ss)
	}
	if state == StateExpectStep {
	} else if state == StateInEndRangeNumber || state == StateInStepNumber || state == StateInNumber {
		err := cc.SetGroupNumber(state, i, fieldi, &numbers, s[j:i])
		if err != nil {
			return CrontabLine{}, err
		}
	} else if (state == StateInEndRangeName || state == StateInName) && len(cc.Fields[fieldi].RangeNames) > 0 {
		err := cc.SetGroupName(state, i, fieldi, &numbers, s[j:i])
		if err != nil {
			return CrontabLine{}, err
		}
	} else {
		return CrontabLine{}, &ErrorParse{Index: i, State: state}
	}
	markers[fieldi][listi] = numbers
	markers.Sort()
	return markers, nil
}
