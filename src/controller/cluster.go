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

const ClusterControllerKey = "ClusterController"

type ClusterController interface {
	StratifiedLevelConfigController
}

type clusterControllerImpl struct {
	LevelConfigService service.ClusterLevelConfigService `inject:"ClusterLevelConfigService"`
}

func (cc *clusterControllerImpl) GetAll(c echo.Context) error {
	resp, err := cc.LevelConfigService.GetAll()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (cc *clusterControllerImpl) Get(c echo.Context) error {
	cluster := c.Param("name")

	if cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := cc.LevelConfigService.Get(cluster)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (cc *clusterControllerImpl) MergedQuery(c echo.Context) error {
	key := c.QueryParam("key")
	cluster := c.Param("cluster")

	if cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := cc.LevelConfigService.MergedQuery(cluster, key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (cc *clusterControllerImpl) Update(c echo.Context) error {
	cluster := c.Param("cluster")

	if cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !gjson.ValidBytes(data) {
		return c.NoContent(http.StatusBadRequest)
	}

	appConfig, err := cc.LevelConfigService.Update(cluster, data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, appConfig)
}

func (cc *clusterControllerImpl) Rollback(c echo.Context) error {
	version := c.Param("version")
	cluster := c.Param("cluster")
	ver, err := strconv.Atoi(version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "version must be an int")
	}
	config, err := cc.LevelConfigService.Rollback(cluster, ver)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, config)
}

func (cc *clusterControllerImpl) Query(c echo.Context) error {
	key := c.QueryParam("key")
	cluster := c.Param("cluster")

	if cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := cc.LevelConfigService.Query(cluster, key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (cc *clusterControllerImpl) Create(c echo.Context) error {
	cluster := c.Param("cluster")
	if cluster == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	if strings.Contains(cluster, state.IdSeparator) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("the cluster name cannot contain a %s character", state.IdSeparator))
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !gjson.ValidBytes(data) {
		return c.NoContent(http.StatusBadRequest)
	}

	clusterConfig, err := cc.LevelConfigService.Create(cluster, data)
	if err != nil {
		return err

	}

	return c.JSON(http.StatusCreated, clusterConfig)
}

func (cc *clusterControllerImpl) GetRoutes() []model.Route {
	return []model.Route{
		model.NewRoute(http.MethodPost, "/config/:cluster", true, cc.Create),
		model.NewRoute(http.MethodPut, "/config/:cluster/rollback/:version", true, cc.Rollback),
		model.NewRoute(http.MethodPut, "/config/:cluster", true, cc.Update),
		model.NewRoute(http.MethodGet, "/config/:cluster", true, cc.Query),
		model.NewRoute(http.MethodGet, "/config/clusters", true, cc.GetAll),
		model.NewRoute(http.MethodGet, "/config/clusters/:name", true, cc.Get),
		model.NewRoute(http.MethodGet, "/config/:cluster", false, cc.MergedQuery),
	}
}
