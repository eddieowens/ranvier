package controller

import (
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/eddieowens/ranvier/server/app/service"
	"github.com/labstack/echo"
	"net/http"
)

const ConfigControllerKey = "ConfigController"

const ConfigQueryRoute = "/config/:name"

type ConfigController interface {
	Controller
	Query(c echo.Context) error
}

type configControllerImpl struct {
	ConfigControllerService service.ConfigControllerService `inject:"ConfigControllerService"`
}

func (g *configControllerImpl) Query(c echo.Context) error {
	query := c.QueryParam("query")
	name := c.Param("name")

	data, err := g.ConfigControllerService.Query(name, query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (g *configControllerImpl) GetRoutes() []model.Route {
	return []model.Route{
		model.NewRoute(http.MethodGet, ConfigQueryRoute, false, g.Query),
	}
}
