package src

import (
	"config-manager/src/filewatcher"
	"config-manager/src/router"
	"config-manager/src/state"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const AppKey = "App"

type App interface {
	Run()
}

type appImpl struct {
	Router      router.Router           `inject:"Router"`
	FileWatcher filewatcher.FileWatcher `inject:"FileWatcher"`
	State       state.LevelConfigState  `inject:"LevelConfigState"`
}

func (a *appImpl) Run() {
	e := echo.New()

	e.Use(middleware.Logger(), middleware.Recover())

	a.Router.RegisterAll(e)

	if err := a.FileWatcher.Start(); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
