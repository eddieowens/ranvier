package validator

import (
	"fmt"
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/commons"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var dns1123Regexp = regexp.MustCompile("[a-z0-9]([-a-z0-9]*[a-z0-9])?")

type ValidationErrorMutator func(*ValidationError)

type TagErrorFactory func(e validator.FieldError) string

type TagErrorMap map[string]TagErrorFactory

var baseTagErrorMap = TagErrorMap{
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

var tagMutatorMap = map[string]ValidationErrorMutator{}

type Validator interface {
	Struct(strct interface{}) error
	OnError(vem ValidationErrorMutator, tags ...string)
	TagErrors() TagErrorMap
}

type validatorImpl struct {
	validator *validator.Validate
}

func (v *validatorImpl) OnError(vem ValidationErrorMutator, tags ...string) {
	if tags == nil {
		tagMutatorMap["all"] = vem
	}
	for _, t := range tags {
		tagMutatorMap[t] = vem
	}
}

func (v *validatorImpl) TagErrors() TagErrorMap {
	return baseTagErrorMap
}

func (v *validatorImpl) Struct(s interface{}) error {
	err := v.validator.Struct(s)
	if err != nil {
		vErrs := newValidationErrorFromValidator(err.(validator.ValidationErrors))
		allMutator := tagMutatorMap["all"]
		for i, vErr := range vErrs {
			if mutator, ok := tagMutatorMap[vErr.OriginalError.Tag()]; ok {
				mutator(&vErrs[i])
			} else if allMutator != nil {
				allMutator(&vErrs[i])
			}
		}
		return vErrs
	}
	return nil
}

func Factory(_ axon.Injector, _ axon.Args) axon.Instance {
	v := validator.New()

	err := v.RegisterValidation("file", func(fl validator.FieldLevel) bool {
		_, err := os.Stat(fl.Field().String())
		return err == nil
	})
	if err != nil {
		panic(err)
	}

	err = v.RegisterValidation("dns_1123", func(fl validator.FieldLevel) bool {
		return dns1123Regexp.MatchString(fl.Field().String())
	})
	if err != nil {
		panic(err)
	}

	err = v.RegisterValidation("ext", func(fl validator.FieldLevel) bool {
		ext := filepath.Ext(fl.Field().String())
		if ext != "" {
			ext = ext[1:]
		}
		params := strings.Split(fl.Param(), " ")
		return commons.StringIncludes(ext, params)
	})
	if err != nil {
		panic(err)
	}

	return axon.StructPtr(&validatorImpl{validator: v})
}
