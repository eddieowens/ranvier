package compiler

import (
	"fmt"
	"github.com/eddieowens/ranvier/lang/compiler"
	"github.com/eddieowens/ranvier/lang/domain"
	"github.com/eddieowens/ranvier/lang/tests/integration"
	json "github.com/json-iterator/go"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type CompilerTest struct {
	integration.Integration
	comp compiler.Compiler
}

type TestConfigDb struct {
	Pool int    `json:"pool,omitempty"`
	Url  string `json:"url,omitempty"`
}

type TestConfig struct {
	Db TestConfigDb `json:"db"`
}

func (c *CompilerTest) SetupTest() {
	c.comp = c.Injector.GetStructPtr(compiler.Key).(compiler.Compiler)
}

func (c *CompilerTest) TestCompileAll() {
	// -- Given
	//
	outputDir := path.Join(c.Resources(), "output")
	_ = os.MkdirAll(outputDir, os.ModePerm)
	defer os.RemoveAll(outputDir)
	fp := path.Join(c.Resources(), "integration-valid")
	opts := compiler.CompileAllOptions{
		CompileOptions: compiler.CompileOptions{
			OutputDirectory: outputDir,
			DryRun:          false,
		},
		Force: false,
	}

	expectedConfig := map[string]TestConfig{
		"users-prod": {
			Db: TestConfigDb{
				Url:  "username@pg.mycompany.com:5432",
				Pool: 20,
			},
		},
		"users-staging": {
			Db: TestConfigDb{
				Pool: 5,
				Url:  "username@pg.staging.mycompany.com:5432",
			},
		},
		"users-users": {
			Db: TestConfigDb{
				Pool: 3,
			},
		},
	}

	expected := map[string]domain.Schema{
		"users-prod": {
			Name:    "users-prod",
			Extends: []string{path.Join(fp, "prod.json")},
			Path:    path.Join(fp, "users", "prod.json"),
		},
		"users-staging": {
			Name:    "users-staging",
			Extends: []string{path.Join(fp, "staging.json")},
			Path:    path.Join(fp, "users", "staging.json"),
		},
		"users-users": {
			Name: "users-users",
			Path: path.Join(fp, "users", "users.json"),
		},
	}

	// -- When
	//
	p, err := c.comp.CompileAll(fp, opts)

	// -- Then
	//
	if c.NoError(err) {
		c.Equal(fp, p.Path())
		for k, v := range expected {
			config := p.Schemas()[k].Config
			actual := p.Schemas()[k]
			actual.Config = nil

			var actualConfig TestConfig
			_ = json.Unmarshal(config, &actualConfig)

			c.Equal(v, actual)
			c.Equal(expectedConfig[k], actualConfig)
			c.EqualConfigFromFile(path.Join(outputDir, fmt.Sprintf("%s.json", v.Name)), expectedConfig[k])
		}
	}
}

func (c *CompilerTest) TestParse() {
	// -- Given
	//
	fp := path.Join(c.Resources(), "integration-valid")

	d, err := json.Marshal(TestConfig{
		Db: TestConfigDb{
			Pool: 3,
		},
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	expected := domain.Schema{
		Name:   "users-users",
		Path:   path.Join(fp, "users", "users.json"),
		Config: d,
	}

	// -- When
	//
	actual, err := c.comp.Parse("users/users.json", compiler.ParseOptions{
		Root: fp,
	})

	// -- Then
	//
	if c.NoError(err) {
		c.EqualSchemas(expected, *actual)
	}
}

func (c *CompilerTest) EqualSchemas(expected, actual domain.Schema) bool {
	var expConfig, actConfig interface{}
	err := json.Unmarshal(expected.Config, &expConfig)
	if c.NoError(err) {
		err = json.Unmarshal(actual.Config, &actConfig)
		if c.NoError(err) {
			if c.Equal(expConfig, actConfig) {
				expected.Config = nil
				actual.Config = nil
				return c.Equal(expected, actual)
			}
		}
	}
	return false
}

func (c *CompilerTest) EqualConfigFromFile(file string, expected TestConfig) {
	d, err := ioutil.ReadFile(file)
	ext := path.Ext(file)[1:]
	if c.NoError(err) {
		var actual TestConfig
		switch ext {
		case domain.Toml:
			err = toml.Unmarshal(d, &actual)
		case domain.Yml, domain.Yaml:
			err = yaml.Unmarshal(d, &actual)
		default:
			err = json.Unmarshal(d, &actual)
		}

		if c.NoError(err) {
			c.Equal(expected, actual)
		}
	}
}

func TestCompiler(t *testing.T) {
	suite.Run(t, new(CompilerTest))
}
