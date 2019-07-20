package validator

import (
	"github.com/eddieowens/axon"
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

func (v *validatorImpl) Struct(strct interface{}) error {
	return v.validator.Struct(strct)
}

func ValidatorFactory(_ axon.Injector, _ axon.Args) axon.Instance {
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
		return strings.Contains(fl.Param(), filepath.Ext(fl.Field().String()))
	})
	if err != nil {
		panic(err)
	}

	return axon.StructPtr(&validatorImpl{validator: v})
}
