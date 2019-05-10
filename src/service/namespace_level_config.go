package service

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/two-rabbits/ranvier/src/exchange/response"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/state"
	"net/http"
)

const NamespaceLevelConfigServiceKey = "NamespaceLevelConfigService"

type NamespaceLevelConfigService interface {
	Query(cluster string, namespace string, query string) (config model.LevelConfig, err error)
	MergedQuery(cluster, namespace, query string) (config model.LevelConfig, err error)
	Create(cluster string, namespace string, data []byte) (model.LevelConfig, error)
	Rollback(cluster string, namespace string, version int) (config model.LevelConfig, err error)
	Update(cluster, namespace string, data []byte) (config model.LevelConfig, err error)
	GetAll(clusterName string) (resp response.NamespaceLevelConfigMeta, err error)
}

type namespaceLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	IdService          state.IdService    `inject:"IdService"`
	MappingService     MappingService     `inject:"MappingService"`
}

func (n *namespaceLevelConfigServiceImpl) GetAll(clusterName string) (resp response.NamespaceLevelConfigMeta, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(clusterName)); !exists {
		return resp, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not n valid cluster", clusterName))
	}

	global, _ := n.LevelConfigService.Query(model.Global, state.GlobalId, "")
	cluster, _ := n.LevelConfigService.Query(model.Cluster, n.IdService.ClusterId(clusterName), "")
	namespaces := n.LevelConfigService.GetAll(model.Namespace)

	resp.Data.Global = n.MappingService.ToLevelConfigMetaData(&global)
	resp.Data.Cluster = n.MappingService.ToLevelConfigMetaData(&cluster)

	ns := make([]response.LevelConfigMetaData, len(namespaces))
	for i := range namespaces {
		ns[i] = n.MappingService.ToLevelConfigMetaData(&namespaces[i])
	}
	resp.Data.Namespaces = ns

	return resp, nil
}

func (n *namespaceLevelConfigServiceImpl) Update(cluster, namespace string, data []byte) (config model.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(cluster)); !exists {
		return config, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not n valid cluster", cluster))
	}

	id := n.IdService.NamespaceId(namespace, cluster)

	return n.LevelConfigService.Update(model.Namespace, id, data)
}

func (n *namespaceLevelConfigServiceImpl) Query(cluster string, namespace string, query string) (config model.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(cluster)); !exists {
		return config, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not n valid cluster", cluster))
	}

	id := n.IdService.NamespaceId(namespace, cluster)

	return n.LevelConfigService.Query(model.Namespace, id, query)
}

func (n *namespaceLevelConfigServiceImpl) MergedQuery(cluster, namespace, query string) (config model.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(cluster)); !exists {
		return config, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not n valid cluster", cluster))
	}

	id := n.IdService.NamespaceId(namespace, cluster)

	return n.LevelConfigService.MergedQuery(model.Namespace, id, query)
}

func (n *namespaceLevelConfigServiceImpl) Create(cluster string, namespace string, data []byte) (config model.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(cluster)); !exists {
		return config, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not n valid cluster", cluster))
	}

	id := n.IdService.NamespaceId(namespace, cluster)

	return n.LevelConfigService.Create(model.Namespace, id, data)
}

func (n *namespaceLevelConfigServiceImpl) Rollback(cluster string, namespace string, version int) (config model.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(cluster)); !exists {
		return config, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not n valid cluster", cluster))
	}

	id := n.IdService.NamespaceId(namespace, cluster)

	return n.LevelConfigService.Rollback(model.Namespace, id, version)
}
