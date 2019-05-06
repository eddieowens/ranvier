package controller

import (
	"config-manager/src/model"
	"config-manager/src/service"
	"github.com/labstack/echo"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
)

const GlobalControllerKey = "GlobalController"

type GlobalController interface {
	LevelConfigController
}

type globalControllerImpl struct {
	LevelConfigService service.GlobalLevelConfigService `inject:"GlobalLevelConfigService"`
}

func (g *globalControllerImpl) MergedQuery(c echo.Context) error {
	return g.Query(c)
}

func (g *globalControllerImpl) Update(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !gjson.ValidBytes(data) {
		return c.NoContent(http.StatusBadRequest)
	}

	appConfig, err := g.LevelConfigService.Update(data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, appConfig)
}

func (g *globalControllerImpl) Rollback(c echo.Context) error {
	version := c.Param("version")
	ver, err := strconv.Atoi(version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "version must be an int")
	}
	config, err := g.LevelConfigService.Rollback(ver)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, config)
}

func (g *globalControllerImpl) Query(c echo.Context) error {
	key := c.QueryParam("key")

	data, err := g.LevelConfigService.Query(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (g *globalControllerImpl) Create(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !gjson.ValidBytes(data) {
		return c.NoContent(http.StatusBadRequest)
	}

	globalConfig, err := g.LevelConfigService.Create(data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, globalConfig)
}

func (g *globalControllerImpl) GetRoutes() []model.Route {
	return []model.Route{
		model.NewRoute(http.MethodPost, "/api/admin/config", g.Create),
		model.NewRoute(http.MethodPut, "/api/admin/config/rollback/:version", g.Rollback),
		model.NewRoute(http.MethodPut, "/api/admin/config", g.Update),
		model.NewRoute(http.MethodGet, "/api/admin/config", g.Query),
		model.NewRoute(http.MethodGet, "/api/config", g.MergedQuery),
	}
}
