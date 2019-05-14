package service

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/two-rabbits/ranvier/src/exchange/response"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/state"
	"net/http"
)

const ApplicationLevelConfigServiceKey = "ApplicationLevelConfigService"

type ApplicationLevelConfigService interface {
	Query(cluster string, namespace string, application string, query string) (resp response.LevelConfig, err error)
	MergedQuery(cluster, namespace, application string, query string) (resp response.LevelConfig, err error)
	Create(cluster string, namespace string, application string, data []byte) (resp response.LevelConfig, err error)
	Rollback(cluster string, namespace string, application string, version int) (resp response.LevelConfig, err error)
	Update(cluster, namespace, application string, data []byte) (resp response.LevelConfig, err error)
	GetAll(clusterName string, namespaceName string) (resp response.ApplicationsLevelConfigMeta, err error)
	Get(clusterName string, namespaceName string, appName string) (resp response.ApplicationLevelConfigMeta, err error)
	Delete(clusterName string, namespaceName string, appName string) (resp response.LevelConfig, err error)
}

type applicationLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	IdService          state.IdService    `inject:"IdService"`
	MappingService     MappingService     `inject:"MappingService"`
	LevelService       state.LevelService `inject:"LevelService"`
}

func (a *applicationLevelConfigServiceImpl) Delete(clusterName string, namespaceName string, appName string) (resp response.LevelConfig, err error) {
	if err := a.exists(clusterName, namespaceName); err != nil {
		return resp, err
	}

	config, err := a.LevelConfigService.Delete(model.Application, a.IdService.ApplicationId(appName, namespaceName, clusterName))
	if err != nil {
		return resp, err
	}

	return a.MappingService.ToLevelConfig(config), err
}

func (a *applicationLevelConfigServiceImpl) Get(clusterName string, namespaceName string, appName string) (resp response.ApplicationLevelConfigMeta, err error) {
	if err := a.exists(clusterName, namespaceName); err != nil {
		return resp, err
	}

	global, _ := a.LevelConfigService.Query(model.Global, state.GlobalId, "")
	cluster, _ := a.LevelConfigService.Query(model.Cluster, a.IdService.ClusterId(clusterName), "")
	namespace, _ := a.LevelConfigService.Query(model.Namespace, a.IdService.NamespaceId(namespaceName, clusterName), "")

	appId := a.IdService.ApplicationId(appName, namespaceName, clusterName)
	application, exists := a.LevelConfigService.Query(model.Application, appId, "")
	if !exists {
		return resp, NewKeyNotFoundError(appId.String())
	}

	resp.Data = &response.ApplicationLevelConfigMetaData{}

	resp.Data.Global = a.MappingService.ToLevelConfigMetaData(global)
	resp.Data.Cluster = a.MappingService.ToLevelConfigMetaData(cluster)
	resp.Data.Namespace = a.MappingService.ToLevelConfigMetaData(namespace)
	resp.Data.Application = a.MappingService.ToLevelConfigMetaData(application)

	return resp, nil
}

func (a *applicationLevelConfigServiceImpl) GetAll(clusterName string, namespaceName string) (resp response.ApplicationsLevelConfigMeta, err error) {
	if err := a.exists(clusterName, namespaceName); err != nil {
		return resp, err
	}

	global, _ := a.LevelConfigService.Query(model.Global, state.GlobalId, "")
	cluster, _ := a.LevelConfigService.Query(model.Cluster, a.IdService.ClusterId(clusterName), "")
	namespace, _ := a.LevelConfigService.Query(model.Namespace, a.IdService.NamespaceId(namespaceName, clusterName), "")
	applications := a.LevelConfigService.GetAll(model.Application)

	resp.Data = &response.ApplicationsLevelConfigMetaData{}

	resp.Data.Global = a.MappingService.ToLevelConfigMetaData(global)
	resp.Data.Cluster = a.MappingService.ToLevelConfigMetaData(cluster)
	resp.Data.Namespace = a.MappingService.ToLevelConfigMetaData(namespace)

	apps := make([]response.LevelConfigMetaData, len(applications))
	for i := range applications {
		app := a.MappingService.ToLevelConfigMetaData(&applications[i])
		apps[i] = *app
	}
	resp.Data.Applications = apps

	return resp, nil
}

func (a *applicationLevelConfigServiceImpl) MergedQuery(cluster, namespace, application string, query string) (resp response.LevelConfig, err error) {
	if err := a.exists(cluster, namespace); err != nil {
		return resp, err
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	config, err := a.LevelConfigService.MergedQuery(model.Application, id, query)
	if err != nil {
		return resp, err
	}

	return a.MappingService.ToLevelConfig(config), nil
}

func (a *applicationLevelConfigServiceImpl) Update(cluster, namespace, application string, data []byte) (resp response.LevelConfig, err error) {
	if err := a.exists(cluster, namespace); err != nil {
		return resp, err
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	config, err := a.LevelConfigService.Update(model.Application, id, data)
	if err != nil {
		return resp, err
	}
	return a.MappingService.ToLevelConfig(config), nil
}

func (a *applicationLevelConfigServiceImpl) Query(cluster string, namespace string, application string, query string) (resp response.LevelConfig, err error) {
	if err := a.exists(cluster, namespace); err != nil {
		return resp, err
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	config, exists := a.LevelConfigService.Query(model.Application, id, query)
	if !exists {
		return resp, NewKeyNotFoundError(query)
	}

	return a.MappingService.ToLevelConfig(config), nil
}

func (a *applicationLevelConfigServiceImpl) Create(cluster string, namespace string, application string, data []byte) (resp response.LevelConfig, err error) {
	if err = a.exists(cluster, namespace); err != nil {
		return
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	config, err := a.LevelConfigService.Create(model.Application, id, data)
	if err != nil {
		return resp, err
	}

	return a.MappingService.ToLevelConfig(config), nil
}

func (a *applicationLevelConfigServiceImpl) Rollback(cluster string, namespace string, application string, version int) (resp response.LevelConfig, err error) {
	if err = a.exists(cluster, namespace); err != nil {
		return
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	config, err := a.LevelConfigService.Rollback(model.Application, id, version)
	if err != nil {
		return resp, err
	}

	return a.MappingService.ToLevelConfig(config), nil
}

func (a *applicationLevelConfigServiceImpl) exists(cluster string, namespace string) error {
	if !a.LevelConfigService.Exists(model.Cluster, a.IdService.ClusterId(cluster)) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid cluster", cluster))
	}
	if !a.LevelConfigService.Exists(model.Namespace, a.IdService.NamespaceId(namespace, cluster)) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid namespace", namespace))
	}
	return nil
}
