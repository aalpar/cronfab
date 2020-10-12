package cronfab

import (
	"fmt"
)

type ErrorBadIndex struct {
	FieldName string
	Value     int
}

func (i *ErrorBadIndex) Error() string {
	return fmt.Sprint("invalid index ", i.Value, " for ", i.FieldName)
}

type ErrorBadName struct {
	FieldName string
	Value     string
}

func (i *ErrorBadName) Error() string {
	return fmt.Sprint("invalid name \"", i.Value, "\" for ", i.FieldName)
}

type ErrorParse struct {
	Index int
	State State
}

func (q *ErrorParse) Error() string {
	return fmt.Sprintf("%s at %d", StateString(q.State), q.Index)
}
