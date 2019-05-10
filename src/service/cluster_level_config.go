package service

import (
	"github.com/two-rabbits/ranvier/src/exchange/response"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/state"
)

const ClusterLevelConfigServiceKey = "ClusterLevelConfigService"

type ClusterLevelConfigService interface {
	Query(cluster string, query string) (model.LevelConfig, error)
	MergedQuery(cluster string, query string) (model.LevelConfig, error)
	Create(cluster string, data []byte) (model.LevelConfig, error)
	Update(cluster string, data []byte) (config model.LevelConfig, err error)
	Rollback(cluster string, version int) (model.LevelConfig, error)
	GetAll() (resp response.ClusterLevelConfigMeta, err error)
}

type clusterLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	IdService          state.IdService    `inject:"IdService"`
	MappingService     MappingService     `inject:"MappingService"`
}

func (c *clusterLevelConfigServiceImpl) GetAll() (resp response.ClusterLevelConfigMeta, err error) {
	global, _ := c.LevelConfigService.Query(model.Global, state.GlobalId, "")
	clusters := c.LevelConfigService.GetAll(model.Cluster)

	resp.Data.Global = c.MappingService.ToLevelConfigMetaData(&global)

	clstrs := make([]response.LevelConfigMetaData, len(clusters))
	for i := range clusters {
		clstrs[i] = c.MappingService.ToLevelConfigMetaData(&clusters[i])
	}
	resp.Data.Clusters = clstrs

	return resp, nil
}

func (c *clusterLevelConfigServiceImpl) MergedQuery(cluster string, query string) (model.LevelConfig, error) {
	return c.LevelConfigService.MergedQuery(model.Cluster, c.IdService.ClusterId(cluster), query)
}

func (c *clusterLevelConfigServiceImpl) Update(cluster string, data []byte) (config model.LevelConfig, err error) {
	return c.LevelConfigService.Update(model.Cluster, c.IdService.ClusterId(cluster), data)
}

func (c *clusterLevelConfigServiceImpl) Query(cluster string, query string) (model.LevelConfig, error) {
	return c.LevelConfigService.Query(model.Cluster, c.IdService.ClusterId(cluster), query)
}

func (c *clusterLevelConfigServiceImpl) Create(cluster string, data []byte) (model.LevelConfig, error) {
	return c.LevelConfigService.Create(model.Cluster, c.IdService.ClusterId(cluster), data)
}

func (c *clusterLevelConfigServiceImpl) Rollback(cluster string, version int) (model.LevelConfig, error) {
	return c.LevelConfigService.Rollback(model.Cluster, c.IdService.ClusterId(cluster), version)
}
