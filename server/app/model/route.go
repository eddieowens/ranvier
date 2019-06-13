package model

import "github.com/labstack/echo"

type Route struct {
	Method      string
	Path        string
	HandlerFunc echo.HandlerFunc
	IsAdmin     bool
}

func NewRoute(method string, path string, isAdmin bool, handler echo.HandlerFunc) Route {
	return Route{
		Method:      method,
		Path:        path,
		HandlerFunc: handler,
		IsAdmin:     isAdmin,
	}
}
