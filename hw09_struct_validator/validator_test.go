package hw09structvalidator

import (
	"encoding/json"
	"errors"
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

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	intSlice struct {
		rangesMinMax []int `validate:"min:10|max:20"`
		rangesIn     []int `validate:"in:256,1024"`
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
				{Field: "Age", Err: ErrValidationMin},
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
				{Field: "ID", Err: ErrValidationLen},
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
				{Field: "Code", Err: ErrValidationIn},
			},
		},
		{
			in: intSlice{
				rangesMinMax: []int{10, 11, 12},
				rangesIn:     []int{256, 1025},
			},
			expectedErr: ValidationErrors{
				{Field: "rangesIn", Err: ErrValidationIn},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			var vErr, exprErr ValidationErrors
			err := Validate(tt.in)

			if errors.As(err, &vErr) && errors.As(tt.expectedErr, &exprErr) {
				if !errorsMatch(vErr, exprErr) {
					t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErr)
				}
			}
		})
	}
}

func errorsMatch(err1, err2 ValidationErrors) bool {
	if len(err1) != len(err2) {
		return false
	}

	for i := range err1 {
		if !errors.Is(err1[i].Err, err2[i].Err) {
			return false
		}
	}

	return true
}
