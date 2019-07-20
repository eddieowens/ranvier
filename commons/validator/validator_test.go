package validator

import (
	"fmt"
	"github.com/eddieowens/axon"
	"github.com/stretchr/testify/suite"
	"os"
	"path"
	"testing"
)

type ValidatorTest struct {
	suite.Suite
	validator Validator
}

func (v *ValidatorTest) SetupTest() {
	inj := axon.NewInjector(axon.NewBinder(
		new(Package),
	))
	v.validator = inj.GetStructPtr(Key).(Validator)
}

func (v *ValidatorTest) TestExtValid() {
	// -- Given
	//
	type s struct {
		Ext string `validate:"ext=json toml"`
	}

	fp := path.Join(os.TempDir(), "what.json")
	_, _ = os.Create(fp)
	defer os.Remove(fp)

	given := s{
		Ext: fp,
	}

	// -- When
	//
	err := v.validator.Struct(given)

	// -- Then
	//
	v.NoError(err)
}

func (v *ValidatorTest) TestExtInvalid() {
	// -- Given
	//
	type s struct {
		Ext string `validate:"ext=txt toml"`
	}

	fp := path.Join(os.TempDir(), "what.json")
	_, _ = os.Create(fp)
	defer os.Remove(fp)

	given := s{
		Ext: fp,
	}
	expected := fmt.Sprintf("%s does not have a valid file extension. Valid extensions are txt toml.", fp)

	// -- When
	//
	err := v.validator.Struct(given)

	// -- Then
	//
	v.EqualError(err, expected)
}

func (v *ValidatorTest) TestExtInvalidFileNotFound() {
	// -- Given
	//
	type s struct {
		Ext string `validate:"ext=txt toml"`
	}

	fp := path.Join(os.TempDir(), "what.json")

	given := s{
		Ext: fp,
	}
	expected := fmt.Sprintf("%s does not have a valid file extension. Valid extensions are txt toml.", fp)

	// -- When
	//
	err := v.validator.Struct(given)

	// -- Then
	//
	v.EqualError(err, expected)
}

func TestValidatorTest(t *testing.T) {
	suite.Run(t, new(ValidatorTest))
}
