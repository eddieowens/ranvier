package compiler

import (
	"github.com/eddieowens/ranvier/commons"
	"github.com/eddieowens/ranvier/commons/validator"
	"github.com/eddieowens/ranvier/lang/domain"
	"github.com/eddieowens/ranvier/lang/semantics"
	"github.com/eddieowens/ranvier/lang/services"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const Key = "Compiler"

type CompileAllOptions struct {
	CompileOptions
	Force bool
}

type CompileOptions struct {
	ParseOptions
	// The directory that the file will be output to. If the directory does not exist, it will be created.
	OutputDirectory string `validate:"required_without=DryRun"`
	DryRun          bool
}

type ParseOptions struct {
	// The root directory of the file to parse which determines the prepended name of the file. The file that is being
	// parsed must lie within this directory.
	Root string `validate:"required,file"`
}

type Compiler interface {
	Compile(filepath string, options CompileOptions) (*domain.Schema, error)
	CompileAll(path string, options CompileAllOptions) (SchemaPack, error)
	ValidateSemantics(manifest *domain.Schema) error
	Parse(fp string, options ParseOptions) (*domain.Schema, error)
}

type compilerImpl struct {
	JsonMerger       commons.JsonMerger     `inject:"JsonMerger"`
	Analyzer         semantics.Analyzer     `inject:"Analyzer"`
	FileCollector    services.FileCollector `inject:"FileCollector"`
	FileService      services.FileService   `inject:"FileService"`
	Packer           SchemaPacker           `inject:"SchemaPacker"`
	ValidatorService validator.Validator    `inject:"Validator"`
}

func (c *compilerImpl) CompileAll(path string, options CompileAllOptions) (SchemaPack, error) {
	options.CompileOptions.ParseOptions.Root = path
	if err := c.ValidatorService.Struct(options); err != nil {
		return nil, err
	}
	files := c.FileCollector.Collect(path)
	files = c.FileService.SubtractPaths(path, files)

	pErr := NewSchemaPackError()
	pack := NewSchemaPack(path)
	for _, f := range files {
		s, err := c.Compile(f, options.CompileOptions)
		if err != nil {
			pErr.AddError(err)
			if !options.Force {
				return nil, pErr
			}
		}
		err = c.Packer.AddSchema(pack, s)
		if err != nil {
			pErr.AddError(err)
			if !options.Force {
				return nil, pErr
			}
		}
	}

	if len(pErr.Errors()) > 0 {
		return pack, pErr
	} else {
		return pack, nil
	}
}

func (c *compilerImpl) ValidateSemantics(manifest *domain.Schema) error {
	return c.Analyzer.Semantics(manifest)
}

func (c *compilerImpl) Parse(fp string, options ParseOptions) (*domain.Schema, error) {
	if err := c.ValidatorService.Struct(options); err != nil {
		return nil, err
	}

	schemaPath := path.Join(options.Root, fp)
	d, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	schma := domain.Schema{}
	err = json.Unmarshal(d, &schma)
	if err != nil {
		return nil, err
	}

	schma.Name = ToSchemaName(fp)

	var extendedConfig []byte
	for i, v := range schma.Extends {
		schma.Extends[i] = path.Join(options.Root, v)
		mani, err := c.Parse(v, options)
		if err != nil {
			return mani, err
		}

		if extendedConfig == nil {
			extendedConfig = mani.Config
		} else {
			extendedConfig, err = c.JsonMerger.MergeJson(mani.Config, extendedConfig)
			if err != nil {
				return nil, err
			}
		}
	}

	var config []byte
	if extendedConfig == nil {
		config = schma.Config
	} else if schma.Config == nil {
		config = extendedConfig
	} else {
		config, err = c.JsonMerger.MergeJson(extendedConfig, schma.Config)
		if err != nil {
			return nil, err
		}
	}

	schma.Config = config
	schma.Path = schemaPath
	if schma.Type == "" {
		schma.Type = "json"
	}

	return &schma, nil
}

func (c *compilerImpl) Compile(fp string, options CompileOptions) (*domain.Schema, error) {
	if err := c.ValidatorService.Struct(options); err != nil {
		return nil, err
	}
	m, err := c.Parse(fp, options.ParseOptions)
	if err != nil {
		return nil, err
	}

	err = c.ValidateSemantics(m)
	if err != nil {
		return nil, err
	}

	if !options.DryRun {
		err := c.FileService.ToFile(options.OutputDirectory, m)
		if err != nil {
			return nil, err
		}
	}

	return m, err
}

// Converts the Schema's filepath into its corresponding name.
func ToSchemaName(fp string) string {
	fp = path.Join(string(os.PathSeparator), fp)
	dirs := strings.Split(fp, string(os.PathSeparator))
	dirs = dirs[1:]
	if len(dirs) <= 0 {
		return ""
	}
	lastInd := len(dirs) - 1
	lastElem := dirs[lastInd]
	dirs[lastInd] = strings.TrimSuffix(lastElem, filepath.Ext(lastElem))

	return strings.Join(dirs, "-")
}
