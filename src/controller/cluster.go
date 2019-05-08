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

const ClusterControllerKey = "ClusterController"

type ClusterController interface {
	LevelConfigController
}

type clusterControllerImpl struct {
	LevelConfigService service.ClusterLevelConfigService `inject:"ClusterLevelConfigService"`
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
		model.NewRoute(http.MethodPost, "/api/admin/config/:cluster", cc.Create),
		model.NewRoute(http.MethodPut, "/api/admin/config/:cluster/rollback/:version", cc.Rollback),
		model.NewRoute(http.MethodPut, "/api/admin/config/:cluster", cc.Update),
		model.NewRoute(http.MethodGet, "/api/admin/config/:cluster", cc.Query),
		model.NewRoute(http.MethodGet, "/api/config/:cluster", cc.MergedQuery),
	}
}
