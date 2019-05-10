package service

import (
	"github.com/two-rabbits/ranvier/src/exchange/response"
	"github.com/two-rabbits/ranvier/src/model"
)

const MappingServiceKey = "MappingService"

type MappingService interface {
	ToLevelConfigMeta(config *model.LevelConfig) response.LevelConfigMeta
}

type mappingServiceImpl struct {
}

func (m *mappingServiceImpl) ToLevelConfigMeta(config *model.LevelConfig) response.LevelConfigMeta {
	return response.LevelConfigMeta{
		Name:    config.Name,
		Version: config.Version,
	}
}
