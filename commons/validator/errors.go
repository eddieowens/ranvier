package validator

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

type errorMessageFactory func(e validator.FieldError) string

var errorMap = map[string]errorMessageFactory{
	"oneof": func(e validator.FieldError) string {
		return fmt.Sprintf("%s is invalid. Valid values are %s.", e.Value(), e.Param())
	},
	"ext": func(e validator.FieldError) string {
		return fmt.Sprintf("%s does not have a valid file extension. Valid extensions are %s.", e.Value(), e.Param())
	},
	"file": func(e validator.FieldError) string {
		return fmt.Sprintf("Could not find file %s.", e.Value())
	},
	"required": func(e validator.FieldError) string {
		return fmt.Sprintf("%s is required.", strings.ToLower(e.Field()))
	},
	"dns_1123": func(e validator.FieldError) string {
		return fmt.Sprintf("%s is invalid. The %s must start with an alphanumeric character, end with an alphanumeric "+
			"character and can only contain '-' special charcaters.", e.Value(), strings.ToLower(e.Field()))
	},
	"default": func(e validator.FieldError) string {
		return fmt.Sprintf("%s is an invalid %s", e.Value(), strings.ToLower(e.Field()))
	},
}

type ValidationError struct {
	errors       []error
	errorStrings []string
}

func (v *ValidationError) Error() string {
	return strings.Join(v.errorStrings, "\n")
}

func (v *ValidationError) Errors() []error {
	return v.errors
}

func (v *ValidationError) AddError(e error) {
	v.errors = append(v.errors, e)
	v.errorStrings = append(v.errorStrings, e.Error())
}

func newValidationError(errors ...error) error {
	errs := make([]error, len(errors))
	errsStr := make([]string, len(errors))

	for i, e := range errors {
		errs[i] = e
		errsStr[i] = e.Error()
	}

	return &ValidationError{
		errors:       errs,
		errorStrings: errsStr,
	}
}

func newValidationErrorFromValidator(errs validator.ValidationErrors) error {
	eSlice := make([]error, len(errs))
	for i, es := range errs {
		eSlice[i] = errors.New(validationErrorMsg(es))
	}
	return newValidationError(eSlice...)
}

func validationErrorMsg(f validator.FieldError) string {
	if v, ok := errorMap[f.Tag()]; ok {
		return v(f)
	} else {
		def := errorMap["default"]
		return def(f)
	}
}
