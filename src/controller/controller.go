package controller

import (
	"config-manager/src/model"
	"github.com/labstack/echo"
)

type Controller interface {
	GetRoutes() []model.Route
}

type LevelConfigController interface {
	Controller
	Query(c echo.Context) error
	Create(c echo.Context) error
	Rollback(c echo.Context) error
	Update(c echo.Context) error
}
