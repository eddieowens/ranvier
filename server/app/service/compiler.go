package service

import (
	"encoding/json"
	"github.com/eddieowens/ranvier/lang/compiler"
	"github.com/eddieowens/ranvier/lang/domain"
	"github.com/eddieowens/ranvier/server/app/configuration"
	"github.com/eddieowens/ranvier/server/app/model"
	"path/filepath"
)

const CompilerServiceKey = "CompilerService"

type CompilerService interface {
	Compile(fp string) (*model.Config, *domain.Schema, error)
}

type compilerServiceImpl struct {
	Compiler compiler.Compiler    `inject:"Compiler"`
	Config   configuration.Config `inject:"Config"`
}

func (c *compilerServiceImpl) Compile(fp string) (*model.Config, *domain.Schema, error) {
	fp, _ = filepath.Rel(c.Config.Git.Directory, fp)
	s, err := c.Compiler.Compile(fp, compiler.CompileOptions{
		ParseOptions: compiler.ParseOptions{
			Root: c.Config.Git.Directory,
		},
		OutputDirectory: c.Config.Compiler.OutputDirectory,
	})

	if err != nil {
		return nil, nil, err
	}

	var data interface{}
	_ = json.Unmarshal(s.Config, &data)

	config := &model.Config{
		Name: s.Name,
		Data: data,
	}

	return config, s, nil
}
