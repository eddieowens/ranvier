package semantics

import (
	"fmt"
	"github.com/eddieowens/ranvier/lang/domain"
	"github.com/eddieowens/ranvier/lang/semantics"
	"github.com/eddieowens/ranvier/lang/tests/unit"
	"github.com/stretchr/testify/suite"
	"path"
	"testing"
)

type ValidatorTest struct {
	unit.Unit
	analyzer semantics.Analyzer
}

func (v *ValidatorTest) SetupTest() {
	v.analyzer = v.Injector.Get(semantics.AnalyzerKey).GetStructPtr().(semantics.Analyzer)
}

func (v *ValidatorTest) TestValidExtends() {
	// -- Given
	//
	m := domain.ParsedSchema{
		Schema: domain.Schema{
			Name: "name",
			Extends: []string{
				path.Join(v.Resources(), "final.json"),
			},
			Config: []byte(""),
		},
		Dependencies: nil,
	}

	// -- When
	//
	err := v.analyzer.Semantics(&m)

	// -- Then
	//
	v.NoError(err)
}

func (v *ValidatorTest) TestInvalidExtendsExt() {
	// -- Given
	//
	m := domain.ParsedSchema{
		Schema: domain.Schema{
			Name: "name",
			Extends: []string{
				path.Join(v.Resources(), "final.jso"),
			},
			Config: []byte(""),
			Path:   "made-up.json",
		},
		Dependencies: nil,
	}

	expectedErrMsg := fmt.Sprintf("Failed to compile made-up.json due to field extends[0]: %s/final.jso does not "+
		"have a valid file extension. Valid extensions are json toml yaml yml.", v.Resources())

	// -- When
	//
	err := v.analyzer.Semantics(&m)

	// -- Then
	//
	v.EqualError(err, expectedErrMsg)
}

func (v *ValidatorTest) TestInvalidFilepath() {
	// -- Given
	//

	m := domain.ParsedSchema{
		Schema: domain.Schema{
			Name: "name",
			Extends: []string{
				path.Join(v.Resources(), "not-exist.yml"),
			},
			Config: []byte(""),
			Path:   "made-up.yml",
		},
		Dependencies: nil,
	}

	expectedErrMsg := fmt.Sprintf("Failed to compile made-up.yml due to field extends[0]: Could not find file %s/not-exist.yml.", v.Resources())

	// -- When
	//
	err := v.analyzer.Semantics(&m)

	// -- Then
	//
	v.EqualError(err, expectedErrMsg)
}

func TestValidator(t *testing.T) {
	suite.Run(t, new(ValidatorTest))
}
