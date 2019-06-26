package controller

import (
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/labstack/echo"
	"net/http"
)

const HealthCheckControllerKey = "HealthCheckController"

const GetHealthRoute = "/health"

type HealthCheckController interface {
	Controller
	GetHealth(c echo.Context) error
}

type healthCheckControllerImpl struct {
}

func (h *healthCheckControllerImpl) GetHealth(c echo.Context) error {
	return c.NoContent(200)
}

func (h *healthCheckControllerImpl) GetRoutes() []model.Route {
	return []model.Route{
		model.NewRoute(http.MethodGet, GetHealthRoute, false, h.GetHealth),
	}
}
