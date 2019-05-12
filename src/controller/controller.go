package controller

import (
	"github.com/labstack/echo"
	"github.com/two-rabbits/ranvier/src/model"
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
	MergedQuery(c echo.Context) error
	GetAll(c echo.Context) error
	Delete(c echo.Context) error
}

type StratifiedLevelConfigController interface {
	LevelConfigController
	Get(c echo.Context) error
}
