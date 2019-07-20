package validator

import (
	"errors"
	"fmt"
	"github.com/eddieowens/ranvier/lang/domain"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

const errorTemplate = "Failed to compile %s due to field %s: %s"

type errorMessageFactory func(e validator.FieldError) string

var errorMap = map[string]errorMessageFactory{
	"oneof": func(e validator.FieldError) string {
		return fmt.Sprintf("%s is invalid. Valid values are %s.", e.Value(), strings.Join(strings.Split(e.Param(), ","), ", "))
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

func NewValidationError(errors ...error) error {
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

func NewSchemaValidationError(errs validator.ValidationErrors, file string) error {
	eSlice := make([]error, len(errs))
	for i, es := range errs {
		eSlice[i] = errors.New(validationErrorMsg(es, file))
	}
	return NewValidationError(eSlice...)
}

func validationErrorMsg(f validator.FieldError, file string) string {
	field := strings.ToLower(f.Field())
	if v, ok := errorMap[f.Tag()]; ok {
		return fmt.Sprintf(errorTemplate, file, field, v(f))
	} else {
		def := errorMap["default"]
		return fmt.Sprintf(errorTemplate, file, field, def(f))
	}
}
