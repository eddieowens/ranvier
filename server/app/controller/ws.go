package controller

import (
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/eddieowens/ranvier/server/app/ws"
	"github.com/labstack/echo"
	"net/http"
)

const WebsocketControllerKey = "WebsocketController"

type WebsocketController interface {
	Controller
	Connect(c echo.Context) error
}

type websocketControllerImpl struct {
	Websocket ws.Websocket `inject:"Websocket"`
}

func (w *websocketControllerImpl) Connect(c echo.Context) error {
	name := c.Param("config_name")
	if name == "" {
		return echo.NewHTTPError(http.StatusNotFound, "a config name is required")
	}

	return w.Websocket.Connect(name, c.Response(), c.Request(), nil)
}

func (w *websocketControllerImpl) GetRoutes() []model.Route {
	return []model.Route{
		model.NewRoute(http.MethodGet, "/config/ws/:config_name", w.Connect),
	}
}
