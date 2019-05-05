package service

import (
	"config-manager/src/model"
	"config-manager/src/state"
)

const ClusterLevelConfigServiceKey = "ClusterLevelConfigService"

type ClusterLevelConfigService interface {
	Query(cluster string, query string) (*model.LevelConfig, error)
	Create(cluster string, data []byte) (model.LevelConfig, error)
	Update(cluster string, data []byte) (config model.LevelConfig, err error)
	Rollback(cluster string, version int) (model.LevelConfig, error)
}

type clusterLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	IdService          state.IdService    `inject:"IdService"`
}

func (c *clusterLevelConfigServiceImpl) Update(cluster string, data []byte) (config model.LevelConfig, err error) {
	return c.LevelConfigService.Update(model.Cluster, c.IdService.ClusterId(cluster), data)
}

func (c *clusterLevelConfigServiceImpl) Query(cluster string, query string) (*model.LevelConfig, error) {
	return c.LevelConfigService.Query(model.Cluster, c.IdService.ClusterId(cluster), query)
}

func (c *clusterLevelConfigServiceImpl) Create(cluster string, data []byte) (model.LevelConfig, error) {
	return c.LevelConfigService.Create(model.Cluster, c.IdService.ClusterId(cluster), data)
}

func (c *clusterLevelConfigServiceImpl) Rollback(cluster string, version int) (model.LevelConfig, error) {
	return c.LevelConfigService.Rollback(model.Cluster, c.IdService.ClusterId(cluster), version)
}
