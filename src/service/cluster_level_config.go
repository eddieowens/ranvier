package service

import (
	"github.com/two-rabbits/ranvier/src/exchange/response"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/state"
)

const ClusterLevelConfigServiceKey = "ClusterLevelConfigService"

type ClusterLevelConfigService interface {
	Query(cluster string, query string) (resp response.LevelConfig, err error)
	MergedQuery(cluster string, query string) (resp response.LevelConfig, err error)
	Create(cluster string, data []byte) (resp response.LevelConfig, err error)
	Update(cluster string, data []byte) (resp response.LevelConfig, err error)
	Rollback(cluster string, version int) (resp response.LevelConfig, err error)
	GetAll() (resp response.ClustersLevelConfigMeta, err error)
	Get(name string) (resp response.ClusterLevelConfigMeta, err error)
	Delete(clusterName string) (resp response.LevelConfig, err error)
}

type clusterLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	IdService          state.IdService    `inject:"IdService"`
	MappingService     MappingService     `inject:"MappingService"`
}

func (c *clusterLevelConfigServiceImpl) Delete(clusterName string) (resp response.LevelConfig, err error) {
	config, err := c.LevelConfigService.Delete(model.Cluster, c.IdService.ClusterId(clusterName))
	if err != nil {
		return resp, err
	}

	return c.MappingService.ToLevelConfig(config), nil
}

func (c *clusterLevelConfigServiceImpl) Get(name string) (resp response.ClusterLevelConfigMeta, err error) {
	global, _ := c.LevelConfigService.Query(model.Global, state.GlobalId, "")

	clusterId := c.IdService.ClusterId(name)
	cluster, exists := c.LevelConfigService.Query(model.Cluster, c.IdService.ClusterId(name), "")
	if !exists {
		return resp, NewKeyNotFoundError(clusterId.String())
	}
	resp.Data = &response.ClusterLevelConfigMetaData{}

	resp.Data.Global = c.MappingService.ToLevelConfigMetaData(global)
	resp.Data.Cluster = c.MappingService.ToLevelConfigMetaData(cluster)

	return resp, nil
}

func (c *clusterLevelConfigServiceImpl) GetAll() (resp response.ClustersLevelConfigMeta, err error) {
	global, _ := c.LevelConfigService.Query(model.Global, state.GlobalId, "")
	clusters := c.LevelConfigService.GetAll(model.Cluster)
	resp.Data = &response.ClustersLevelConfigMetaData{}

	resp.Data.Global = c.MappingService.ToLevelConfigMetaData(global)

	clstrs := make([]response.LevelConfigMetaData, len(clusters))
	for i := range clusters {
		clstr := c.MappingService.ToLevelConfigMetaData(&clusters[i])
		clstrs[i] = *clstr
	}
	resp.Data.Clusters = clstrs

	return resp, nil
}

func (c *clusterLevelConfigServiceImpl) MergedQuery(cluster string, query string) (resp response.LevelConfig, err error) {
	config, err := c.LevelConfigService.MergedQuery(model.Cluster, c.IdService.ClusterId(cluster), query)
	if err != nil {
		return resp, err
	}

	return c.MappingService.ToLevelConfig(config), nil
}

func (c *clusterLevelConfigServiceImpl) Update(cluster string, data []byte) (resp response.LevelConfig, err error) {
	config, err := c.LevelConfigService.Update(model.Cluster, c.IdService.ClusterId(cluster), data)
	if err != nil {
		return resp, err
	}

	return c.MappingService.ToLevelConfig(config), nil
}

func (c *clusterLevelConfigServiceImpl) Query(cluster string, query string) (resp response.LevelConfig, err error) {
	clusterId := c.IdService.ClusterId(cluster)
	config, exists := c.LevelConfigService.Query(model.Cluster, clusterId, query)
	if !exists {
		return resp, NewKeyNotFoundError(clusterId.String())
	}

	return c.MappingService.ToLevelConfig(config), nil
}

func (c *clusterLevelConfigServiceImpl) Create(cluster string, data []byte) (resp response.LevelConfig, err error) {
	config, err := c.LevelConfigService.Create(model.Cluster, c.IdService.ClusterId(cluster), data)
	if err != nil {
		return resp, err
	}

	return c.MappingService.ToLevelConfig(config), nil
}

func (c *clusterLevelConfigServiceImpl) Rollback(cluster string, version int) (resp response.LevelConfig, err error) {
	config, err := c.LevelConfigService.Rollback(model.Cluster, c.IdService.ClusterId(cluster), version)
	if err != nil {
		return resp, err
	}

	return c.MappingService.ToLevelConfig(config), nil
}
