package hw09structvalidator

import (
	"encoding/json"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	AppSlice struct {
		Versions []string `validate:"len:5"`
	}

	AppRegexp struct {
		Version string `validate:"regexp:^\\d+.\\d+$"`
	}

	AppRegexpSlice struct {
		Versions []string `validate:"regexp:^\\d+.\\d+$"`
	}

	AppIn struct {
		Version string `validate:"in:1.1,1.2,1.3"`
	}

	AppInNumber struct {
		Version int `validate:"in:10,20"`
	}

	AppInNumberSlice struct {
		Versions []int `validate:"in:10,20"`
	}

	AppInSlice struct {
		Versions []string `validate:"in:1.1,1.2,1.3"`
	}

	AppMinMax struct {
		Version int `validate:"min:10|max:20"`
	}

	AppMinMaxSlice struct {
		Versions []int `validate:"min:10|max:20"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	InvalidValidator struct {
		Field string `validate:"invalid"`
	}

	InvalidValidatorType struct {
		Field string `validate:"min:10"`
	}
)

func TestValidate(t *testing.T) { //nolint:funlen
	tests := []struct {
		msg         string
		in          interface{}
		expectedErr error
	}{
		{
			msg:         "Not struct",
			in:          "test string",
			expectedErr: ErrorNotStruct,
		},
		{
			msg: "Multi validators",
			in: User{
				ID:    "123",
				Age:   17,
				Email: "invalid",
				Role:  "admin",
				meta:  nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err: &ValidationErrorLengthMismatch{
						expected: 36,
						str:      "123",
					},
				},
				ValidationError{
					Field: "Age",
					Err: &ValidationErrorMinValue{
						min:   18,
						value: 17,
					},
				},
				ValidationError{
					Field: "Email",
					Err: &ValidationErrorRegexpDoesNotMatch{
						str:    "invalid",
						regexp: "^\\w+@\\w+\\.\\w+$",
					},
				},
			},
		},
		{
			msg:         "No validators",
			in:          Token{},
			expectedErr: nil,
		},
		{
			msg:         "len validator for string",
			in:          App{"1.2.3"},
			expectedErr: nil,
		},
		{
			msg: "len validator for string",
			in:  App{"1.2.31"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err: &ValidationErrorLengthMismatch{
						expected: 5,
						str:      "1.2.31",
					},
				},
			},
		},
		{
			msg:         "len validator for []string",
			in:          AppSlice{[]string{"1.2.3", "1.2.4"}},
			expectedErr: nil,
		},
		{
			msg: "len validator for []string",
			in:  AppSlice{[]string{"1.2.3", "1.2.31", "1.2.41"}},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorLengthMismatch{
						expected: 5,
						str:      "1.2.31",
					},
				},
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorLengthMismatch{
						expected: 5,
						str:      "1.2.41",
					},
				},
			},
		},
		{
			msg:         "Regexp validator for string",
			in:          AppRegexp{"1.2"},
			expectedErr: nil,
		},
		{
			msg: "Regexp validator for string",
			in:  AppRegexp{"1.2.3"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err: &ValidationErrorRegexpDoesNotMatch{
						str:    "1.2.3",
						regexp: "^\\d+.\\d+$",
					},
				},
			},
		},
		{
			msg:         "Regexp validator for []string",
			in:          AppRegexpSlice{Versions: []string{"1.2", "1.3"}},
			expectedErr: nil,
		},
		{
			msg: "Regexp validator for []string",
			in:  AppRegexpSlice{Versions: []string{"1.2", "1.2.3", "1.2.4"}},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorRegexpDoesNotMatch{
						str:    "1.2.3",
						regexp: "^\\d+.\\d+$",
					},
				},
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorRegexpDoesNotMatch{
						str:    "1.2.4",
						regexp: "^\\d+.\\d+$",
					},
				},
			},
		},
		{
			msg:         "In validator for string",
			in:          AppIn{"1.2"},
			expectedErr: nil,
		},
		{
			msg: "In validator for string",
			in:  AppIn{"1.5"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err: &ValidationErrorInMismatch{
						str: "1.5",
						in:  "1.1,1.2,1.3",
					},
				},
			},
		},
		{
			msg:         "In validator for []string",
			in:          AppInSlice{[]string{"1.1", "1.2"}},
			expectedErr: nil,
		},
		{
			msg: "In validator for []string",
			in:  AppInSlice{[]string{"1.1", "1.2", "1.5", "a.b.c"}},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorInMismatch{
						str: "1.5",
						in:  "1.1,1.2,1.3",
					},
				},
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorInMismatch{
						str: "a.b.c",
						in:  "1.1,1.2,1.3",
					},
				},
			},
		},
		{
			msg:         "In validator for int",
			in:          AppInNumber{15},
			expectedErr: nil,
		},
		{
			msg: "In validator for int",
			in:  AppInNumber{25},
			expectedErr: &ValidationErrors{
				ValidationError{
					Field: "Version",
					Err: &ValidationErrorNumberInMismatch{
						value: 25,
						min:   10,
						max:   20,
					},
				},
			},
		},
		{
			msg:         "In validator for []int",
			in:          AppInNumberSlice{[]int{15, 20}},
			expectedErr: nil,
		},
		{
			msg: "In validator for []int",
			in:  AppInNumberSlice{[]int{15, 20, 21, 14, -1}},
			expectedErr: &ValidationErrors{
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorNumberInMismatch{
						value: 21,
						min:   10,
						max:   20,
					},
				},
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorNumberInMismatch{
						value: -1,
						min:   10,
						max:   20,
					},
				},
			},
		},
		{
			msg:         "Min/max validator for int",
			in:          AppMinMax{15},
			expectedErr: nil,
		},
		{
			msg: "Min validator for int",
			in:  AppMinMax{9},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err: &ValidationErrorMinValue{
						min:   10,
						value: 9,
					},
				},
			},
		},
		{
			msg: "Min validator for []int",
			in:  AppMinMaxSlice{[]int{10, 15, 5, -1}},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorMinValue{
						min:   10,
						value: 5,
					},
				},
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorMinValue{
						min:   10,
						value: -1,
					},
				},
			},
		},
		{
			msg: "Max validator for int",
			in:  AppMinMax{21},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err: &ValidationErrorMaxValue{
						max:   20,
						value: 21,
					},
				},
			},
		},
		{
			msg: "Max validator for []int",
			in:  AppMinMaxSlice{[]int{10, 15, 13, 21, 22}},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorMaxValue{
						max:   20,
						value: 21,
					},
				},
				ValidationError{
					Field: "Versions",
					Err: &ValidationErrorMaxValue{
						max:   20,
						value: 22,
					},
				},
			},
		},
		{
			msg: "Unsupported validator params",
			in: Response{
				Code: 0,
				Body: "",
			},
			expectedErr: &ErrorUnsupportedValidatorParams{
				validator: "in",
				params:    "200,404,500",
			},
		},
		{
			msg:         "Unsuppoted validator",
			in:          InvalidValidator{},
			expectedErr: &ErrorUnsupportedValidator{"invalid"},
		},
		{
			msg: "Unsupported validator type",
			in:  InvalidValidatorType{},
			expectedErr: &ErrorUnsupportedValidatorType{
				validator:    "min",
				providedType: "string",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if err == nil && tt.expectedErr == nil {
				return
			}
			if err == nil && tt.expectedErr != nil {
				t.Errorf("Expected '%v', got '%v'", tt.expectedErr, err)
				return
			}
			if err != nil && tt.expectedErr == nil {
				t.Errorf("Expected '%v', got '%v'", tt.expectedErr, err)
				return
			}
			if err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected '%v', got '%v'", tt.expectedErr, err)
			}
		})
	}
}
