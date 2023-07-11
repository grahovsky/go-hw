package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "12345678-1234-1234-1234-123456789abc",
				Name:   "John Doe",
				Age:    25,
				Email:  "some@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID:     "12345678-1234-1234-1234-123456789abc",
				Name:   "John Doe",
				Age:    17,
				Email:  "johndoe@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: fmt.Errorf("value must be greater than or equal to 18")},
			},
		},
		{
			in: User{
				ID:     "12345678",
				Name:   "John Doe",
				Age:    30,
				Email:  "johndoe@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("string length must be 36")},
			},
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: Response{
				Code: 201,
				Body: "some body",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("value must be one of 200,404,500")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if !errorsMatch(err, tt.expectedErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErr)
			}
		})
	}
}

func errorsMatch(err1, err2 error) bool {
	if len(err1.(ValidationErrors)) != len(err2.(ValidationErrors)) {
		return false
	}

	for i := range err1.(ValidationErrors) {
		if err1.(ValidationErrors)[i].Field != err2.(ValidationErrors)[i].Field || err1.(ValidationErrors)[i].Err.Error() != err2.(ValidationErrors)[i].Err.Error() {
			return false
		}
	}

	return true
}
