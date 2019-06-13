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
	topic := c.Param("topic")
	if topic == "" {
		return echo.NewHTTPError(http.StatusNotFound, "topic could not be found")
	}

	return w.Websocket.Connect(topic, c.Response(), c.Request(), nil)
}

func (w *websocketControllerImpl) GetRoutes() []model.Route {
	return []model.Route{
		model.NewRoute(http.MethodGet, "/config/ws/:topic", false, w.Connect),
	}
}
