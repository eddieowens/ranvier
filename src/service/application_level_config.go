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
	Query(cluster string, namespace string, application string, query string) (config model.LevelConfig, err error)
	MergedQuery(cluster, namespace, application string, query string) (config model.LevelConfig, err error)
	Create(cluster string, namespace string, application string, data []byte) (model.LevelConfig, error)
	Rollback(cluster string, namespace string, application string, version int) (config model.LevelConfig, err error)
	Update(cluster, namespace, application string, data []byte) (config model.LevelConfig, err error)
	GetAll(clusterName string, namespaceName string) (resp response.ApplicationLevelConfigMeta, err error)
}

type applicationLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	IdService          state.IdService    `inject:"IdService"`
	MappingService     MappingService     `inject:"MappingService"`
	LevelService       state.LevelService `inject:"LevelService"`
}

func (a *applicationLevelConfigServiceImpl) GetAll(clusterName string, namespaceName string) (resp response.ApplicationLevelConfigMeta, err error) {
	if err := a.exists(clusterName, namespaceName); err != nil {
		return resp, err
	}

	global, _ := a.LevelConfigService.Query(model.Global, state.GlobalId, "")
	cluster, _ := a.LevelConfigService.Query(model.Cluster, a.IdService.ClusterId(clusterName), "")
	namespace, _ := a.LevelConfigService.Query(model.Namespace, a.IdService.NamespaceId(namespaceName, clusterName), "")
	applications := a.LevelConfigService.GetAll(model.Application)

	resp.Data.Global = a.MappingService.ToLevelConfigMetaData(&global)
	resp.Data.Cluster = a.MappingService.ToLevelConfigMetaData(&cluster)
	resp.Data.Namespace = a.MappingService.ToLevelConfigMetaData(&namespace)

	apps := make([]response.LevelConfigMetaData, len(applications))
	for i := range applications {
		apps[i] = a.MappingService.ToLevelConfigMetaData(&applications[i])
	}
	resp.Data.Applications = apps

	return resp, nil
}

func (a *applicationLevelConfigServiceImpl) MergedQuery(cluster, namespace, application string, query string) (config model.LevelConfig, err error) {
	if err := a.exists(cluster, namespace); err != nil {
		return config, err
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	return a.LevelConfigService.MergedQuery(model.Application, id, query)
}

func (a *applicationLevelConfigServiceImpl) Update(cluster, namespace, application string, data []byte) (config model.LevelConfig, err error) {
	if err := a.exists(cluster, namespace); err != nil {
		return config, err
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	return a.LevelConfigService.Update(model.Application, id, data)
}

func (a *applicationLevelConfigServiceImpl) Query(cluster string, namespace string, application string, query string) (config model.LevelConfig, err error) {
	if err := a.exists(cluster, namespace); err != nil {
		return config, err
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	return a.LevelConfigService.Query(model.Application, id, query)
}

func (a *applicationLevelConfigServiceImpl) Create(cluster string, namespace string, application string, data []byte) (config model.LevelConfig, err error) {
	if err = a.exists(cluster, namespace); err != nil {
		return
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	return a.LevelConfigService.Create(model.Application, id, data)
}

func (a *applicationLevelConfigServiceImpl) Rollback(cluster string, namespace string, application string, version int) (config model.LevelConfig, err error) {
	if err = a.exists(cluster, namespace); err != nil {
		return
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	return a.LevelConfigService.Rollback(model.Application, id, version)
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
