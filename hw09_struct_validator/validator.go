package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

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
			var validationErrors *ValidationErrors
			if errors.As(err, &validationErrors) {
				errs = append(errs, *validationErrors...)
				continue
			}
			var validationError *ValidationError
			if errors.As(err, &validationError) {
				errs = append(errs, *validationError)
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
