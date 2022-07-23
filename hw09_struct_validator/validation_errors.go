package hw09structvalidator

import "fmt"

type ValidationErrorLengthMismatch struct {
	expected int
	str      string
}

func (e *ValidationErrorLengthMismatch) Error() string {
	return fmt.Sprintf("expected length %d, actual length %d, source '%s'", e.expected, len(e.str), e.str)
}

type ValidationErrorRegexpDoesNotMatch struct {
	str    string
	regexp string
}

func (e *ValidationErrorRegexpDoesNotMatch) Error() string {
	return fmt.Sprintf("expected '%s' to match regexp '%s'", e.str, e.regexp)
}

type ValidationErrorInMismatch struct {
	str string
	in  string
}

func (e *ValidationErrorInMismatch) Error() string {
	return fmt.Sprintf("expected '%s' to be in '%s'", e.str, e.in)
}

type ValidationErrorMinValue struct {
	min   int
	value int
}

func (e *ValidationErrorMinValue) Error() string {
	return fmt.Sprintf("expected %d >= %d", e.value, e.min)
}

type ValidationErrorMaxValue struct {
	max   int
	value int
}

func (e *ValidationErrorMaxValue) Error() string {
	return fmt.Sprintf("expected %d <= %d", e.value, e.max)
}
