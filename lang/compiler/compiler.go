package compiler

import (
	"github.com/eddieowens/ranvier/commons"
	"github.com/eddieowens/ranvier/lang/domain"
	"github.com/eddieowens/ranvier/lang/semantics"
	"github.com/eddieowens/ranvier/lang/services"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"path"
)

const CompilerKey = "Compiler"

type CompileAllOptions struct {
	CompileOptions
	Force bool
}

type CompileOptions struct {
	OutputDirectory string `validate:"filepath,required_without=DryRun"`
	DryRun          bool
}

type Compiler interface {
	Compile(filepath string, options *CompileOptions) (*domain.Schema, error)
	CompileAll(path string, options *CompileAllOptions) (Pack, error)
	ValidateSemantics(manifest *domain.Schema) error
	Parse(filepath string) (*domain.Schema, error)
}

type compilerImpl struct {
	JsonMerger    commons.JsonMerger     `inject:"JsonMerger"`
	Analyzer      semantics.Analyzer     `inject:"Analyzer"`
	FileCollector services.FileCollector `inject:"FileCollector"`
	FileService   services.FileService   `inject:"FileService"`
	Packer        Packer                 `inject:"Packer"`
}

func (c *compilerImpl) CompileAll(path string, options *CompileAllOptions) (Pack, error) {
	files := c.FileCollector.Collect(path)
	pErr := NewPackError()
	pack := NewPack(path)
	for _, f := range files {
		s, err := c.Compile(f, &options.CompileOptions)
		if err != nil {
			pErr.AddError(err)
			if !options.Force {
				return nil, pErr
			}
		}
		if !s.IsAbstract {
			err := c.Packer.AddSchema(pack, s)
			if err != nil {
				pErr.AddError(err)
				if !options.Force {
					return nil, pErr
				}
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

func (c *compilerImpl) Parse(filepath string) (*domain.Schema, error) {
	d, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	schma := domain.Schema{}
	err = json.Unmarshal(d, &schma)
	if err != nil {
		return nil, err
	}

	var extendedConfig []byte
	for i, v := range schma.Extends {
		dir, _ := path.Split(filepath)
		extPath := path.Join(dir, v)
		schma.Extends[i] = extPath
		mani, err := c.Parse(extPath)
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
	schma.Path = filepath
	schma.IsAbstract = schma.Name == ""
	if schma.Type == "" {
		schma.Type = "json"
	}

	return &schma, nil
}

func (c *compilerImpl) Compile(filepath string, options *CompileOptions) (*domain.Schema, error) {
	m, err := c.Parse(filepath)
	if err != nil {
		return nil, err
	}

	err = c.ValidateSemantics(m)
	if err != nil {
		return nil, err
	}

	if !options.DryRun && !m.IsAbstract {
		err := c.FileService.ToFile(options.OutputDirectory, m)
		if err != nil {
			return nil, err
		}
	}

	return m, err
}
