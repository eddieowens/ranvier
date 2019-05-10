package service

import (
	"github.com/two-rabbits/ranvier/src/exchange/response"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/state"
)

const GlobalLevelConfigServiceKey = "GlobalLevelConfigService"

type GlobalLevelConfigService interface {
	Query(query string) (model.LevelConfig, error)
	Create(data []byte) (model.LevelConfig, error)
	Rollback(version int) (model.LevelConfig, error)
	Update(data []byte) (config model.LevelConfig, err error)
	GetAll() (resp response.LevelConfigMeta, err error)
}

type globalLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	MappingService     MappingService     `inject:"MappingService"`
}

func (g *globalLevelConfigServiceImpl) GetAll() (resp response.LevelConfigMeta, err error) {
	globalConfig := g.LevelConfigService.GetAll(model.Global)
	if len(globalConfig) > 0 {
		meta := g.MappingService.ToLevelConfigMeta(&globalConfig[0])
		return meta, nil
	}
	return resp, nil
}

func (g *globalLevelConfigServiceImpl) Update(data []byte) (config model.LevelConfig, err error) {
	return g.LevelConfigService.Update(model.Global, state.GlobalId, data)
}

func (g *globalLevelConfigServiceImpl) Query(query string) (model.LevelConfig, error) {
	return g.LevelConfigService.Query(model.Global, state.GlobalId, query)
}

func (g *globalLevelConfigServiceImpl) Create(data []byte) (model.LevelConfig, error) {
	return g.LevelConfigService.Create(model.Global, state.GlobalId, data)
}

func (g *globalLevelConfigServiceImpl) Rollback(version int) (model.LevelConfig, error) {
	return g.LevelConfigService.Rollback(model.Global, state.GlobalId, version)
}
