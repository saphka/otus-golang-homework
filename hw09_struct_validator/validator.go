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
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("validation errors, total %d:", len(v)))

	for _, err := range v {
		sb.WriteString(fmt.Sprintf(" field %s error: %v;", err.Field, err.Err))
	}

	return sb.String()
}

var (
	ErrNotStruct           = errors.New("value is not a structure")
	ErrBadValidateParam    = errors.New("bad validate tag parameter")
	ErrUnknownValidateType = errors.New("unknown validator type")
	ErrWrongArgumentType   = errors.New("cannot validate value of this type")
	ErrLenConstraint       = errors.New("string length must be equal to")
	ErrRegexpConstraint    = errors.New("string must match regexp")
	ErrInConstraint        = errors.New("value must be in")
	ErrMinConstraint       = errors.New("value must be not be less than")
	ErrMaxConstraint       = errors.New("value must be not be greater than")
)

func Validate(v interface{}) error {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	var result ValidationErrors

	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		fieldType := valueOf.Type().Field(i)
		if !fieldType.IsExported() {
			continue
		}
		valid, ok := fieldType.Tag.Lookup("validate")
		if !ok {
			continue
		}

		validator, err := getCompositeValidator(valid)
		if err != nil {
			return err
		}

		result = validateField(result, field, validator, fieldType.Name)
	}

	if len(result) > 0 {
		return result
	}
	return nil
}

func validateField(
	result ValidationErrors,
	field reflect.Value,
	validator compositeValidatorFunc,
	name string,
) ValidationErrors {
	switch field.Kind() { //nolint:exhaustive
	case reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			result = append(result, validateField(nil, field.Index(i), validator, fmt.Sprintf("%s[%d]", name, i))...)
		}
	case reflect.Int, reflect.String:
		err := validator(field)
		if len(err) > 0 {
			for _, singleErr := range err {
				result = append(result, ValidationError{Field: name, Err: singleErr})
			}
		}
	}
	return result
}

type compositeValidatorFunc func(v reflect.Value) []error

func getCompositeValidator(valid string) (compositeValidatorFunc, error) {
	split := strings.Split(valid, "|")
	validators := make([]validatorFunc, 0, len(split))
	for _, singleValid := range split {
		validator, err := getValidator(singleValid)
		if err != nil {
			return nil, err
		}
		validators = append(validators, validator)
	}
	return func(v reflect.Value) []error {
		result := make([]error, 0, len(validators))
		for _, validator := range validators {
			err := validator(v)
			if err != nil {
				result = append(result, err)
			}
		}
		if len(result) > 0 {
			return result
		}
		return nil
	}, nil
}

type validatorFunc func(v reflect.Value) error

func getValidator(valid string) (validatorFunc, error) {
	split := strings.SplitN(valid, ":", 2)
	if len(split) < 2 {
		return nil, ErrBadValidateParam
	}

	switch split[0] {
	case "len":
		return getLenValidator(split[1])
	case "regexp":
		return getRegexpValidator(split[1])
	case "in":
		return getInValidator(split[1])
	case "min":
		return getMinValidator(split[1])
	case "max":
		return getMaxValidator(split[1])
	default:
		return nil, ErrUnknownValidateType
	}
}

func getLenValidator(param string) (validatorFunc, error) {
	length, err := strconv.Atoi(param)
	if err != nil {
		return nil, ErrBadValidateParam
	}
	return func(v reflect.Value) error {
		if v.Kind() != reflect.String {
			return ErrWrongArgumentType
		}
		if len(v.String()) != length {
			return fmt.Errorf("%w %d", ErrLenConstraint, length)
		}
		return nil
	}, nil
}

func getRegexpValidator(param string) (validatorFunc, error) {
	reg, err := regexp.Compile(param)
	if err != nil {
		return nil, ErrBadValidateParam
	}
	return func(v reflect.Value) error {
		if v.Kind() != reflect.String {
			return ErrWrongArgumentType
		}
		if !reg.MatchString(v.String()) {
			return fmt.Errorf("%w %s", ErrRegexpConstraint, param)
		}
		return nil
	}, nil
}

func getInValidator(param string) (validatorFunc, error) {
	split := strings.Split(param, ",")
	check := make(map[string]struct{}, len(split))
	for _, val := range split {
		check[val] = struct{}{}
	}
	return func(v reflect.Value) error {
		var val string
		switch v.Kind() { //nolint:exhaustive
		case reflect.Int:
			val = strconv.FormatInt(v.Int(), 10)
		case reflect.String:
			val = v.String()
		default:
			return ErrWrongArgumentType
		}
		if _, ok := check[val]; !ok {
			return fmt.Errorf("%w (%s)", ErrInConstraint, param)
		}
		return nil
	}, nil
}

func getMinValidator(param string) (validatorFunc, error) {
	min, err := strconv.Atoi(param)
	if err != nil {
		return nil, ErrBadValidateParam
	}
	return func(v reflect.Value) error {
		if v.Kind() != reflect.Int {
			return ErrWrongArgumentType
		}
		if v.Int() < int64(min) {
			return fmt.Errorf("%w %d", ErrMinConstraint, min)
		}
		return nil
	}, nil
}

func getMaxValidator(param string) (validatorFunc, error) {
	max, err := strconv.Atoi(param)
	if err != nil {
		return nil, ErrBadValidateParam
	}
	return func(v reflect.Value) error {
		if v.Kind() != reflect.Int {
			return ErrWrongArgumentType
		}
		if v.Int() > int64(max) {
			return fmt.Errorf("%w %d", ErrMaxConstraint, max)
		}
		return nil
	}, nil
}
