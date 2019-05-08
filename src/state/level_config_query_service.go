package state

import (
	"github.com/json-iterator/go"
	"github.com/two-rabbits/ranvier/src/collections"
	"github.com/two-rabbits/ranvier/src/model"
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
