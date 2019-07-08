package semantics

import (
	"errors"
	"fmt"
	"github.com/eddieowens/ranvier/lang/domain"
	"github.com/eddieowens/ranvier/lang/services"
	"github.com/eddieowens/ranvier/lang/tests/unit"
	"github.com/stretchr/testify/suite"
	"path"
	"testing"
)

type ValidatorTest struct {
	unit.Unit
	validator services.Validator
}

func (v *ValidatorTest) SetupTest() {
	v.validator = v.Injector.Get(services.ValidatorKey).GetStructPtr().(services.Validator)
}

func (v *ValidatorTest) TestValidExtends() {
	// -- Given
	//
	m := domain.Schema{
		Name: "name",
		Extends: []string{
			path.Join(v.Resources(), "final.json"),
		},
		Config: []byte(""),
		Type:   domain.Json,
	}

	// -- When
	//
	err := v.validator.Schema(&m)

	// -- Then
	//
	v.NoError(err)
}

func (v *ValidatorTest) TestInvalidExtendsExt() {
	// -- Given
	//
	m := domain.Schema{
		Name: "name",
		Extends: []string{
			path.Join(v.Resources(), "final.jso"),
		},
		Config: []byte(""),
		Path:   "made-up.json",
		Type:   domain.Json,
	}

	// -- When
	//
	err := v.validator.Schema(&m)

	// -- Then
	//
	expected := []error{
		errors.New(
			fmt.Sprintf("Failed to compile made-up.json due to field extends[0]: Extension %s/final.jso is invalid. Valid extensions are toml, json, yaml, yml.", v.Resources()),
		),
	}
	if v.Error(err) {
		errs := err.(*services.ValidationError)
		v.ElementsMatch(expected, errs.Errors())
	}
}

func (v *ValidatorTest) TestInvalidFilepath() {
	// -- Given
	//
	m := domain.Schema{
		Name: "name",
		Extends: []string{
			path.Join(v.Resources(), "not-exist.yml"),
		},
		Config: []byte(""),
		Path:   "made-up.yml",
		Type:   domain.Json,
	}

	// -- When
	//
	err := v.validator.Schema(&m)

	// -- Then
	//
	expected := []error{
		errors.New(
			fmt.Sprintf("Failed to compile made-up.yml due to field extends[0]: Could not find file %s/not-exist.yml.", v.Resources()),
		),
	}

	if v.Error(err) {
		errs := err.(*services.ValidationError)
		v.ElementsMatch(expected, errs.Errors())
	}
}

func TestValidator(t *testing.T) {
	suite.Run(t, new(ValidatorTest))
}
