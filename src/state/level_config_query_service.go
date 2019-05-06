package state

import (
	"config-manager/src/collections"
	"config-manager/src/model"
	"github.com/json-iterator/go"
)

const LevelConfigQueryServiceKey = "LevelConfigQueryService"

type LevelConfigQueryService interface {
	Query(levelConfig model.LevelConfig, query string) (config model.LevelConfig, exists bool)
}

type levelConfigQueryServiceImpl struct {
	Json jsoniter.API `inject:"Json"`
}

func (l *levelConfigQueryServiceImpl) Query(levelConfig model.LevelConfig, query string) (config model.LevelConfig, exists bool) {
	var raw interface{}
	if query == "" {
		raw = levelConfig.Config.GetAll()
	} else {
		raw, exists = levelConfig.Config.Get(query)
		if !exists {
			return config, exists
		}
	}
	config = levelConfig
	b, _ := l.Json.Marshal(raw)
	config.Config = collections.NewJsonMap(b)
	return config, true
}
