package compiler

import (
	"github.com/eddieowens/ranvier/lang/compiler"
	"github.com/eddieowens/ranvier/lang/tests/unit"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CompilerTest struct {
	unit.Unit
}

func (c *CompilerTest) SetupTest() {
}

func (c *CompilerTest) TestToSchemaNameUpDirectory() {
	// -- Given
	//
	given := "../users/staging.json"
	expected := "users-staging"

	// -- When
	//
	actual := compiler.ToSchemaName(given)

	// -- Then
	//
	c.Equal(expected, actual)
}

func (c *CompilerTest) TestToSchemaNameDotInMiddle() {
	// -- Given
	//
	given := "../users/./staging.json"
	expected := "users-staging"

	// -- When
	//
	actual := compiler.ToSchemaName(given)

	// -- Then
	//
	c.Equal(expected, actual)
}

func (c *CompilerTest) TestToSchemaName() {
	// -- Given
	//
	given := "users/staging.json"
	expected := "users-staging"

	// -- When
	//
	actual := compiler.ToSchemaName(given)

	// -- Then
	//
	c.Equal(expected, actual)
}

func (c *CompilerTest) TestToSchemaNameLong() {
	// -- Given
	//
	given := "../one/two/three/four/five.json"
	expected := "one-two-three-four-five"

	// -- When
	//
	actual := compiler.ToSchemaName(given)

	// -- Then
	//
	c.Equal(expected, actual)
}

func (c *CompilerTest) TestToSchemaNameRoot() {
	// -- Given
	//
	given := "staging.json"
	expected := "staging"

	// -- When
	//
	actual := compiler.ToSchemaName(given)

	// -- Then
	//
	c.Equal(expected, actual)
}

func TestCompilerTest(t *testing.T) {
	suite.Run(t, new(CompilerTest))
}
