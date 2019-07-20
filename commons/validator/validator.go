package validator

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/commons"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var dns1123Regexp = regexp.MustCompile("[a-z0-9]([-a-z0-9]*[a-z0-9])?")

type Validator interface {
	Struct(strct interface{}) error
}

type validatorImpl struct {
	validator *validator.Validate
}

func (v *validatorImpl) Struct(s interface{}) error {
	err := v.validator.Struct(s)
	if err != nil {
		return newValidationErrorFromValidator(err.(validator.ValidationErrors))
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
