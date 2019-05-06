package service

import (
	"config-manager/src/model"
	"config-manager/src/state"
)

const GlobalLevelConfigServiceKey = "GlobalLevelConfigService"

type GlobalLevelConfigService interface {
	Query(query string) (model.LevelConfig, error)
	Create(data []byte) (model.LevelConfig, error)
	Rollback(version int) (model.LevelConfig, error)
	Update(data []byte) (config model.LevelConfig, err error)
}

type globalLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
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
