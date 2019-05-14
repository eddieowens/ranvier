package service

import (
	"github.com/two-rabbits/ranvier/src/exchange/response"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/state"
)

const GlobalLevelConfigServiceKey = "GlobalLevelConfigService"

type GlobalLevelConfigService interface {
	Query(query string) (resp response.LevelConfig, err error)
	Create(data []byte) (resp response.LevelConfig, err error)
	Rollback(version int) (resp response.LevelConfig, err error)
	Update(data []byte) (resp response.LevelConfig, err error)
	GetAll() (resp response.LevelConfigMeta, err error)
	Delete() (resp response.LevelConfig, err error)
}

type globalLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	MappingService     MappingService     `inject:"MappingService"`
}

func (g *globalLevelConfigServiceImpl) Delete() (resp response.LevelConfig, err error) {
	config, err := g.LevelConfigService.Delete(model.Global, state.GlobalId)
	if err != nil {
		return resp, err
	}

	return g.MappingService.ToLevelConfig(config), nil
}

func (g *globalLevelConfigServiceImpl) GetAll() (resp response.LevelConfigMeta, err error) {
	globalConfig := g.LevelConfigService.GetAll(model.Global)
	if len(globalConfig) > 0 {
		meta := g.MappingService.ToLevelConfigMeta(&globalConfig[0])
		return meta, nil
	}
	return resp, nil
}

func (g *globalLevelConfigServiceImpl) Update(data []byte) (resp response.LevelConfig, err error) {
	config, err := g.LevelConfigService.Update(model.Global, state.GlobalId, data)
	if err != nil {
		return resp, err
	}

	return g.MappingService.ToLevelConfig(config), nil
}

func (g *globalLevelConfigServiceImpl) Query(query string) (resp response.LevelConfig, err error) {
	config, exists := g.LevelConfigService.Query(model.Global, state.GlobalId, query)
	if !exists {
		return resp, NewKeyNotFoundError(query)
	}

	return g.MappingService.ToLevelConfig(config), nil
}

func (g *globalLevelConfigServiceImpl) Create(data []byte) (resp response.LevelConfig, err error) {
	config, err := g.LevelConfigService.Create(model.Global, state.GlobalId, data)
	if err != nil {
		return resp, err
	}

	return g.MappingService.ToLevelConfig(config), nil
}

func (g *globalLevelConfigServiceImpl) Rollback(version int) (resp response.LevelConfig, err error) {
	config, err := g.LevelConfigService.Rollback(model.Global, state.GlobalId, version)
	if err != nil {
		return resp, err
	}

	return g.MappingService.ToLevelConfig(config), nil
}
