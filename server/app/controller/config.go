package controller

import (
	"github.com/eddieowens/ranvier/server/app/configuration"
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/eddieowens/ranvier/server/app/service"
	"github.com/labstack/echo"
	"net/http"
)

const ConfigControllerKey = "ConfigController"

const ConfigCreateRoute = "/config"
const ConfigQueryRoute = "/config/:name"
const ConfigUpdateRoute = "/config/:name"
const ConfigDeleteRoute = "/config/:name"

type ConfigController interface {
	Controller
	Query(ctx echo.Context) error
	Create(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
}

type configControllerImpl struct {
	ConfigControllerService service.ConfigControllerService `inject:"ConfigControllerService"`
	Config                  configuration.Config            `inject:"Config"`
}

func (c *configControllerImpl) Create(ctx echo.Context) error {
	config := new(model.Config)
	if err := ctx.Bind(config); err != nil {
		return err
	}

	conf, err := c.ConfigControllerService.Create(config)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, conf)
}

func (c *configControllerImpl) Update(ctx echo.Context) error {
	name := ctx.Param("name")
	if name == "" {
		return echo.NewHTTPError(http.StatusNotFound, "the name of the config is required")
	}

	var config interface{}
	if err := ctx.Bind(&config); err != nil {
		return err
	}

	conf, err := c.ConfigControllerService.Update(&model.Config{
		Name: name,
		Data: config,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, conf)
}

func (c *configControllerImpl) Delete(ctx echo.Context) error {
	name := ctx.Param("name")
	if name == "" {
		return echo.NewHTTPError(http.StatusNotFound, "the name of the config is required")
	}

	conf, err := c.ConfigControllerService.Delete(name)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, conf)
}

func (c *configControllerImpl) Query(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	name := ctx.Param("name")

	data, err := c.ConfigControllerService.Query(name, query)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, data)
}

func (c *configControllerImpl) GetRoutes() []model.Route {
	routes := []model.Route{
		model.NewRoute(http.MethodGet, ConfigQueryRoute, c.Query),
	}
	if c.Config.Env == "dev" {
		routes = append(routes, model.NewRoute(http.MethodPost, ConfigCreateRoute, c.Create))
		routes = append(routes, model.NewRoute(http.MethodPut, ConfigUpdateRoute, c.Update))
		routes = append(routes, model.NewRoute(http.MethodDelete, ConfigDeleteRoute, c.Delete))
	}

	return routes
}
