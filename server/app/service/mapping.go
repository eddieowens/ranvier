package service

import (
	"github.com/eddieowens/ranvier/server/app/exchange/response"
	"github.com/eddieowens/ranvier/server/app/model"
)

const MappingServiceKey = "MappingService"

type MappingService interface {
	ToLevelConfigMeta(config *model.Config) response.ConfigMeta
	ToLevelConfigMetaData(config *model.Config) *response.ConfigMetaData
	ToResponse(config *model.Config) *response.Config
}

type mappingServiceImpl struct {
}

func (m *mappingServiceImpl) ToLevelConfigMetaData(config *model.Config) *response.ConfigMetaData {
	if config == nil {
		return nil
	}
	return &response.ConfigMetaData{
		Name: config.Name,
	}
}

func (m *mappingServiceImpl) ToResponse(config *model.Config) *response.Config {
	if config == nil {
		return nil
	}
	return &response.Config{
		Data: config,
	}
}

func (m *mappingServiceImpl) ToLevelConfigMeta(config *model.Config) response.ConfigMeta {
	return response.ConfigMeta{
		Data: m.ToLevelConfigMetaData(config),
	}
}
