package router

import (
	"github.com/labstack/echo"
	"github.com/two-rabbits/ranvier/src/controller"
)

const RouterKey = "Router"

type Router interface {
	RegisterAll(e *echo.Echo)
}

type routerImpl struct {
	Global      controller.GlobalController      `inject:"GlobalController"`
	Cluster     controller.ClusterController     `inject:"ClusterController"`
	Namespace   controller.NamespaceController   `inject:"NamespaceController"`
	Application controller.ApplicationController `inject:"ApplicationController"`
}

func (r *routerImpl) getControllers() []controller.Controller {
	return []controller.Controller{
		r.Global,
		r.Cluster,
		r.Namespace,
		r.Application,
	}
}

func (r *routerImpl) RegisterAll(e *echo.Echo) {
	api := e.Group("/api")
	admin := api.Group("/admin")
	for _, c := range r.getControllers() {
		for _, route := range c.GetRoutes() {
			if route.IsAdmin {
				admin.Add(route.Method, route.Path, route.HandlerFunc)
			} else {
				api.Add(route.Method, route.Path, route.HandlerFunc)
			}
		}
	}
}
