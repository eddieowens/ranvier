package app

import (
	"fmt"
	"github.com/eddieowens/ranvier/lang/domain"
	"github.com/eddieowens/ranvier/server/app/poller"
	"github.com/eddieowens/ranvier/server/app/router"
	"github.com/eddieowens/ranvier/server/app/service"
	"github.com/eddieowens/ranvier/server/app/swagger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"regexp"
	"strings"
)

const AppKey = "App"

type App interface {
	Run()
}

type appImpl struct {
	Router        router.Router         `inject:"Router"`
	GitPoller     poller.GitPoller      `inject:"GitPoller"`
	ConfigService service.ConfigService `inject:"ConfigService"`
}

func (a *appImpl) Run() {
	e := echo.New()

	e.Use(middleware.Logger(), middleware.Recover())
	e.GET("/swagger/*", swagger.Handler())

	a.Router.RegisterAll(e)

	err := a.GitPoller.Start(a.ConfigService.SetFromFile, *regexp.MustCompile(fmt.Sprintf(".+(%s)", strings.Join(domain.SupportedFileTypes, "|"))))
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
