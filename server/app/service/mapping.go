package service

import (
	"github.com/two-rabbits/ranvier/server/app/exchange/response"
	"github.com/two-rabbits/ranvier/server/app/model"
)

const MappingServiceKey = "MappingService"

type MappingService interface {
	ToLevelConfigMeta(config *model.Config) response.ConfigMeta
	ToLevelConfigMetaData(config *model.Config) *response.LevelConfigMetaData
	ToConfig(config *model.Config) response.Config
	ToLevelConfigData(config *model.Config) *response.LevelConfigData
}

type mappingServiceImpl struct {
}

func (m *mappingServiceImpl) ToLevelConfigData(config *model.Config) *response.LevelConfigData {
	if config == nil {
		return nil
	}

	return &response.LevelConfigData{
		Name:   config.Name,
		Config: config.Data,
	}
}

func (m *mappingServiceImpl) ToLevelConfigMetaData(config *model.Config) *response.LevelConfigMetaData {
	if config == nil {
		return nil
	}
	return &response.LevelConfigMetaData{
		Name: config.Name,
	}
}

func (m *mappingServiceImpl) ToConfig(config *model.Config) response.Config {
	if config == nil {
		return response.Config{
			Data: nil,
		}
	}
	return response.Config{
		Data: m.ToLevelConfigData(config),
	}
}

func (m *mappingServiceImpl) ToLevelConfigMeta(config *model.Config) response.ConfigMeta {
	return response.ConfigMeta{
		Data: m.ToLevelConfigMetaData(config),
	}
}
