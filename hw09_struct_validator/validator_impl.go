package hw09structvalidator

import (
	"reflect"
	"strings"
)

type ValidatorType string

const (
	Regexp ValidatorType = "regexp"
	Len    ValidatorType = "len"
	In     ValidatorType = "in"
	Min    ValidatorType = "min"
	Max    ValidatorType = "max"
)

func getType(validatorType string) (ValidatorType, bool) {
	if validatorType == string(Regexp) {
		return Regexp, true
	}
	if validatorType == string(Len) {
		return Len, true
	}
	if validatorType == string(In) {
		return In, true
	}
	if validatorType == string(Min) {
		return Min, true
	}
	if validatorType == string(Max) {
		return Max, true
	}
	return "", false
}

func getValidatorString(fName string, v string, t ValidatorType, params string) (*func() error, error) {
	switch t { //nolint:exhaustive
	case Regexp:
		reg, err := getParamRegexp(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			ok := reg.MatchString(v)
			if !ok {
				return &ValidationError{
					Field: fName,
					Err: &ValidationErrorRegexpDoesNotMatch{
						str:    v,
						regexp: params,
					},
				}
			}
			return nil
		}
		return &f, nil
	case Len:
		length, err := getParamInt(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			if len(v) != length {
				return &ValidationError{
					Field: fName,
					Err: &ValidationErrorLengthMismatch{
						expected: length,
						str:      v,
					},
				}
			}
			return nil
		}
		return &f, nil
	case In:
		inValues := strings.Split(params, ",")
		if len(inValues) == 0 {
			return nil, &ErrorUnsupportedValidatorParams{
				validator: t,
				params:    params,
			}
		}
		f := func() error {
			for _, pValue := range inValues {
				if pValue == v {
					return nil
				}
			}
			return &ValidationError{
				Field: fName,
				Err: &ValidationErrorInMismatch{
					str: v,
					in:  params,
				},
			}
		}
		return &f, nil
	default:
		return nil, &ErrorUnsupportedValidatorType{
			validator:    t,
			providedType: "string",
		}
	}
}

//nolint:gocognit
func getValidatorStrSlice(fName string, v []string, t ValidatorType, params string) (*func() error, error) {
	switch t { //nolint:exhaustive
	case Regexp:
		reg, err := getParamRegexp(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			errs := make(ValidationErrors, 0)
			for _, s := range v {
				ok := reg.MatchString(s)
				if !ok {
					errs = append(errs, ValidationError{
						Field: fName,
						Err: &ValidationErrorRegexpDoesNotMatch{
							str:    s,
							regexp: params,
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
	case Len:
		length, err := getParamInt(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			errs := make(ValidationErrors, 0)
			for _, s := range v {
				if len(s) != length {
					errs = append(errs, ValidationError{
						Field: fName,
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
	case In:
		prms, err := getParamInStr(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			errs := make(ValidationErrors, 0)
			for _, s := range v {
				found := false
				for _, pValue := range prms {
					if pValue == s {
						found = true
						break
					}
				}
				if !found {
					errs = append(errs, ValidationError{
						Field: fName,
						Err: &ValidationErrorInMismatch{
							str: s,
							in:  params,
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
	default:
		return nil, &ErrorUnsupportedValidatorType{
			validator:    t,
			providedType: "[]string",
		}
	}
}

//nolint:gocognit
func getValidatorIntSlice(fName string, v []int, t ValidatorType, params string) (*func() error, error) {
	switch t { //nolint:exhaustive
	case Min:
		min, err := getParamInt(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			errs := make(ValidationErrors, 0)
			for _, i := range v {
				if i < min {
					errs = append(errs, ValidationError{
						Field: fName,
						Err: &ValidationErrorMinValue{
							min:   min,
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
	case Max:
		max, err := getParamInt(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			errs := make(ValidationErrors, 0)
			for _, i := range v {
				if i > max {
					errs = append(errs, ValidationError{
						Field: fName,
						Err: &ValidationErrorMaxValue{
							max:   max,
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
	case In:
		min, max, err := getParamInInt(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			errs := make(ValidationErrors, 0)
			for _, i := range v {
				if i < min || i > max {
					errs = append(errs, ValidationError{
						Field: fName,
						Err: &ValidationErrorNumberInMismatch{
							value: i,
							min:   min,
							max:   max,
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
	default:
		return nil, &ErrorUnsupportedValidatorType{
			validator:    t,
			providedType: "[]int",
		}
	}
}

func getValidatorInt(fName string, v int, t ValidatorType, params string) (*func() error, error) {
	switch t { //nolint:exhaustive
	case Min:
		cmp, err := getParamInt(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			if v < cmp {
				return &ValidationError{
					Field: fName,
					Err: &ValidationErrorMinValue{
						min:   cmp,
						value: v,
					},
				}
			}
			return nil
		}
		return &f, nil
	case Max:
		cmp, err := getParamInt(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			if v > cmp {
				return &ValidationError{
					Field: fName,
					Err: &ValidationErrorMaxValue{
						max:   cmp,
						value: v,
					},
				}
			}
			return nil
		}
		return &f, nil
	case In:
		min, max, err := getParamInInt(params, t)
		if err != nil {
			return nil, err
		}
		f := func() error {
			if v < min || v > max {
				return &ValidationError{
					Field: fName,
					Err: &ValidationErrorNumberInMismatch{
						value: v,
						min:   min,
						max:   max,
					},
				}
			}
			return nil
		}
		return &f, nil
	default:
		return nil, &ErrorUnsupportedValidatorType{
			validator:    t,
			providedType: "int",
		}
	}
}

func getValidator(fName string, v reflect.Value, t ValidatorType, params string) (*func() error, error) {
	switch v.Kind() { //nolint:exhaustive
	case reflect.Int:
		vInt := v.Int()
		return getValidatorInt(fName, int(vInt), t, params)
	case reflect.String:
		vString := v.String()
		return getValidatorString(fName, vString, t, params)
	case reflect.Slice:
		vSliceInt, ok := v.Interface().([]int)
		if ok {
			return getValidatorIntSlice(fName, vSliceInt, t, params)
		}
		vSliceStr, ok := v.Interface().([]string)
		if ok {
			return getValidatorStrSlice(fName, vSliceStr, t, params)
		}
		return nil, ErrorTypeConvert
	}

	return nil, &ErrorUnsupportedValidatorType{t, v.Type().String()}
}

func runValidation(fieldName string, validator string, value reflect.Value) error {
	parts := strings.Split(validator, ":")
	if len(parts) != 2 {
		return &ErrorUnsupportedValidator{validator}
	}

	validationType, ok := getType(parts[0])
	if !ok {
		return &ErrorUnsupportedValidator{
			validator: parts[0],
		}
	}
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

	err = (*validatorFunc)()
	return err
}
