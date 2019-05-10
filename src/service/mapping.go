package service

import (
	"github.com/two-rabbits/ranvier/src/exchange/response"
	"github.com/two-rabbits/ranvier/src/model"
)

const MappingServiceKey = "MappingService"

type MappingService interface {
	ToLevelConfigMetaData(config *model.LevelConfig) response.LevelConfigMetaData
}

type mappingServiceImpl struct {
}

func (m *mappingServiceImpl) ToLevelConfigMetaData(config *model.LevelConfig) response.LevelConfigMetaData {
	return response.LevelConfigMetaData{
		Name:    config.Name,
		Version: config.Version,
	}
}
