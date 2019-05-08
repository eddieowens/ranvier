package controller

import (
	"github.com/labstack/echo"
	"github.com/tidwall/gjson"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/service"
	"io/ioutil"
	"net/http"
	"strconv"
)

const ApplicationControllerKey = "ApplicationController"

type ApplicationController interface {
	LevelConfigController
}

type applicationControllerImpl struct {
	LevelConfigService service.ApplicationLevelConfigService `inject:"ApplicationLevelConfigService"`
}

func (a *applicationControllerImpl) MergedQuery(c echo.Context) error {
	key := c.QueryParam("key")
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	application := c.Param("application")

	if application == "" || namespace == "" || cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := a.LevelConfigService.MergedQuery(cluster, namespace, application, key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (a *applicationControllerImpl) Update(c echo.Context) error {
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	application := c.Param("application")

	if application == "" || namespace == "" || cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !gjson.ValidBytes(data) {
		return c.NoContent(http.StatusBadRequest)
	}

	appConfig, err := a.LevelConfigService.Update(cluster, namespace, application, data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, appConfig)
}

func (a *applicationControllerImpl) Rollback(c echo.Context) error {
	version := c.Param("version")
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	application := c.Param("application")

	if application == "" || namespace == "" || cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	ver, err := strconv.Atoi(version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "version must be an int")
	}

	config, err := a.LevelConfigService.Rollback(cluster, namespace, application, ver)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, config)
}

func (a *applicationControllerImpl) Query(c echo.Context) error {
	key := c.QueryParam("key")
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	application := c.Param("application")

	if application == "" || namespace == "" || cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := a.LevelConfigService.Query(cluster, namespace, application, key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (a *applicationControllerImpl) Create(c echo.Context) error {
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	application := c.Param("application")

	if application == "" || namespace == "" || cluster == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !gjson.ValidBytes(data) {
		return c.NoContent(http.StatusBadRequest)
	}

	appConfig, err := a.LevelConfigService.Create(cluster, namespace, application, data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, appConfig)
}

func (a *applicationControllerImpl) GetRoutes() []model.Route {
	return []model.Route{
		model.NewRoute(http.MethodPost, "/api/admin/config/:cluster/:namespace/:application", a.Create),
		model.NewRoute(http.MethodPut, "/api/admin/config/:cluster/:namespace/:application/rollback/:version", a.Rollback),
		model.NewRoute(http.MethodPut, "/api/admin/config/:cluster/:namespace/:application", a.Update),
		model.NewRoute(http.MethodGet, "/api/admin/config/:cluster/:namespace/:application", a.Query),
		model.NewRoute(http.MethodGet, "/api/config/:cluster/:namespace/:application", a.MergedQuery),
	}
}
