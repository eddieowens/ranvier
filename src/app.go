package src

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/two-rabbits/ranvier/src/docs"
	"github.com/two-rabbits/ranvier/src/filewatcher"
	"github.com/two-rabbits/ranvier/src/router"
	"github.com/two-rabbits/ranvier/src/state"
	"github.com/two-rabbits/ranvier/src/swagger"
)

const AppKey = "App"

type App interface {
	Run()
}

type appImpl struct {
	Router      router.Router           `inject:"Router"`
	FileWatcher filewatcher.FileWatcher `inject:"FileWatcher"`
	State       state.ConfigState       `inject:"ConfigState"`
}

func (a *appImpl) Run() {
	e := echo.New()

	e.Use(middleware.Logger(), middleware.Recover())
	e.GET("/swagger/*", swagger.Handler())

	a.Router.RegisterAll(e)

	if err := a.FileWatcher.Start(); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
