package service

import (
	"github.com/json-iterator/go"
	"github.com/oliveagle/jsonpath"
	"github.com/two-rabbits/ranvier/server/model"
	"github.com/two-rabbits/ranvier/server/state"
)

const ConfigQueryServiceKey = "ConfigQueryService"

type ConfigQueryService interface {
	Query(name string, query string) (*model.Config, error)
}

type configQueryServiceImpl struct {
	Json      jsoniter.API    `inject:"Json"`
	ConfigMap state.ConfigMap `inject:"ConfigMap"`
}

func (l *configQueryServiceImpl) Query(name string, query string) (*model.Config, error) {
	config, exists := l.ConfigMap.Get(name)
	if !exists {
		return nil, nil
	}

	raw, err := jsonpath.JsonPathLookup(config.Data, query)
	if err != nil {
		return nil, err
	}

	return &model.Config{
		Name: name,
		Data: raw,
	}, nil
}
