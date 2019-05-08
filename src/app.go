package src

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/two-rabbits/ranvier/src/filewatcher"
	"github.com/two-rabbits/ranvier/src/router"
	"github.com/two-rabbits/ranvier/src/state"
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
