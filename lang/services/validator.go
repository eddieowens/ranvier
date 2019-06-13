package services

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/commons"
	"github.com/eddieowens/ranvier/lang/domain"
	"gopkg.in/go-playground/validator.v9"
)

const ValidatorKey = "Validator"

type Validator interface {
	Validate(manifest *domain.Schema) error
}

type validatorImpl struct {
	validator  *validator.Validate
	FileFilter FileFilter    `inject:"FileFilter"`
	Filer      commons.Filer `inject:"Filer"`
}

func (v *validatorImpl) Validate(manifest *domain.Schema) error {
	err := v.validator.Struct(*manifest)
	if err != nil {
		return NewSchemaValidationError(err.(validator.ValidationErrors), manifest.Path)
	}
	return err
}

func validatorFactory(inj axon.Injector, _ axon.Args) axon.Instance {
	v := validator.New()

	filer := inj.GetStructPtr(FilerKey).(commons.Filer)
	err := v.RegisterValidation("filepath", filepathValidator(filer))
	if err != nil {
		panic(err)
	}

	filter := inj.GetStructPtr(FileFilterKey).(FileFilter)
	err = v.RegisterValidation("ext", extValidator(filter))
	if err != nil {
		panic(err)
	}

	return axon.StructPtr(&validatorImpl{validator: v})
}

func filepathValidator(filer commons.Filer) validator.Func {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field()
		if fp, ok := val.Interface().(string); ok {
			return filer.Exists(fp)
		}

		return false
	}
}

func extValidator(filter FileFilter) validator.Func {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field()
		if fp, ok := val.Interface().(string); ok {
			if fp == "" {
				return true
			}

			return filter.Filter(fp)
		}
		return false
	}
}
