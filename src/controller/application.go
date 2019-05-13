package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/tidwall/gjson"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/service"
	"github.com/two-rabbits/ranvier/src/state"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const ApplicationControllerKey = "ApplicationController"

type ApplicationController interface {
	StratifiedLevelConfigController
}

type applicationControllerImpl struct {
	LevelConfigService service.ApplicationLevelConfigService `inject:"ApplicationLevelConfigService"`
}

func (a *applicationControllerImpl) Delete(c echo.Context) error {
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	application := c.Param("application")

	if application == "" || namespace == "" || cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := a.LevelConfigService.Delete(cluster, namespace, application)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (a *applicationControllerImpl) GetAll(c echo.Context) error {
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")

	if namespace == "" || cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	resp, err := a.LevelConfigService.GetAll(cluster, namespace)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *applicationControllerImpl) Get(c echo.Context) error {
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	application := c.Param("name")

	if application == "" || namespace == "" || cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := a.LevelConfigService.Get(cluster, namespace, application)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
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

// @Summary Query on a user-defined configuration.
// @Description Retrieves the entirety or a query result of user-defined configuration via a key param. The key param. They key param is a valid gjson (https://github.com/tidwall/gjson#path-syntax) query. For example, to retrieve a nested key, the key param would be outer_key.inner_key.
// @Produce json
// @Router /config/{cluster}/{namespace}/{application} [get]
// @Param key path string false "a gjson valid query"
// @Success 200 {object} model.LevelConfig
// @Failure 404 {object} model.HTTPError
// @Failure 500 {object} model.HTTPError
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

	if strings.Contains(application, state.IdSeparator) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("the application name cannot contain a %s character", state.IdSeparator))
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
		model.NewRoute(http.MethodPost, "/config/:cluster/:namespace/:application", true, a.Create),
		model.NewRoute(http.MethodPut, "/config/:cluster/:namespace/:application/rollback/:version", true, a.Rollback),
		model.NewRoute(http.MethodPut, "/config/:cluster/:namespace/:application", true, a.Update),
		model.NewRoute(http.MethodGet, "/config/:cluster/:namespace/:application", true, a.Query),
		model.NewRoute(http.MethodDelete, "/config/:cluster/:namespace/:application", true, a.Delete),
		model.NewRoute(http.MethodGet, "/config/:cluster/:namespace/applications", true, a.GetAll),
		model.NewRoute(http.MethodGet, "/config/:cluster/:namespace/applications/:name", true, a.Get),
		model.NewRoute(http.MethodGet, "/config/:cluster/:namespace/:application", false, a.MergedQuery),
	}
}
