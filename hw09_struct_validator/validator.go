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

var (
	ErrInvalidInputType = errors.New("input parameter type must be a struct")

	ErrInvalidInteger   = errors.New("invalid integer")
	ErrInvalidRegexp    = errors.New("invalig regexp")
	ErrInvalidString    = errors.New("invalid string")
	ErrInvalidFieldType = errors.New("invalid field type")

	ErrValidationMin = errors.New("value less than a required minimum")
	ErrValidationMax = errors.New("value more than a required maximum")
	ErrValidationIn  = errors.New("value not found in list of possible values")
	ErrValidationLen = errors.New("invalid value length")
)

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
		return ErrInvalidInputType
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
			errs := validateField(field, fieldValue, validator)
			for _, err := range errs {
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: field.Name,
						Err:   err,
					})
				}
			}
		}
	}

	return validationErrors
}

func validateField(field reflect.StructField, fieldValue reflect.Value, validator string) []error {
	validatorParts := strings.SplitN(validator, ":", 2)
	validatorType := validatorParts[0]
	validatorArgs := ""
	if len(validatorParts) > 1 {
		validatorArgs = validatorParts[1]
	}

	switch field.Type.Kind() { //nolint: exhaustive
	case reflect.String:
		return []error{validateStringField(fieldValue, validatorType, validatorArgs)}
	case reflect.Int, reflect.Int64:
		return []error{validateIntField(fieldValue, validatorType, validatorArgs)}
	case reflect.Slice:
		errs := []error{}
		var valFunc func(reflect.Value, string, string) error

		elemType := field.Type.Elem()
		switch elemType.Kind() {
		case reflect.Int, reflect.Int64:
			valFunc = validateIntField
		case reflect.String:
			valFunc = validateStringField
		default:
			valFunc = validateOther
		}

		for i := 0; i < fieldValue.Len(); i++ {
			elem := fieldValue.Index(i)
			errs = append(errs, valFunc(elem, validatorType, validatorArgs))
		}
		return errs
	default:
		return []error{ErrInvalidFieldType}
	}
}

func validateStringField(fieldValue reflect.Value, validatorType, validatorArgs string) error {
	switch validatorType {
	case "len":
		strLen, err := strconv.Atoi(validatorArgs)
		if err != nil {
			return err
		}
		if len(fieldValue.String()) != strLen {
			return errors.Join(ErrValidationLen, fmt.Errorf("length must be %d", strLen))
		}
		return nil
	case "regexp":
		r, err := regexp.Compile(validatorArgs)
		if err != nil {
			return err
		}
		if !r.MatchString(fieldValue.String()) {
			return errors.Join(ErrInvalidRegexp, fmt.Errorf("string does not match regex pattern"))
		}
	case "in":
		validValues := strings.Split(validatorArgs, ",")
		fieldValueStr := fieldValue.String()
		for _, v := range validValues {
			if v == fieldValueStr {
				return nil
			}
		}
		return errors.Join(ErrValidationIn, fmt.Errorf("string must be one of %s", validatorArgs))
	default:
		return ErrInvalidString
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
			return errors.Join(ErrValidationMin, fmt.Errorf("value must be greater than or equal to %d", minValue))
		}
	case "max":
		maxValue, err := strconv.Atoi(validatorArgs)
		if err != nil {
			return err
		}
		if fieldValue.Int() > int64(maxValue) {
			return errors.Join(ErrValidationMax, fmt.Errorf("value must be less than or equal to %d", maxValue))
		}
	case "in":
		validValues := strings.Split(validatorArgs, ",")
		fieldValueInt := fieldValue.Int()
		for _, v := range validValues {
			if intValue, err := strconv.ParseInt(v, 10, 64); err == nil && intValue == fieldValueInt {
				return nil
			}
		}
		return errors.Join(ErrValidationIn, fmt.Errorf("value must be one of %s", validatorArgs))
	default:
		return nil
	}
	return nil
}

func validateOther(fieldValue reflect.Value, validatorType, validatorArgs string) error {
	switch validatorType {
	case "len":
		sliceLen, err := strconv.Atoi(validatorArgs)
		if err != nil {
			return err
		}

		if len(fieldValue.Bytes()) != sliceLen {
			return errors.Join(ErrValidationLen, fmt.Errorf("length must be %d", sliceLen))
		}
	default:
	}
	return nil
}
