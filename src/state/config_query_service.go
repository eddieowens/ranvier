package state

import (
	"github.com/json-iterator/go"
	"github.com/two-rabbits/ranvier/src/collections"
	"github.com/two-rabbits/ranvier/src/model"
)

const ConfigQueryServiceKey = "ConfigQueryService"

type ConfigQueryService interface {
	Query(config model.Config, query string) (out model.Config, exists bool)
}

type configQueryServiceImpl struct {
	Json jsoniter.API `inject:"Json"`
}

func (l *configQueryServiceImpl) Query(config model.Config, query string) (out model.Config, exists bool) {
	var raw interface{}
	if query == "" {
		raw = config.Config.GetAll()
	} else {
		raw, exists = config.Config.Get(query)
		if !exists {
			return out, exists
		}
	}
	out = config
	b, _ := l.Json.Marshal(raw)
	out.Config = collections.NewJsonMap(b)
	return out, true
}
