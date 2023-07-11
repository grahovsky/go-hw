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

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, ve := range v {
		sb.WriteString(fmt.Sprintf("%s: %s\n", ve.Field, ve.Err.Error()))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}

	var validationErrors ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		validatorTag := field.Tag.Get("validate")
		if validatorTag == "" {
			continue
		}

		validators := strings.Split(validatorTag, "|")

		for _, validator := range validators {
			err := validateField(fieldValue, validator)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		}
	}

	return validationErrors
}

func validateField(fieldValue reflect.Value, validator string) error {
	validatorParts := strings.SplitN(validator, ":", 2)
	validatorType := validatorParts[0]
	validatorArgs := ""
	if len(validatorParts) > 1 {
		validatorArgs = validatorParts[1]
	}

	switch fieldValue.Kind() {
	case reflect.String:
		return validateStringField(fieldValue, validatorType, validatorArgs)
	case reflect.Int, reflect.Int64:
		return validateIntField(fieldValue, validatorType, validatorArgs)
	case reflect.Slice:
		switch fieldValues := fieldValue.Interface().(type) {
		case []string:
			for _, fv := range fieldValues {
				if err := validateStringField(reflect.ValueOf(fv), validatorType, validatorArgs); err != nil {
					return err
				}
			}
		default:
			return validateSliceField(fieldValue, validatorType, validatorArgs)
		}
	}

	return nil
}

func validateStringField(fieldValue reflect.Value, validatorType, validatorArgs string) error {
	switch validatorType {
	case "len":
		strLen, err := strconv.Atoi(validatorArgs)
		if err != nil {
			return err
		}
		if len(fieldValue.String()) != strLen {
			return fmt.Errorf("string length must be %d", strLen)
		}
	case "regexp":
		r, err := regexp.Compile(validatorArgs)
		if err != nil {
			return err
		}
		if !r.MatchString(fieldValue.String()) {
			return fmt.Errorf("string does not match regex pattern")
		}
	case "in":
		validValues := strings.Split(validatorArgs, ",")
		fieldValueStr := fieldValue.String()
		for _, v := range validValues {
			if v == fieldValueStr {
				return nil
			}
		}
		return fmt.Errorf("string must be one of %s", validatorArgs)
	}
	return nil
}

func validateIntField(fieldValue reflect.Value, validatorType, validatorArgs string) error {
	switch validatorType {
	case "min":
		minValue, err := strconv.Atoi(validatorArgs)
		if err != nil {
			return err
		}
		if fieldValue.Int() < int64(minValue) {
			return fmt.Errorf("value must be greater than or equal to %d", minValue)
		}
	case "max":
		maxValue, err := strconv.Atoi(validatorArgs)
		if err != nil {
			return err
		}
		if fieldValue.Int() > int64(maxValue) {
			return fmt.Errorf("value must be less than or equal to %d", maxValue)
		}
	case "in":
		validValues := strings.Split(validatorArgs, ",")
		fieldValueInt := fieldValue.Int()
		for _, v := range validValues {
			if intValue, err := strconv.ParseInt(v, 10, 64); err == nil && intValue == fieldValueInt {
				return nil
			}
		}
		return fmt.Errorf("value must be one of %s", validatorArgs)
	}
	return nil
}

func validateSliceField(fieldValue reflect.Value, validatorType, validatorArgs string) error {
	switch validatorType {
	case "len":
		sliceLen, err := strconv.Atoi(validatorArgs)
		if err != nil {
			return err
		}

		if fieldValue.Len() != sliceLen {
			return fmt.Errorf("slice length must be %d", sliceLen)
		}
	}
	return nil
}
