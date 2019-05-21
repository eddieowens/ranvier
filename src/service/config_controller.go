package service

import (
	"github.com/two-rabbits/ranvier/src/exchange/response"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/state"
)

const ConfigControllerServiceKey = "ConfigControllerService"

type ConfigControllerService interface {
	Query(query string) (resp response.LevelConfig, err error)
	Create(data []byte) (resp response.LevelConfig, err error)
	Rollback(version int) (resp response.LevelConfig, err error)
	Update(data []byte) (resp response.LevelConfig, err error)
	GetAll() (resp response.LevelConfigMeta, err error)
	Delete() (resp response.LevelConfig, err error)
}

type configControllerServiceImpl struct {
	ConfigService  ConfigService  `inject:"ConfigService"`
	MappingService MappingService `inject:"MappingService"`
}

func (g *configControllerServiceImpl) Delete() (resp response.LevelConfig, err error) {
	config, err := g.ConfigService.Delete(state.GlobalId)
	if err != nil {
		return resp, err
	}

	return g.MappingService.ToConfig(config), nil
}

func (g *configControllerServiceImpl) GetAll() (resp response.LevelConfigMeta, err error) {
	globalConfig := g.ConfigService.GetAll()
	if len(globalConfig) > 0 {
		meta := g.MappingService.ToLevelConfigMeta(&globalConfig[0])
		return meta, nil
	}
	return resp, nil
}

func (g *configControllerServiceImpl) Update(data []byte) (resp response.LevelConfig, err error) {
	config, err := g.ConfigService.Update(state.GlobalId, data)
	if err != nil {
		return resp, err
	}

	return g.MappingService.ToConfig(config), nil
}

func (g *configControllerServiceImpl) Query(query string) (resp response.LevelConfig, err error) {
	config, exists := g.ConfigService.Query(state.GlobalId, query)
	if !exists {
		return resp, NewKeyNotFoundError(query)
	}

	return g.MappingService.ToConfig(config), nil
}

func (g *configControllerServiceImpl) Create(data []byte) (resp response.LevelConfig, err error) {
	config, err := g.ConfigService.Create(state.GlobalId, data)
	if err != nil {
		return resp, err
	}

	return g.MappingService.ToConfig(config), nil
}

func (g *configControllerServiceImpl) Rollback(version int) (resp response.LevelConfig, err error) {
	config, err := g.ConfigService.Rollback(state.GlobalId, version)
	if err != nil {
		return resp, err
	}

	return g.MappingService.ToConfig(config), nil
}
