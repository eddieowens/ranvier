package router

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/server/app/controller"
	"github.com/labstack/echo"
)

const RouterKey = "Router"

type Router interface {
	RegisterAll(e *echo.Echo)
}

type routerImpl struct {
	Controllers []axon.Instance `inject:"Controllers"`
}

func (r *routerImpl) RegisterAll(e *echo.Echo) {
	api := e.Group("/api")
	admin := api.Group("/admin")
	for _, inst := range r.Controllers {
		c := inst.GetValue().(controller.Controller)
		for _, route := range c.GetRoutes() {
			if route.IsAdmin {
				admin.Add(route.Method, route.Path, route.HandlerFunc)
			} else {
				api.Add(route.Method, route.Path, route.HandlerFunc)
			}
		}
	}
}
