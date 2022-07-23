package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

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

type ErrorUnsupportedValidator struct {
	validator string
}

func (e *ErrorUnsupportedValidator) Error() string {
	return "unsupported validator: " + e.validator
}

type ErrorUnsupportedValidatorParams struct {
	validator string
	params    string
}

func (e *ErrorUnsupportedValidatorParams) Error() string {
	return "unsupported validator params: " + e.validator + " " + e.params
}

type ErrorUnsupportedValidatorType struct {
	validator    string
	providedType string
	expectedType string
}

func (e *ErrorUnsupportedValidatorType) Error() string {
	return fmt.Sprintf("unsupported validator type %s, provided type %s, expected type %s", e.validator, e.providedType, e.expectedType)
}

var ErrorNotStruct = errors.New("provided value is not struct")

func Validate(v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return ErrorNotStruct
	}

	errs := make(ValidationErrors, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		value := reflect.ValueOf(v).FieldByName(field.Name)

		validators := strings.Split(tag, "|")
		for _, validator := range validators {
			err := runValidation(field.Name, validator, value)
			if err == nil {
				continue
			}
			if errors.Is(err, ValidationErrors{}) {
				errs := err.(ValidationErrors)
				for _, e := range errs {
					errs = append(errs, e)
				}
				continue
			}
			if errors.Is(err, ValidationError{}) {
				errs = append(errs, err.(ValidationError))
				continue
			}
			return err
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func getValidator(fieldName string, value reflect.Value, validationType string, validationParams string) (*func(params string, value reflect.Value) error, error) {
	switch validationType {
	case "regexp":
		reg, err := regexp.Compile(validationParams)
		if err != nil {
			return nil, &ErrorUnsupportedValidatorParams{
				validator: validationType,
				params:    validationParams,
			}
		}
		if value.Kind() == reflect.String {
			f := func(p string, v reflect.Value) error {
				s := v.String()
				ok := reg.MatchString(s)
				if !ok {
					return &ValidationError{
						Field: fieldName,
						Err: &ValidationErrorRegexpDoesNotMatch{
							str:    s,
							regexp: validationParams,
						},
					}
				}

				return nil
			}
			return &f, nil
		}
		if value.Type().String() == "[]string" {
			f := func(p string, v reflect.Value) error {
				slice, ok := value.Interface().([]string)
				if !ok {
					return &ValidationError{
						Field: fieldName,
						Err:   errors.New("unable to get slice of string"),
					}
				}

				errs := make(ValidationErrors, 0)

				for _, s := range slice {
					ok := reg.MatchString(s)
					if !ok {
						errs = append(errs, ValidationError{
							Field: fieldName,
							Err: &ValidationErrorRegexpDoesNotMatch{
								str:    s,
								regexp: validationParams,
							},
						})
					}
				}
				if len(errs) > 0 {
					return errs
				}
				return nil
			}
			return &f, nil
		}
		return nil, &ErrorUnsupportedValidatorType{
			validator:    validationType,
			providedType: value.Kind().String(),
			expectedType: reflect.String.String(),
		}
	case "len":
		length, err := strconv.Atoi(validationParams)
		if err != nil {
			return nil, &ErrorUnsupportedValidatorParams{
				validator: validationType,
				params:    validationParams,
			}
		}
		if value.Kind() == reflect.String {
			f := func(p string, v reflect.Value) error {
				s := v.String()
				if len(s) != length {
					return &ValidationError{
						Field: fieldName,
						Err: &ValidationErrorLengthMismatch{
							expected: length,
							str:      s,
						},
					}
				}
				return nil
			}
			return &f, nil
		}
		if value.Type().String() == "[]string" {
			f := func(p string, v reflect.Value) error {
				slice, ok := value.Interface().([]string)
				if !ok {
					return &ValidationError{
						Field: fieldName,
						Err:   errors.New("unable to get slice of string"),
					}
				}

				errs := make(ValidationErrors, 0)

				for _, s := range slice {
					if len(s) != length {
						errs = append(errs, ValidationError{
							Field: fieldName,
							Err: &ValidationErrorLengthMismatch{
								expected: length,
								str:      s,
							},
						})
					}
				}

				if len(errs) > 0 {
					return errs
				}

				return nil
			}
			return &f, nil
		}
		return nil, &ErrorUnsupportedValidatorType{
			validator:    validationType,
			providedType: value.Kind().String(),
			expectedType: reflect.String.String(),
		}
	case "in":
		if value.Kind() == reflect.String {
			f := func(p string, v reflect.Value) error {
				values := strings.Split(p, ",")
				if len(values) == 0 {
					return &ErrorUnsupportedValidatorParams{
						validator: validationType,
						params:    p,
					}
				}

				s := v.String()
				for _, pValue := range values {
					if pValue == s {
						return nil
					}
				}

				return &ValidationError{
					Field: fieldName,
					Err: &ValidationErrorInMismatch{
						str: s,
						in:  validationParams,
					},
				}
			}
			return &f, nil
		}
		if value.Type().String() == "[]string" {
			f := func(p string, v reflect.Value) error {
				values := strings.Split(p, ",")
				if len(values) == 0 {
					return &ErrorUnsupportedValidatorParams{
						validator: validationType,
						params:    p,
					}
				}

				slice, ok := value.Interface().([]string)
				if !ok {
					return &ValidationError{
						Field: fieldName,
						Err:   errors.New("unable to get slice of string"),
					}
				}
				errs := make(ValidationErrors, 0)
				for _, s := range slice {
					found := false
					for _, pValue := range values {
						if pValue == s {
							found = true
							break
						}
					}
					if !found {
						errs = append(errs, ValidationError{
							Field: fieldName,
							Err: &ValidationErrorInMismatch{
								str: s,
								in:  p,
							},
						})
					}
				}
				if len(errs) > 0 {
					return errs
				}
				return nil
			}
			return &f, nil
		}
		if value.Kind() == reflect.Int {
			f := func(p string, v reflect.Value) error {
				values := strings.Split(p, ",")
				if len(values) != 2 {
					return &ErrorUnsupportedValidatorParams{
						validator: validationType,
						params:    p,
					}
				}

				min, err := strconv.Atoi(values[0])
				if err != nil {
					return &ErrorUnsupportedValidatorParams{
						validator: validationType,
						params:    p,
					}
				}
				max, err := strconv.Atoi(values[1])
				if err != nil {
					return &ErrorUnsupportedValidatorParams{
						validator: validationType,
						params:    p,
					}
				}

				i := v.Int()
				if i < int64(min) {
					return &ValidationError{
						Field: fieldName,
						Err:   errors.New(fmt.Sprintf("%d is less than %d", i, min)),
					}
				}
				if i > int64(max) {
					return &ValidationError{
						Field: fieldName,
						Err:   errors.New(fmt.Sprintf("%d is more than %d", i, max)),
					}
				}
				return nil
			}
			return &f, nil
		}
		if value.Type().String() == "[]int" {
			f := func(p string, v reflect.Value) error {
				values := strings.Split(p, ",")
				if len(values) != 2 {
					return &ErrorUnsupportedValidatorParams{
						validator: validationType,
						params:    p,
					}
				}

				min, err := strconv.Atoi(values[0])
				if err != nil {
					return &ErrorUnsupportedValidatorParams{
						validator: validationType,
						params:    p,
					}
				}
				max, err := strconv.Atoi(values[1])
				if err != nil {
					return &ErrorUnsupportedValidatorParams{
						validator: validationType,
						params:    p,
					}
				}

				slice, ok := value.Interface().([]int)
				if !ok {
					return &ValidationError{
						Field: fieldName,
						Err:   errors.New("unable to get slice of int"),
					}
				}

				for k, i := range slice {
					if i < min {
						return &ValidationError{
							Field: fieldName,
							Err:   errors.New(fmt.Sprintf("%d is less than %d at index %d", i, min, k)),
						}
					}
					if i > max {
						return &ValidationError{
							Field: fieldName,
							Err:   errors.New(fmt.Sprintf("%d is more than %d at index %d", i, max, k)),
						}
					}
				}
				return nil
			}
			return &f, nil
		}
		return nil, &ErrorUnsupportedValidatorType{
			validator:    validationType,
			providedType: value.Kind().String(),
			expectedType: reflect.String.String(),
		}
	case "min":
		cmp, err := strconv.Atoi(validationParams)
		if err != nil {
			return nil, &ErrorUnsupportedValidatorParams{
				validator: validationType,
				params:    validationParams,
			}
		}
		if value.Kind() == reflect.Int {
			f := func(p string, v reflect.Value) error {
				i := v.Int()
				if i < int64(cmp) {
					return &ValidationError{
						Field: fieldName,
						Err: &ValidationErrorMinValue{
							min:   cmp,
							value: int(i),
						},
					}
				}
				return nil
			}
			return &f, nil
		}
		if value.Type().String() == "[]int" {
			f := func(p string, v reflect.Value) error {
				slice, ok := value.Interface().([]int)
				if !ok {
					return &ValidationError{
						Field: fieldName,
						Err:   errors.New("unable to get slice of int"),
					}
				}
				errs := make(ValidationErrors, 0)
				for _, i := range slice {
					if i < cmp {
						errs = append(errs, ValidationError{
							Field: fieldName,
							Err: &ValidationErrorMinValue{
								min:   cmp,
								value: i,
							},
						})
					}
				}
				if len(errs) > 0 {
					return errs
				}
				return nil
			}
			return &f, nil
		}
		return nil, &ErrorUnsupportedValidatorType{
			validator:    validationType,
			providedType: value.Kind().String(),
			expectedType: reflect.Int.String(),
		}
	case "max":
		cmp, err := strconv.Atoi(validationParams)
		if err != nil {
			return nil, &ErrorUnsupportedValidatorParams{
				validator: validationType,
				params:    validationParams,
			}
		}
		if value.Kind() == reflect.Int {
			f := func(p string, v reflect.Value) error {
				i := v.Int()
				if i > int64(cmp) {
					return &ValidationError{
						Field: fieldName,
						Err: &ValidationErrorMaxValue{
							max:   cmp,
							value: int(i),
						},
					}
				}
				return nil
			}
			return &f, nil
		}
		if value.Type().String() == "[]int" {
			f := func(p string, v reflect.Value) error {
				slice, ok := value.Interface().([]int)
				if !ok {
					return &ValidationError{
						Field: fieldName,
						Err:   errors.New("unable to get slice of int"),
					}
				}
				errs := make(ValidationErrors, 0)
				for _, i := range slice {
					if i > cmp {
						errs = append(errs, ValidationError{
							Field: fieldName,
							Err: &ValidationErrorMaxValue{
								max:   cmp,
								value: i,
							},
						})
					}
				}
				if len(errs) > 0 {
					return errs
				}
				return nil
			}
			return &f, nil
		}
		return nil, &ErrorUnsupportedValidatorType{
			validator:    validationType,
			providedType: value.Kind().String(),
			expectedType: reflect.Int.String(),
		}
	default:
		return nil, &ErrorUnsupportedValidator{validationType}
	}
}

func runValidation(fieldName string, validator string, value reflect.Value) error {
	parts := strings.Split(validator, ":")
	if len(parts) != 2 {
		return &ErrorUnsupportedValidator{validator}
	}

	validationType := parts[0]
	validationParams := parts[1]
	if validationParams == "" {
		return &ErrorUnsupportedValidatorParams{
			validator: validationType,
			params:    validationParams,
		}
	}

	validatorFunc, err := getValidator(fieldName, value, validationType, validationParams)
	if err != nil {
		return err
	}

	err = (*validatorFunc)(validationParams, value)
	return err
}
