package controller

import (
	"github.com/labstack/echo"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/service"
	"net/http"
)

const ConfigControllerKey = "ConfigController"

type ConfigController interface {
	Controller
	Query(c echo.Context) error
}

type configControllerImpl struct {
	ConfigControllerService service.ConfigControllerService `inject:"ConfigControllerService"`
}

func (g *configControllerImpl) Query(c echo.Context) error {
	key := c.QueryParam("key")
	name := c.QueryParam("name")

	data, err := g.ConfigControllerService.Query(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (g *configControllerImpl) GetRoutes() []model.Route {
	return []model.Route{
		model.NewRoute(http.MethodGet, "/config/:name", true, g.Query),
	}
}
