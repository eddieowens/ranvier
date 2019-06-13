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
	c.comp = c.Injector.GetStructPtr(compiler.CompilerKey).(compiler.Compiler)
}

func (c *CompilerTest) TestCompileAll() {
	// -- Given
	//
	outputDir := path.Join(c.Resources(), "output")
	opts := &compiler.CompileAllOptions{
		CompileOptions: compiler.CompileOptions{
			DryRun:          false,
			OutputDirectory: outputDir,
		},
		Force: false,
	}

	defer os.RemoveAll(outputDir)

	fp := path.Join(c.Resources(), "integration-valid")

	expectedConfig := map[string]TestConfig{
		"ProdUsers": {
			Db: TestConfigDb{
				Url:  "username@pg.mycompany.com:5432",
				Pool: 20,
			},
		},
		"StagingUsers": {
			Db: TestConfigDb{
				Pool: 5,
				Url:  "username@pg.staging.mycompany.com:5432",
			},
		},
		"Users": {
			Db: TestConfigDb{
				Pool: 3,
			},
		},
	}

	expected := map[string]domain.Schema{
		"ProdUsers": {
			Name:    "ProdUsers",
			Extends: []string{path.Join(fp, "prod.json")},
			Path:    path.Join(fp, "users", "prod-users.json"),
			Type:    "json",
		},
		"StagingUsers": {
			Name:    "StagingUsers",
			Extends: []string{path.Join(fp, "/staging.json")},
			Path:    path.Join(fp, "users", "staging-users.json"),
			Type:    "yaml",
		},
		"Users": {
			Name: "Users",
			Path: path.Join(fp, "users", "users.json"),
			Type: "toml",
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
			c.EqualConfigFromFile(path.Join(outputDir, fmt.Sprintf("%s.%s", v.Name, v.Type)), expectedConfig[k])
		}
	}
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
