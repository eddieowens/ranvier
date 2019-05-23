package app

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/two-rabbits/ranvier/compiler"
	_ "github.com/two-rabbits/ranvier/server/app/docs"
	"github.com/two-rabbits/ranvier/server/app/poller"
	"github.com/two-rabbits/ranvier/server/app/router"
	"github.com/two-rabbits/ranvier/server/app/swagger"
	"regexp"
	"strings"
)

const AppKey = "App"

type App interface {
	Run()
}

type appImpl struct {
	Router    router.Router    `inject:"Router"`
	GitPoller poller.GitPoller `inject:"GitPoller"`
}

func onUpdate(filepath string) {
	fmt.Println("update ", filepath)
}

func (a *appImpl) Run() {
	e := echo.New()

	e.Use(middleware.Logger(), middleware.Recover())
	e.GET("/swagger/*", swagger.Handler())

	a.Router.RegisterAll(e)

	err := a.GitPoller.Start(onUpdate, *regexp.MustCompile(fmt.Sprintf(".+(%s)", strings.Join(compiler.SupportedFileTypes, "|"))))
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
