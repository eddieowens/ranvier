package compiler

import (
	"github.com/json-iterator/go"
	"github.com/two-rabbits/ranvier/compiler/model"
	"github.com/two-rabbits/ranvier/src/service"
	"io/ioutil"
)

const ConfigCompilerKey = "ConfigCompiler"

type ConfigCompiler interface {
	Compile(filepath string) (m model.ConfigManifest, err error)
}

type configCompilerImpl struct {
	Json         jsoniter.API         `inject:"Json"`
	MergeService service.MergeService `inject:"MergeService"`
}

func (c *configCompilerImpl) Compile(filepath string) (m model.ConfigManifest, err error) {
	d, err := ioutil.ReadFile(filepath)
	if err != nil {
		return m, err
	}

	err = c.Json.Unmarshal(d, &m)
	if err != nil {
		return m, err
	}

	var extension []byte
	for _, v := range m.Extends {
		mani, err := c.Compile(v)
		if err != nil {
			return mani, err
		}

		if extension == nil {
			extension = mani.Config
		} else {
			extension, err = c.MergeService.MergeJson(mani.Config, extension)
			if err != nil {
				return m, err
			}
		}
	}

	config, err := c.MergeService.MergeJson(extension, m.Config)
	if err != nil {
		return m, err
	}

	m.Config = config

	return m, err
}
