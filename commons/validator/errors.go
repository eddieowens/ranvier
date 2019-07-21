package validator

import (
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errMsgs := make([]string, len(v))
	for i, e := range v {
		errMsgs[i] = e.Error()
	}
	return strings.Join(errMsgs, "\n")
}

type ValidationError struct {
	OriginalError validator.FieldError
	Msg           string
}

func (v ValidationError) Error() string {
	return v.Msg
}

func newValidationErrorFromValidator(errs validator.ValidationErrors) ValidationErrors {
	ve := make(ValidationErrors, len(errs))
	for i, e := range errs {
		ve[i] = ValidationError{
			Msg:           validationErrorMsg(e),
			OriginalError: e,
		}
	}
	return ve
}

func validationErrorMsg(f validator.FieldError) string {
	if v, ok := baseTagErrorMap[f.Tag()]; ok {
		return v(f)
	} else {
		def := baseTagErrorMap["default"]
		return def(f)
	}
}
