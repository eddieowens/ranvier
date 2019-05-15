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

// @Summary Update the user-defined configuration.
// @Description Updates the user-defined configuration by merging the provided request body with the current configuration. Any key conflicts will be resolved with the provided request body taking precedence.
// @Produce json
// @Accept json
// @Router /config/{cluster}/{namespace}/{application} [put]
// @Param cluster path string true "The name of the cluster that the application resides in"
// @Param namespace path string true "The name of the namespace that the application resides in"
// @Param application path string true "The name of the application"
// @Success 200 {object} response.LevelConfig
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Tags Admin
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

// @Summary Rollback user-defined configuration to a particular version.
// @Description Accepts the integer version of a particular configuration and rolls back configuration to that state. A new version will be created after rolling back.
// @Produce json
// @Router /config/{cluster}/{namespace}/{application}/rollback/{version} [put]
// @Param cluster path string true "The name of the cluster that the application resides in"
// @Param namespace path string true "The name of the namespace that the application resides in"
// @Param application path string true "The name of the application"
// @Param version path int true "The version to rollback to"
// @Success 200 {object} response.LevelConfig
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Tags Admin
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
// @Description Retrieves the entirety or a query result of user-defined configuration via a key param. The key param is a valid gjson (https://github.com/tidwall/gjson#path-syntax) query. For example, to retrieve a nested key, the key param would be outer_key.inner_key.
// @Produce json
// @Router /config/{cluster}/{namespace}/{application} [get]
// @Param cluster path string true "The name of the cluster that the application resides in"
// @Param namespace path string true "The name of the namespace that the application resides in"
// @Param application path string true "The name of the application"
// @Param key query string false "A gjson valid query"
// @Success 200 {object} response.LevelConfig
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Tags Admin
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
