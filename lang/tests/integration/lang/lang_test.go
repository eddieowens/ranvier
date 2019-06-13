package lang

import (
	"github.com/eddieowens/ranvier/lang"
	"github.com/eddieowens/ranvier/lang/tests/integration"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RanvierTest struct {
	integration.Integration
}

func (r *RanvierTest) SetupTest() {
}

func (r *RanvierTest) TestNewCompiler() {
	// -- Given
	//

	// -- When
	//
	compiler := lang.NewCompiler()

	// -- Then
	//
	r.NotNil(compiler)
}

func TestRanvierTest(t *testing.T) {
	suite.Run(t, new(RanvierTest))
}
