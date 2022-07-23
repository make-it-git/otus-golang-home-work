package hw09structvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

func getParamRegexp(validationParams string, validationType ValidatorType) (*regexp.Regexp, error) {
	reg, err := regexp.Compile(validationParams)
	if err != nil {
		return nil, &ErrorUnsupportedValidatorParams{
			validator: validationType,
			params:    validationParams,
		}
	}
	return reg, nil
}

func getParamInt(validationParams string, validationType ValidatorType) (int, error) {
	cmp, err := strconv.Atoi(validationParams)
	if err != nil {
		return 0, &ErrorUnsupportedValidatorParams{
			validator: validationType,
			params:    validationParams,
		}
	}
	return cmp, nil
}

func getParamInStr(validationParams string, validationType ValidatorType) ([]string, error) {
	params := strings.Split(validationParams, ",")
	if len(params) == 0 {
		return nil, &ErrorUnsupportedValidatorParams{
			validator: validationType,
			params:    validationParams,
		}
	}
	return params, nil
}

func getParamInInt(validationParams string, validationType ValidatorType) (int, int, error) {
	params := strings.Split(validationParams, ",")
	if len(params) != 2 {
		return 0, 0, &ErrorUnsupportedValidatorParams{
			validator: validationType,
			params:    validationParams,
		}
	}
	min, err := getParamInt(params[0], validationType)
	if err != nil {
		return 0, 0, err
	}
	max, err := getParamInt(params[1], validationType)
	if err != nil {
		return 0, 0, err
	}
	return min, max, nil
}
