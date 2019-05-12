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
	GetAll(clusterName string, namespaceName string) (resp response.ApplicationsLevelConfigMeta, err error)
	Get(clusterName string, namespaceName string, appName string) (resp response.ApplicationLevelConfigMeta, err error)
	Delete(clusterName string, namespaceName string, appName string) (resp model.LevelConfig, err error)
}

type applicationLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	IdService          state.IdService    `inject:"IdService"`
	MappingService     MappingService     `inject:"MappingService"`
	LevelService       state.LevelService `inject:"LevelService"`
}

func (a *applicationLevelConfigServiceImpl) Delete(clusterName string, namespaceName string, appName string) (resp model.LevelConfig, err error) {
	if err := a.exists(clusterName, namespaceName); err != nil {
		return resp, err
	}

	return a.LevelConfigService.Delete(model.Application, a.IdService.ApplicationId(appName, namespaceName, clusterName))
}

func (a *applicationLevelConfigServiceImpl) Get(clusterName string, namespaceName string, appName string) (resp response.ApplicationLevelConfigMeta, err error) {
	if err := a.exists(clusterName, namespaceName); err != nil {
		return resp, err
	}

	global, _ := a.LevelConfigService.Query(model.Global, state.GlobalId, "")
	cluster, _ := a.LevelConfigService.Query(model.Cluster, a.IdService.ClusterId(clusterName), "")
	namespace, _ := a.LevelConfigService.Query(model.Namespace, a.IdService.NamespaceId(namespaceName, clusterName), "")
	application, err := a.LevelConfigService.Query(model.Application, a.IdService.ApplicationId(appName, namespaceName, clusterName), "")
	if err != nil {
		return resp, err
	}

	resp.Global = a.MappingService.ToLevelConfigMeta(&global)
	resp.Cluster = a.MappingService.ToLevelConfigMeta(&cluster)
	resp.Namespace = a.MappingService.ToLevelConfigMeta(&namespace)
	resp.Application = a.MappingService.ToLevelConfigMeta(&application)

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

	resp.Global = a.MappingService.ToLevelConfigMeta(&global)
	resp.Cluster = a.MappingService.ToLevelConfigMeta(&cluster)
	resp.Namespace = a.MappingService.ToLevelConfigMeta(&namespace)

	apps := make([]response.LevelConfigMeta, len(applications))
	for i := range applications {
		apps[i] = a.MappingService.ToLevelConfigMeta(&applications[i])
	}
	resp.Applications = apps

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
