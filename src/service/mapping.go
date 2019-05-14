package service

import (
	"github.com/two-rabbits/ranvier/src/exchange/response"
	"github.com/two-rabbits/ranvier/src/model"
)

const MappingServiceKey = "MappingService"

type MappingService interface {
	ToLevelConfigMeta(config *model.LevelConfig) response.LevelConfigMeta
	ToLevelConfigMetaData(config *model.LevelConfig) *response.LevelConfigMetaData
	ToLevelConfig(config *model.LevelConfig) response.LevelConfig
	ToLevelConfigData(config *model.LevelConfig) *response.LevelConfigData
}

type mappingServiceImpl struct {
}

func (m *mappingServiceImpl) ToLevelConfigData(config *model.LevelConfig) *response.LevelConfigData {
	if config == nil {
		return nil
	}

	return &response.LevelConfigData{
		Name:    config.Name,
		Version: config.Version,
		Config:  config.Config,
	}
}

func (m *mappingServiceImpl) ToLevelConfigMetaData(config *model.LevelConfig) *response.LevelConfigMetaData {
	if config == nil {
		return nil
	}
	return &response.LevelConfigMetaData{
		Name:    config.Name,
		Version: config.Version,
	}
}

func (m *mappingServiceImpl) ToLevelConfig(config *model.LevelConfig) response.LevelConfig {
	if config == nil {
		return response.LevelConfig{
			Data: nil,
		}
	}
	return response.LevelConfig{
		Data: m.ToLevelConfigData(config),
	}
}

func (m *mappingServiceImpl) ToLevelConfigMeta(config *model.LevelConfig) response.LevelConfigMeta {
	return response.LevelConfigMeta{
		Data: m.ToLevelConfigMetaData(config),
	}
}
