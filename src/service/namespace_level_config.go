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
	Query(cluster string, namespace string, query string) (resp response.LevelConfig, err error)
	MergedQuery(cluster, namespace, query string) (resp response.LevelConfig, err error)
	Create(cluster string, namespace string, data []byte) (resp response.LevelConfig, err error)
	Rollback(cluster string, namespace string, version int) (resp response.LevelConfig, err error)
	Update(cluster, namespace string, data []byte) (resp response.LevelConfig, err error)
	GetAll(clusterName string) (resp response.NamespacesLevelConfigMeta, err error)
	Get(clusterName string, namespaceName string) (resp response.NamespaceLevelConfigMeta, err error)
	Delete(clusterName string, namespaceName string) (resp response.LevelConfig, err error)
}

type namespaceLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	IdService          state.IdService    `inject:"IdService"`
	MappingService     MappingService     `inject:"MappingService"`
}

func (n *namespaceLevelConfigServiceImpl) Delete(clusterName string, namespaceName string) (resp response.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(clusterName)); !exists {
		return resp, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid cluster", clusterName))
	}

	config, err := n.LevelConfigService.Delete(model.Namespace, n.IdService.NamespaceId(namespaceName, clusterName))
	if err != nil {
		return resp, err
	}

	return n.MappingService.ToLevelConfig(config), nil
}

func (n *namespaceLevelConfigServiceImpl) Get(clusterName string, namespaceName string) (resp response.NamespaceLevelConfigMeta, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(clusterName)); !exists {
		return resp, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid cluster", clusterName))
	}

	global, _ := n.LevelConfigService.Query(model.Global, state.GlobalId, "")
	cluster, _ := n.LevelConfigService.Query(model.Cluster, n.IdService.ClusterId(clusterName), "")

	namespaceId := n.IdService.NamespaceId(namespaceName, clusterName)
	namespace, exists := n.LevelConfigService.Query(model.Namespace, namespaceId, "")
	if !exists {
		return resp, NewKeyNotFoundError(namespaceId.String())
	}

	resp.Data = &response.NamespaceLevelConfigMetaData{}

	resp.Data.Global = n.MappingService.ToLevelConfigMetaData(global)
	resp.Data.Cluster = n.MappingService.ToLevelConfigMetaData(cluster)
	resp.Data.Namespace = n.MappingService.ToLevelConfigMetaData(namespace)

	return resp, nil
}

func (n *namespaceLevelConfigServiceImpl) GetAll(clusterName string) (resp response.NamespacesLevelConfigMeta, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(clusterName)); !exists {
		return resp, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid cluster", clusterName))
	}

	global, _ := n.LevelConfigService.Query(model.Global, state.GlobalId, "")
	cluster, _ := n.LevelConfigService.Query(model.Cluster, n.IdService.ClusterId(clusterName), "")
	namespaces := n.LevelConfigService.GetAll(model.Namespace)

	resp.Data = &response.NamespacesLevelConfigMetaData{}

	resp.Data.Global = n.MappingService.ToLevelConfigMetaData(global)
	resp.Data.Cluster = n.MappingService.ToLevelConfigMetaData(cluster)

	nss := make([]response.LevelConfigMetaData, len(namespaces))
	for i := range namespaces {
		ns := n.MappingService.ToLevelConfigMetaData(&namespaces[i])
		nss[i] = *ns
	}
	resp.Data.Namespaces = nss

	return resp, nil
}

func (n *namespaceLevelConfigServiceImpl) Update(cluster, namespace string, data []byte) (resp response.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(cluster)); !exists {
		return resp, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid cluster", cluster))
	}

	id := n.IdService.NamespaceId(namespace, cluster)

	config, err := n.LevelConfigService.Update(model.Namespace, id, data)
	if err != nil {
		return resp, err
	}

	return n.MappingService.ToLevelConfig(config), nil
}

func (n *namespaceLevelConfigServiceImpl) Query(cluster string, namespace string, query string) (resp response.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(cluster)); !exists {
		return resp, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid cluster", cluster))
	}

	id := n.IdService.NamespaceId(namespace, cluster)

	config, exists := n.LevelConfigService.Query(model.Namespace, id, query)
	if !exists {
		return resp, NewKeyNotFoundError(query)
	}

	return n.MappingService.ToLevelConfig(config), nil
}

func (n *namespaceLevelConfigServiceImpl) MergedQuery(cluster, namespace, query string) (resp response.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(cluster)); !exists {
		return resp, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid cluster", cluster))
	}

	id := n.IdService.NamespaceId(namespace, cluster)

	config, err := n.LevelConfigService.MergedQuery(model.Namespace, id, query)
	if err != nil {
		return resp, err
	}

	return n.MappingService.ToLevelConfig(config), nil
}

func (n *namespaceLevelConfigServiceImpl) Create(cluster string, namespace string, data []byte) (resp response.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(cluster)); !exists {
		return resp, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid cluster", cluster))
	}

	id := n.IdService.NamespaceId(namespace, cluster)

	config, err := n.LevelConfigService.Create(model.Namespace, id, data)
	if err != nil {
		return resp, err
	}

	return n.MappingService.ToLevelConfig(config), nil
}

func (n *namespaceLevelConfigServiceImpl) Rollback(cluster string, namespace string, version int) (resp response.LevelConfig, err error) {
	if exists := n.LevelConfigService.Exists(model.Cluster, n.IdService.ClusterId(cluster)); !exists {
		return resp, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid cluster", cluster))
	}

	id := n.IdService.NamespaceId(namespace, cluster)

	config, err := n.LevelConfigService.Rollback(model.Namespace, id, version)
	if err != nil {
		return resp, err
	}

	return n.MappingService.ToLevelConfig(config), nil
}
