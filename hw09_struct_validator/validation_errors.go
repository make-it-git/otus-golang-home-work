package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

var ErrorNotStruct = errors.New("provided value is not struct")

var ErrorTypeConvert = errors.New("unable to convert type")

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("Field: %s, error: %s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, len(v))
	for i, err := range v {
		errs[i] = err.Error()
	}
	return strings.Join(errs, ". ")
}

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

type ValidationErrorNumberInMismatch struct {
	value int
	min   int
	max   int
}

func (e *ValidationErrorNumberInMismatch) Error() string {
	return fmt.Sprintf("expected %d to be between %d and %d", e.value, e.min, e.max)
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

type ErrorUnsupportedValidator struct {
	validator string
}

func (e *ErrorUnsupportedValidator) Error() string {
	return "unsupported validator: " + e.validator
}

type ErrorUnsupportedValidatorParams struct {
	validator ValidatorType
	params    string
}

func (e *ErrorUnsupportedValidatorParams) Error() string {
	return "unsupported validator params: " + string(e.validator) + " " + e.params
}

type ErrorUnsupportedValidatorType struct {
	validator    ValidatorType
	providedType string
}

func (e *ErrorUnsupportedValidatorType) Error() string {
	return fmt.Sprintf("unsupported validator type %s, not defined for type %s",
		e.validator, e.providedType)
}
