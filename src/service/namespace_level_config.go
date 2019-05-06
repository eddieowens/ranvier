package service

import (
	"config-manager/src/model"
	"config-manager/src/state"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

const NamespaceLevelConfigServiceKey = "NamespaceLevelConfigService"

type NamespaceLevelConfigService interface {
	Query(cluster string, namespace string, query string) (config model.LevelConfig, err error)
	MergedQuery(cluster, namespace, query string) (config model.LevelConfig, err error)
	Create(cluster string, namespace string, data []byte) (model.LevelConfig, error)
	Rollback(cluster string, namespace string, version int) (config model.LevelConfig, err error)
	Update(cluster, namespace string, data []byte) (config model.LevelConfig, err error)
}

type namespaceLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	IdService          state.IdService    `inject:"IdService"`
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
