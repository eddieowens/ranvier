package app

import (
	"fmt"
	"github.com/eddieowens/ranvier/lang/domain"
	"github.com/eddieowens/ranvier/server/app/configuration"
	"github.com/eddieowens/ranvier/server/app/poller"
	"github.com/eddieowens/ranvier/server/app/router"
	"github.com/eddieowens/ranvier/server/app/service"
	"github.com/eddieowens/ranvier/server/app/swagger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

const AppKey = "App"

type App interface {
	Run()
}

type appImpl struct {
	Router          router.Router               `inject:"Router"`
	GitPoller       poller.GitPoller            `inject:"GitPoller"`
	ConfigWsService service.ConfigPollerService `inject:"ConfigPollerService"`
	Config          configuration.Config        `inject:"Config"`
}

func (a *appImpl) Run() {
	e := echo.New()

	format := &log.JSONFormatter{
		TimestampFormat: a.Config.Log.TimeFormat,
	}

	log.SetFormatter(format)
	log.SetLevel(resolveLevel(a.Config.Log.Level))

	if log.GetLevel() >= log.DebugLevel {
		e.Use(middleware.Logger(), middleware.Recover())
	}

	e.HideBanner = true
	e.HidePort = true
	e.GET("/swagger/*", swagger.Handler())

	a.Router.RegisterAll(e)

	err := a.GitPoller.Start(a.ConfigWsService.OnUpdate, a.ConfigWsService.OnStart, *regexp.MustCompile(fmt.Sprintf(".+(%s)", strings.Join(domain.SupportedFileTypes, "|"))))
	if err != nil {
		log.Fatal(err)
	}

	log.Info("starting server on port ", a.Config.Server.Port)
	log.Fatal(e.Start(fmt.Sprintf(":%d", a.Config.Server.Port)))
}

func resolveLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "trace":
		return log.TraceLevel
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn", "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		log.Info(level, " is not a valid log level. Setting to info.")
		return log.InfoLevel
	}
}
